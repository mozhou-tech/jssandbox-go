package jssandbox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBug1_StatePollution(t *testing.T) {
	ctx := context.Background()
	code := `
		// If 'old_var' exists from previous call, it's state pollution
		var result = [];
		if (typeof old_var !== 'undefined') {
			result.push({role: 'system', content: 'polluted: ' + old_var});
		}
		result.push({role: 'user', content: 'current: ' + current_var});
		result;
	`
	tpl, err := NewPromptTemplate(ctx, &Config{
		Code: code,
	})
	assert.NoError(t, err)

	// First call - set 'old_var'
	msgs1, err := tpl.Format(ctx, map[string]any{
		"old_var":     "value1",
		"current_var": "call1",
	})
	assert.NoError(t, err)
	// msgs1 will have 2 messages because we passed both variables
	assert.Len(t, msgs1, 2)

	// Second call - ONLY set 'current_var'
	msgs2, err := tpl.Format(ctx, map[string]any{
		"current_var": "call2",
	})
	assert.NoError(t, err)

	// If Bug 1 is fixed, msgs2 should only have 1 message (current_var)
	// If Bug 1 persists, msgs2 will have 2 messages because 'old_var' from first call remains
	assert.Len(t, msgs2, 1, "State pollution: 'old_var' persisted from previous call")
	assert.Equal(t, "current: call2", msgs2[0].Content)
}

func TestBug3_DeleteFailure(t *testing.T) {
	ctx := context.Background()
	// Code that makes 'pollute' non-configurable
	code := `
		Object.defineProperty(globalThis, 'pollute', {
			configurable: false
		});
		[{role: 'user', content: 'done'}]
	`
	tpl, err := NewPromptTemplate(ctx, &Config{
		Code: code,
	})
	assert.NoError(t, err)

	_, err = tpl.Format(ctx, map[string]any{
		"pollute": "initial",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to clean up injected variables")
	assert.Contains(t, err.Error(), "pollute")
}

func TestBug2_InconsistentState(t *testing.T) {
	ctx := context.Background()

	// Test image_url with missing imageURL field
	code := `
		[{role: 'user', multiContent: [{type: 'image_url'}]}]
	`
	tpl, err := NewPromptTemplate(ctx, &Config{
		Code: code,
	})
	assert.NoError(t, err)

	_, err = tpl.Format(ctx, nil)
	// Now it should return an error instead of an inconsistent state
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "imageURL field is missing or invalid")
}
