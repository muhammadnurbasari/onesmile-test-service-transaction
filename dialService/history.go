package dialService

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"google.golang.org/grpc"
)

func ServiceHistory() generate.TransactionsClient {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	DIAL_HISTORY_SERVICE := os.Getenv("DIAL_HISTORY_SERVICE")
	conn, err := grpc.Dial(DIAL_HISTORY_SERVICE, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("could not connect to" + DIAL_HISTORY_SERVICE + " error: " + err.Error())
		os.Exit(1)
	}

	return generate.NewTransactionsClient(conn)
}
