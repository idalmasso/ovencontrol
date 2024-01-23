package main

import (
	"flag"
	"log/slog"

	"github.com/idalmasso/ovencontrol/backend/dummyinterface"
	"github.com/idalmasso/ovencontrol/backend/server"
)

func init() {
	flag.Parse()

}

func main() {
	slog.Info("backend start process")

	//controller := hwinterface.NewController()
	controller := &dummyinterface.DummyController{}

	server := server.MachineServer{}
	server.Init(controller)
	server.ListenAndServe()

}
