package balancer

import (
	"context"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/zrpc/internal/clientinterceptors"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc/balancer"
)

// WrapDoneWithBreakerCtx returns a Done callback that wraps breaker logic.
// It attaches breaker state reporting to the request completion callback.
// If breaker strategy is node-level (from ctx), it uses a per-node breaker;
// otherwise, it uses a NopBreaker that never opens.
func WrapDoneWithBreakerCtx(ctx context.Context, name string, done func(balancer.DoneInfo)) (func(info balancer.DoneInfo), error) {
	brk := breaker.NopBreaker()
	if clientinterceptors.IsBreakerStrategyNode(ctx) {
		brk = breaker.GetBreaker(name)
	}

	promise, err := brk.AllowCtx(ctx)
	if err != nil {
		return nil, err
	}

	return func(info balancer.DoneInfo) {
		if info.Err == nil || codes.Acceptable(info.Err) {
			promise.Accept()
		} else {
			promise.Reject(info.Err.Error())
		}

		if done != nil {
			done(info)
		}
	}, nil
}
