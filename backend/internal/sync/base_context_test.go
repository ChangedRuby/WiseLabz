package sync

import (
	"context"
	"testing"
)

func TestBaseContext(t *testing.T) {
	e := NewEngine(nil, nil, nil, nil, "")
	if e.BaseContext() == nil || e.BaseContext().Err() != nil {
		t.Fatal("default base context must be live")
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.SetBaseContext(ctx)
	cancel()
	if e.BaseContext().Err() == nil {
		t.Fatal("base context must reflect shutdown cancellation")
	}
}
