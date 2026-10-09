package transformer

import (
	"context"

	"github.com/routatic/proxy/pkg/types"
)

type startUsageKey struct{}

/*
 * WithStartUsage attaches the usage to report in message_start. Upstreams
 * only report real usage at the end of a stream, but Claude Code records
 * each content block as it completes, using the message_start usage. A
 * carried-forward estimate keeps those entries from showing zero tokens.
 */
func WithStartUsage(ctx context.Context, u types.Usage) context.Context {
	return context.WithValue(ctx, startUsageKey{}, u)
}

// startUsage returns the usage attached by WithStartUsage, or zero usage.
func startUsage(ctx context.Context) types.Usage {
	if ctx == nil {
		return types.Usage{}
	}
	u, _ := ctx.Value(startUsageKey{}).(types.Usage)
	return u
}
