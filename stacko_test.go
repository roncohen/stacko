// Copyright 2014 Christoffer Hallas.

package stacko

import (
	"testing"
)

// TestNewStacktrace creates a new stacktrace skipping 0 frames. If the creation
// of the stacktrace is correct, this should gives us a stacktrace of 4 frames.
func TestNewStacktrace(t *testing.T) {
	stacktrace, err := NewStacktrace(0)
	if err != nil {
		t.Error(err)
	}

	if len(stacktrace) != 4 {
		t.Error("Stacktrace should be 4 frames")
	}
}
