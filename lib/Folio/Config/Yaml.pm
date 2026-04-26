# ABOUTME: Minimal read-only YAML parser for the config file subset.
# ABOUTME: Handles scalars, one-level maps, comments, and boolean normalisation.
package Folio::Config::Yaml;

use strict;
use warnings;
use utf8;
use open ':std', ':encoding(UTF-8)';

our $VERSION = '0.1.0';

# Parse a YAML file and return a two-level hash.
# Dies with file path and line number on malformed input.
sub parse_file {
    my ($class, $path) = @_;

    open(my $fh, '<:encoding(UTF-8)', $path)
        or die "Error: cannot open config file: $path: $!\n";

    my %result;
    my $current_map_key;
    my $line_num = 0;

    while (my $line = <$fh>) {
        $line_num += 1;
        chomp $line;

        # Skip blank lines
        next if $line =~ /^\s*$/;

        # Skip comment-only lines
        next if $line =~ /^\s*#/;

        # Indented line — belongs to current map
        if ($line =~ /^(\s+)(\S.*)$/) {
            my $indent = $1;
            my $content = $2;

            if (!defined $current_map_key) {
                die "Error: $path:$line_num: unexpected indented line outside a map\n";
            }

            # Tabs are not allowed in YAML indentation
            if ($indent =~ /\t/) {
                die "Error: $path:$line_num: tabs are not allowed in YAML indentation\n";
            }

            my ($key, $value) = _parse_key_value($content, $path, $line_num);
            $result{$current_map_key}{$key} = _normalise_value($value);
            next;
        }

        # Top-level line
        $current_map_key = undef;

        my ($key, $value) = _parse_key_value($line, $path, $line_num);

        if (!defined $value || $value eq '') {
            # This key introduces a map — subsequent indented lines are its children
            $current_map_key = $key;
            $result{$key} = {} if !exists $result{$key};
        } else {
            $result{$key} = _normalise_value($value);
        }
    }

    close $fh;
    return \%result;
}

# Parse "key: value" from a line. Handles quoted values and inline comments.
sub _parse_key_value {
    my ($text, $path, $line_num) = @_;

    # Strip inline comment (but not inside quotes)
    # Simple approach: only strip # preceded by whitespace and not inside quotes
    my $stripped = _strip_inline_comment($text);

    if ($stripped =~ /^([^\s:][^:]*?)\s*:\s*(.*)$/) {
        my $key = $1;
        my $value = $2;
        $key =~ s/\s+$//;
        $value = _unquote($value, $path, $line_num);
        return ($key, $value);
    }

    die "Error: $path:$line_num: invalid YAML syntax: $text\n";
}

# Remove inline comments (# preceded by space, not inside quotes)
sub _strip_inline_comment {
    my ($text) = @_;

    my $in_single = 0;
    my $in_double = 0;
    my $result = '';

    for my $i (0 .. length($text) - 1) {
        my $ch = substr($text, $i, 1);

        if ($ch eq "'" && !$in_double) {
            $in_single = !$in_single;
        } elsif ($ch eq '"' && !$in_single) {
            $in_double = !$in_double;
        } elsif ($ch eq '#' && !$in_single && !$in_double) {
            # Check if preceded by whitespace
            if ($i > 0 && substr($text, $i - 1, 1) =~ /\s/) {
                $result =~ s/\s+$//;
                return $result;
            }
        }

        $result .= $ch;
    }

    return $result;
}

# Remove surrounding quotes. Dies on unmatched quotes.
sub _unquote {
    my ($value, $path, $line_num) = @_;

    $value =~ s/^\s+//;
    $value =~ s/\s+$//;

    return $value if $value eq '';

    if ($value =~ /^"(.*)"$/) {
        return $1;
    }
    if ($value =~ /^'(.*)'$/) {
        return $1;
    }

    # Check for unmatched quotes
    if ($value =~ /^"/ && $value !~ /"$/) {
        die "Error: $path:$line_num: unmatched double quote: $value\n";
    }
    if ($value =~ /^'/ && $value !~ /'$/) {
        die "Error: $path:$line_num: unmatched single quote: $value\n";
    }

    return $value;
}

# Normalise boolean strings to 1/0. Pass through everything else.
sub _normalise_value {
    my ($value) = @_;

    return $value if !defined $value || $value eq '';

    my $lower = lc $value;
    return 1 if $lower eq 'true'  || $lower eq 'yes' || $lower eq 'on';
    return 0 if $lower eq 'false' || $lower eq 'no'  || $lower eq 'off';

    return $value;
}

1;
