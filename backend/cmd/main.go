package main

import (
	"flag"
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/idalmasso/ovencontrol/backend/hwinterface"
	"github.com/idalmasso/ovencontrol/backend/server"
)

func init() {
	flag.Parse()

}

func main() {
	slog.Info("backend start process")
	logger := httplog.NewLogger("oven-logger", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"env": "dev",
		},
		QuietDownRoutes: []string{
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	controller := hwinterface.NewController(hwinterface.WithLogger(logger))
	//controller := dummyinterface.NewDummyController(dummyinterface.WithLogger(logger))

	server := server.NewMachineServer()
	server.Init(controller, logger)
	server.ListenAndServe()

}
