package clientinterceptors

import (
	"context"
	"path"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc"
)

const (
	breakerStrategyService = "service"
	breakerStrategyNode    = "node"
)

// BreakerInterceptor is an interceptor that acts as a circuit breaker.
func BreakerInterceptor(breakerStrategy string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if breakerStrategy == breakerStrategyNode {
			ctx = withBreakerStrategy(ctx, breakerStrategyNode)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		breakerName := path.Join(cc.Target(), method)
		return breaker.DoWithAcceptableCtx(ctx, breakerName, func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}, codes.Acceptable)
	}
}

type withBreakerStrategyKey struct{}

func withBreakerStrategy(ctx context.Context, strategy string) context.Context {
	return context.WithValue(ctx, withBreakerStrategyKey{}, strategy)
}

func getBreakerStrategy(ctx context.Context) string {
	strategy, ok := ctx.Value(withBreakerStrategyKey{}).(string)
	if !ok {
		return breakerStrategyService
	}
	return strategy
}

func IsBreakerStrategyNode(ctx context.Context) bool {
	return getBreakerStrategy(ctx) == breakerStrategyNode
}
