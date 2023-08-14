package dialService

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func ServiceCreditCard() generate.ValidationClient {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	DIAL_CREDITCARD_SERVICE := os.Getenv("DIAL_CREDITCARD_SERVICE")
	conn, err := grpc.Dial(DIAL_CREDITCARD_SERVICE, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("could not connect to" + DIAL_CREDITCARD_SERVICE + " error: " + err.Error())
		os.Exit(1)
	}

	return generate.NewValidationClient(conn)
}
