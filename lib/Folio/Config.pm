# ABOUTME: Config loading with layered merge and namespaced access.
# ABOUTME: Reads script.yaml from local and global locations, merges per-key, applies CLI overrides.
package Folio::Config;

use strict;
use warnings;
use utf8;

use Folio::Config::Yaml;

our $VERSION = '0.1.0';

# Built-in defaults
my %DEFAULTS = (
    'title'                    => undef,
    'subtitle'                 => undef,
    'author'                   => undef,
    'draft-date'               => undef,
    'render-stage-directions'  => 1,
    'render-intro'             => 1,
    'render-footnotes'         => 1,
    'render-character-table'   => 1,
    'folio' => {
        'font'              => 'New Computer Modern',
        'font-size'         => '12pt',
        'margin'            => '25mm',
        'page'              => 'a4',
        'indent'            => '4em',
        'dialogue-spacing'  => '1.6em',
        'direction-spacing' => '1.6em',
        'direction-italic'  => 1,
        'direction-center'  => 0,
        'default-format'    => undef,
    },
);

# Load config from all sources, merge in precedence order, return a config object.
#
# Arguments (named):
#   source_dir => directory containing the source file (for local script.yaml)
#   cli        => hashref of CLI flag overrides (keys match folio: namespace keys)
#
# Returns: a Folio::Config object
sub load {
    my ($class, %args) = @_;

    my $source_dir = $args{source_dir};
    my $cli        = $args{cli} || {};

    # Start with deep copy of defaults
    my %merged = _deep_copy_config(\%DEFAULTS);

    # Layer 1: global config
    my $global_path = "$ENV{HOME}/.config/first-folio/script.yaml";
    if (-f $global_path) {
        my $global = Folio::Config::Yaml->parse_file($global_path);
        _merge_config(\%merged, $global);
    }

    # Layer 2: local config (source file directory)
    if (defined $source_dir && $source_dir ne '') {
        my $local_path = "$source_dir/script.yaml";
        if (-f $local_path) {
            my $local = Folio::Config::Yaml->parse_file($local_path);
            _merge_config(\%merged, $local);
        }
    }

    # Layer 3: CLI flags (these map to folio: namespace keys)
    if (%$cli) {
        for my $key (keys %$cli) {
            next if !defined $cli->{$key};
            if (exists $merged{folio}{$key}) {
                $merged{folio}{$key} = $cli->{$key};
            }
        }
    }

    return bless { config => \%merged }, $class;
}

# Get a config value. Supports dotted access: get('folio.font')
sub get {
    my ($self, $key) = @_;

    if ($key =~ /^(\w[\w-]*)\.\s*(.+)$/) {
        my ($ns, $subkey) = ($1, $2);
        if (ref $self->{config}{$ns} eq 'HASH') {
            return $self->{config}{$ns}{$subkey};
        }
        return undef;
    }

    return $self->{config}{$key};
}

# Deep copy the defaults hash (one level of nesting)
sub _deep_copy_config {
    my ($source) = @_;
    my %copy;
    for my $key (keys %$source) {
        if (ref $source->{$key} eq 'HASH') {
            $copy{$key} = { %{$source->{$key}} };
        } else {
            $copy{$key} = $source->{$key};
        }
    }
    return %copy;
}

# Merge an overlay hash onto a base hash. Scalars override; hashes merge one level deep.
sub _merge_config {
    my ($base, $overlay) = @_;

    for my $key (keys %$overlay) {
        if (ref $overlay->{$key} eq 'HASH') {
            # Merge map keys individually
            $base->{$key} = {} if !exists $base->{$key} || ref $base->{$key} ne 'HASH';
            for my $subkey (keys %{$overlay->{$key}}) {
                $base->{$key}{$subkey} = $overlay->{$key}{$subkey};
            }
        } else {
            $base->{$key} = $overlay->{$key};
        }
    }
}

1;
