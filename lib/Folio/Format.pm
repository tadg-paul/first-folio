# ABOUTME: Extension-to-format mapping and validation for First Folio.
# ABOUTME: Maps file extensions and --to values to canonical format names.
package Folio::Format;

use strict;
use warnings;

our $VERSION = '0.1.0';

my %EXT_TO_FORMAT = (
    '.org'      => 'org',
    '.md'       => 'markdown',
    '.fountain' => 'fountain',
    '.ftn'      => 'fountain',
    '.pdf'      => 'pdf',
);

my %NAME_TO_FORMAT = (
    'org'      => 'org',
    'md'       => 'markdown',
    'markdown' => 'markdown',
    'fountain' => 'fountain',
    'ftn'      => 'fountain',
    'pdf'      => 'pdf',
);

my %READABLE = (
    'org'      => 1,
    'markdown' => 1,
    'fountain' => 1,
);

# Determine format from a file path's extension.
# Returns (format_name, undef) on success or (undef, error_message) on failure.
sub from_extension {
    my ($class, $path) = @_;

    if ($path =~ /(\.[^.]+)$/) {
        my $ext = lc $1;
        if (my $fmt = $EXT_TO_FORMAT{$ext}) {
            return ($fmt, undef);
        }
        return (undef, "Unrecognised file extension: $ext");
    }

    return (undef, "No file extension on: $path");
}

# Determine format from a --to value.
# Returns (format_name, undef) on success or (undef, error_message) on failure.
sub from_name {
    my ($class, $name) = @_;

    my $fmt = $NAME_TO_FORMAT{lc $name};
    if ($fmt) {
        return ($fmt, undef);
    }

    return (undef, "Unrecognised format: $name");
}

# Check whether a format can be read (used as source).
sub is_readable {
    my ($class, $format) = @_;
    return $READABLE{$format} ? 1 : 0;
}

1;
