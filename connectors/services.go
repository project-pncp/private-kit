package connectors

import (
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Pncp() grpc.ClientConnInterface {
	return dial(dialAddress("PNCP_HOST"))
}

func dial(host string) grpc.ClientConnInterface {
	msgSize := 25000000
	grpc.EnableTracing = true
	c, err := grpc.NewClient(host,
		grpc.WithMaxHeaderListSize(uint32(msgSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(msgSize),
			grpc.MaxCallSendMsgSize(msgSize),
		),
	)

	if err != nil {
		panic(err)
	}

	return c
}

func dialAddress(adr string) string {
	d := os.Getenv(adr)
	if len(d) == 0 {
		d = strings.ReplaceAll(adr, "_HOST", "")
		d = strings.ReplaceAll(d, "_", "-")
		d = strings.ToLower(d)
		d = d + ":9090"
	}

	return d
}
