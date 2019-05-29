package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("New is nil")
	} else {
		tracer.Trace("hello tracer")

		if buf.String() != "hello tracer\n" {
			t.Errorf("Error: '%s'", buf.String())
		}
	}

}
