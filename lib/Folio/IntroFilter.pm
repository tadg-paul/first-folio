# ABOUTME: Post-parse filter that distinguishes intro sections from act/scene content.
# ABOUTME: Buffers events before first character dialogue, reorders: front matter, char table, intro, play.
package Folio::IntroFilter;

use strict;
use warnings;

our $VERSION = '0.1.0';

# Wrap an emitter with intro detection. Returns a new emitter hashref.
#
# Buffers all events until the first `character` event. At that point:
# 1. The last act_header before the character event marks the play boundary.
# 2. Events are sorted into three groups:
#    a) Front matter and character table events — emitted first, as-is
#    b) Everything else before the boundary — reclassified as intro_header/intro_text
#    c) Play content from the boundary onwards — emitted as-is
#
# The character table heading: if an act_header immediately precedes
# character_table_start, its text is preserved as the table label.
# Otherwise "Characters" is used. The heading is consumed (not emitted
# as a separate act_header or intro_header).
sub wrap {
    my ($class, $emitter) = @_;

    my @buffer;
    my $boundary_found = 0;

    my $emit = sub {
        my ($event, @args) = @_;
        if (my $cb = $emitter->{$event}) {
            $cb->(@args);
        }
    };

    my $flush_buffer = sub {
        # Find the play boundary: the last act_header before the first character event
        my $boundary_idx = -1;
        for my $i (reverse 0 .. $#buffer) {
            if ($buffer[$i][0] eq 'act_header') {
                $boundary_idx = $i;
                last;
            }
        }

        # Identify character table events and the heading above the table
        my @char_table_events;
        my $char_table_heading;
        my %skip_indices;

        for my $i (0 .. $#buffer) {
            my $event = $buffer[$i][0];
            if ($event eq 'character_table_start' ||
                $event eq 'character_table_row' ||
                $event eq 'character_table_end') {
                push @char_table_events, $buffer[$i];
                $skip_indices{$i} = 1;

                # Check if the event before table_start is an act_header — that's the table heading
                if ($event eq 'character_table_start' && $i > 0 && $buffer[$i - 1][0] eq 'act_header') {
                    $char_table_heading = $buffer[$i - 1][1];
                    $skip_indices{$i - 1} = 1;
                }
            }
        }

        # Phase 1: emit front matter events
        for my $i (0 .. $#buffer) {
            next if $skip_indices{$i};
            my ($event, @args) = @{$buffer[$i]};
            if ($event eq 'front_matter') {
                $emit->($event, @args);
                $skip_indices{$i} = 1;
            }
        }

        # Phase 2: emit character table (with heading label)
        if (@char_table_events) {
            for my $entry (@char_table_events) {
                my ($event, @args) = @$entry;
                if ($event eq 'character_table_start') {
                    # Pass the heading label as a second argument if the emitter supports it
                    $emit->($event, $char_table_heading // 'Characters');
                } else {
                    $emit->($event, @args);
                }
            }
        }

        # Phase 3: emit intro content (everything before boundary, not front_matter or char table)
        for my $i (0 .. $#buffer) {
            next if $skip_indices{$i};
            next if $i >= $boundary_idx;
            my ($event, @args) = @{$buffer[$i]};

            if ($event eq 'act_header') {
                $emit->('intro_header', @args);
            } elsif ($event eq 'dialogue' || $event eq 'stage_direction' || $event eq 'scene_header') {
                $emit->('intro_text', @args);
            }
            # Other events in intro range are silently dropped (shouldn't be any)
        }

        # Phase 4: emit play content (from boundary onwards)
        for my $i ($boundary_idx .. $#buffer) {
            next if $skip_indices{$i};
            my ($event, @args) = @{$buffer[$i]};
            $emit->($event, @args);
        }

        @buffer = ();
        $boundary_found = 1;
    };

    my $wrapped = {};

    my %boundary_events = (
        character => 1,
    );

    for my $event (keys %$emitter) {
        $wrapped->{$event} = sub {
            my @args = @_;

            if ($boundary_found) {
                $emit->($event, @args);
                return;
            }

            if ($boundary_events{$event}) {
                push @buffer, [$event, @args];
                $flush_buffer->();
                return;
            }

            push @buffer, [$event, @args];
        };
    }

    # Ensure intro event callbacks exist
    for my $new_event (qw(intro_header intro_text)) {
        if (!exists $wrapped->{$new_event}) {
            $wrapped->{$new_event} = sub {};
        }
    }

    # Override end to flush remaining buffer (documents with no dialogue)
    $wrapped->{end} = sub {
        if (!$boundary_found && @buffer) {
            # No character found — everything is intro
            # Still emit front matter and char table first
            my @char_events;
            my %skip;
            my $ct_heading;

            for my $i (0 .. $#buffer) {
                my $ev = $buffer[$i][0];
                if ($ev eq 'character_table_start' || $ev eq 'character_table_row' || $ev eq 'character_table_end') {
                    push @char_events, $buffer[$i];
                    $skip{$i} = 1;
                    if ($ev eq 'character_table_start' && $i > 0 && $buffer[$i - 1][0] eq 'act_header') {
                        $ct_heading = $buffer[$i - 1][1];
                        $skip{$i - 1} = 1;
                    }
                }
            }

            for my $i (0 .. $#buffer) {
                next if $skip{$i};
                if ($buffer[$i][0] eq 'front_matter') {
                    $emit->($buffer[$i][0], @{$buffer[$i]}[1..$#{$buffer[$i]}]);
                    $skip{$i} = 1;
                }
            }

            if (@char_events) {
                for my $entry (@char_events) {
                    my ($ev, @args) = @$entry;
                    if ($ev eq 'character_table_start') {
                        $emit->($ev, $ct_heading // 'Characters');
                    } else {
                        $emit->($ev, @args);
                    }
                }
            }

            for my $i (0 .. $#buffer) {
                next if $skip{$i};
                my ($ev, @args) = @{$buffer[$i]};
                if ($ev eq 'act_header') {
                    $emit->('intro_header', @args);
                } elsif ($ev eq 'dialogue' || $ev eq 'stage_direction' || $ev eq 'scene_header') {
                    $emit->('intro_text', @args);
                }
            }
            @buffer = ();
        }
        $emit->('end');
    };

    return $wrapped;
}

1;
