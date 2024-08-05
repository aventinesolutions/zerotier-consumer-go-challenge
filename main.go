package main

import (
	zap "go.uber.org/zap"
	"net/http"
	"os"
)

const (
	ServerAddr = ":4444"
)

// WebHook Secret is provided through the environment via a GCP Secret
var psk = os.Getenv("ZEROTIER_ONE_WEBHOOK_SECRET")

var Logger = zap.NewExample().Sugar()

func main() {
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			panic("unable to defer Zap logging, exiting!")
		}
	}(Logger)
	Logger.Info("starting ZeroTier Consumer Coding Challenge")

	SetRoutes()
	err2 := http.ListenAndServe(ServerAddr, nil)
	if err2 != nil {
		Logger.Errorf("error starting Web Service: %s", err2)
		panic("unable to start Web Service, exiting!")
	}
}
