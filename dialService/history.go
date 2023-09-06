package dialService

import (
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"google.golang.org/grpc"
)

func ServiceHistory(grpcClientConn *grpc.ClientConn) generate.TransactionsClient {
	return generate.NewTransactionsClient(grpcClientConn)
}
