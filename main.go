package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	//"github.com/bradford-hamilton/go-graphql-api/postgres"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Traestan/shorturl/endpoint"
	"github.com/Traestan/shorturl/repository"
	"github.com/Traestan/shorturl/transport"
	"github.com/go-kit/kit/log"
)

var HostSite string = "http://localhost"

func main() {

	var (
		httpPort = flag.String("http.port", "5230", "HTTP listen port")
		//jwtSecret = flag.String("jwt.secret", "4b4e7436167e", "secret for jwt encoding")
		//HostSite = flag.String("host", "http://localhost", "use in qrcode")
	)

	flag.Parse()
	//logger
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")
	//mongo connect
	clientOptions := options.Client().ApplyURI("mongodb://shorturl:123123@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Log("msg", "mongo panic")
	}

	urlEncoder := repository.NewURLEncoder(&repository.URLEncoderConfig{})
	mng := repository.NewModel(client, logger, urlEncoder)
	//endpoint init
	sep := endpoint.MakeEndpoints(&mng)

	errc := make(chan error)

	logger.Log("msg", httpPort)
	//init server
	go func() {
		errc <- http.ListenAndServe(":"+*httpPort, transport.MakeHTTPHandler(sep, &mng, logger))
	}()

	logger.Log("exit", <-errc)
}
