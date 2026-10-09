package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/routatic/proxy/pkg/types"
)

const usageCarryMaxEntries = 256

/*
 * usageCarry remembers the final usage of each conversation's last turn.
 * Consecutive turns of a tool loop have nearly the same prompt, so the
 * previous turn's usage, including its cache split, is a close estimate for
 * the next turn's message_start. It is an estimate only: the terminal
 * message_delta still carries the real usage.
 */
type usageCarry struct {
	mu sync.Mutex
	m  map[string]types.Usage
}

func newUsageCarry() *usageCarry {
	return &usageCarry{m: make(map[string]types.Usage)}
}

// get and put accept a nil receiver so handlers built without a store still work.
func (c *usageCarry) get(key string) types.Usage {
	if c == nil {
		return types.Usage{}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[key]
}

func (c *usageCarry) put(key string, u types.Usage) {
	if c == nil || u.InputTokens+u.CacheReadInputTokens+u.CacheCreationInputTokens == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.m[key]; !ok && len(c.m) >= usageCarryMaxEntries {
		// shortcut: drop everything at the cap, upgrade to LRU if churn matters.
		c.m = make(map[string]types.Usage)
	}
	c.m[key] = u
}

/*
 * conversationKey identifies a conversation by model, system prompt and first
 * message. These stay fixed across turns while later messages grow.
 */
func conversationKey(req *types.MessageRequest) string {
	if req == nil {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(req.Model))
	h.Write(req.System)
	if len(req.Messages) > 0 {
		b, _ := json.Marshal(req.Messages[0])
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
