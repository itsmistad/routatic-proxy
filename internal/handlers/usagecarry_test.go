package handlers

import (
	"encoding/json"
	"testing"

	"github.com/routatic/proxy/pkg/types"
)

func TestUsageCarryKeyAndStore(t *testing.T) {
	msg := func(text string) types.Message {
		return types.Message{Role: "user", Content: json.RawMessage(`"` + text + `"`)}
	}
	a := &types.MessageRequest{Model: "m", Messages: []types.Message{msg("first")}}
	grown := &types.MessageRequest{Model: "m", Messages: []types.Message{msg("first"), msg("second")}}
	other := &types.MessageRequest{Model: "m", Messages: []types.Message{msg("different")}}

	if conversationKey(a) != conversationKey(grown) {
		t.Error("key must stay stable as the conversation grows")
	}
	if conversationKey(a) == conversationKey(other) {
		t.Error("different conversations must not share a key")
	}

	c := newUsageCarry()
	if got := c.get("k"); got != (types.Usage{}) {
		t.Errorf("empty store returned %+v", got)
	}
	c.put("k", types.Usage{})
	if got := c.get("k"); got != (types.Usage{}) {
		t.Error("zero usage must not be stored")
	}
	want := types.Usage{InputTokens: 5, OutputTokens: 2, CacheReadInputTokens: 90}
	c.put("k", want)
	if got := c.get("k"); got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
