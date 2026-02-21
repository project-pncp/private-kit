package decode

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/marco-kit/private-kit/query"
)

func GRPCParams(ctx context.Context, md metadata.MD) context.Context {
	filter := query.Filter{}
	filter.DecodeGRPC(md)
	ctx = context.WithValue(ctx, "filter", filter)
	return ctx
}
