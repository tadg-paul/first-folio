# ABOUTME: Config loading with deep layered merge and dotted-path access.
# ABOUTME: Layers: British preset → global script.yaml → local script.yaml → CLI flags.
package Folio::Config;

use strict;
use warnings;
use utf8;

use File::Basename qw(dirname);
use YAML::Tiny;

our $VERSION = '0.2.0';

# Path to the canonical British preset (embedded in the distribution)
my $PRESET_DIR;
BEGIN {
    require FindBin;
    # When run from bin/folio, RealBin is bin/ so ../examples works.
    # When run from project root scripts, RealBin is the project root.
    my $candidate = "$FindBin::RealBin/../examples";
    if (-d $candidate) {
        $PRESET_DIR = $candidate;
    } elsif (-d "$FindBin::RealBin/examples") {
        $PRESET_DIR = "$FindBin::RealBin/examples";
    } else {
        # Walk up from lib/ to find examples/
        my $lib_dir = __FILE__;
        $lib_dir =~ s{/lib/Folio/Config\.pm$}{};
        $PRESET_DIR = "$lib_dir/examples";
    }
}

# Load config from all sources, merge in precedence order, return a config object.
#
# Layers (lowest to highest priority):
#   1. British preset (canonical defaults)
#   2. Global ~/.config/first-folio/script.yaml
#   3. Local <source_dir>/script.yaml
#   4. CLI flags
#
# Arguments (named):
#   source_dir => directory containing the source file
#   cli        => hashref of CLI flag overrides (flat keys mapped into folio:)
#   style      => 'british' (default) or 'american'
#
# Returns: a Folio::Config object
sub load {
    my ($class, %args) = @_;

    my $source_dir = $args{source_dir};
    my $cli        = $args{cli} || {};
    my $style      = $args{style} || 'british';

    # Layer 1: British preset (always the base)
    my $base = _load_yaml_file("$PRESET_DIR/british-script.yaml")
        or die "Error: cannot load British preset from $PRESET_DIR/british-script.yaml\n";

    # Layer 1b: American overrides if requested
    if ($style eq 'american') {
        my $us = _load_yaml_file("$PRESET_DIR/us-overrides-script.yaml");
        _deep_merge($base, $us) if $us;
    }

    # Layer 2: global config
    my $global_path = "$ENV{HOME}/.config/first-folio/script.yaml";
    if (-f $global_path) {
        my $global = _load_yaml_file($global_path);
        _deep_merge($base, $global) if $global;
    }

    # Layer 3: local config (source file directory)
    if (defined $source_dir && $source_dir ne '') {
        my $local_path = "$source_dir/script.yaml";
        if (-f $local_path) {
            my $local = _load_yaml_file($local_path);
            _deep_merge($base, $local) if $local;
        }
    }

    # Layer 4: CLI flags (mapped into folio: namespace)
    if (%$cli) {
        for my $key (keys %$cli) {
            next if !defined $cli->{$key};
            $base->{folio}{$key} = $cli->{$key} if exists $base->{folio}{$key};
        }
    }

    return bless { config => $base }, $class;
}

# Get a config value using dotted path: get('folio.positioning.speech.speaker.bold')
# Boolean strings (true/false/yes/no/on/off) are normalised to 1/0.
sub get {
    my ($self, $path) = @_;

    my @parts = split /\./, $path;
    my $node = $self->{config};

    for my $part (@parts) {
        if (ref $node eq 'HASH' && exists $node->{$part}) {
            $node = $node->{$part};
        } else {
            return undef;
        }
    }

    return _normalise_bool($node);
}

sub _normalise_bool {
    my ($value) = @_;
    return $value if !defined $value || ref $value;

    my $lower = lc $value;
    return 1 if $lower eq 'true'  || $lower eq 'yes' || $lower eq 'on';
    return 0 if $lower eq 'false' || $lower eq 'no'  || $lower eq 'off';
    return $value;
}

# Convenience: return the entire config hash (for inspection/debugging)
sub as_hash {
    my ($self) = @_;
    return $self->{config};
}

# --- Private helpers ---

sub _load_yaml_file {
    my ($path) = @_;

    my $yaml = YAML::Tiny->read($path);
    if (!$yaml) {
        die "Error: cannot parse YAML: $path: " . YAML::Tiny->errstr . "\n";
    }

    return $yaml->[0];
}

# Deep recursive merge: overlay values override base values.
# Hashes merge recursively; scalars are replaced.
sub _deep_merge {
    my ($base, $overlay) = @_;

    for my $key (keys %$overlay) {
        if (ref $overlay->{$key} eq 'HASH' && ref($base->{$key} // '') eq 'HASH') {
            _deep_merge($base->{$key}, $overlay->{$key});
        } else {
            $base->{$key} = $overlay->{$key};
        }
    }
}

1;
