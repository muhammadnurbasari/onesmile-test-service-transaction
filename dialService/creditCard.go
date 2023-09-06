package dialService

import (
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"google.golang.org/grpc"
)

func ServiceCreditCard(grpcClientConn *grpc.ClientConn) generate.ValidationClient {
	return generate.NewValidationClient(grpcClientConn)
}
