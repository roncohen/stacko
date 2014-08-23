package stacko

import (
  "testing"
)

func TestNewStacktrace (t *testing.T) {
  stacktrace, err := NewStacktrace(0)
  if err != nil {
    t.Error(err)
  }

  if len(stacktrace) != 4 {
    t.Error("Stacktrace should be 4 frames.")
  }
}
