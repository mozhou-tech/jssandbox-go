package jssandbox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSandbox_Delete_NonConfigurable(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// Create a non-configurable property via JavaScript
	_, err := sb.Run(`
		Object.defineProperty(globalThis, 'nonDeletable', {
			value: 42,
			configurable: false
		});
	`)
	assert.NoError(t, err)

	// Verify it exists
	val := sb.Get("nonDeletable")
	assert.NotNil(t, val)
	assert.Equal(t, int64(42), val.ToInteger())

	// Try to delete it using the current Delete implementation
	errDel := sb.Delete("nonDeletable")
	assert.Error(t, errDel, "Deletion should return an error for non-configurable property")

	// Verify it STILL exists if it was non-configurable
	valAfter := sb.Get("nonDeletable")
	assert.NotNil(t, valAfter)
	assert.Equal(t, int64(42), valAfter.ToInteger(), "Property should still exist because it was non-configurable")
}
