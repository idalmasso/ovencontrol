package main

//NOTE: glog levels:
// V=1 for ERROR
// V=2 for WARNING
// V=3 for INFO
// V=4 for DEBUG

import (
	"flag"

	"github.com/golang/glog"
	"github.com/idalmasso/ovencontrol/backend/dummyinterface"
	"github.com/idalmasso/ovencontrol/backend/server"
)

func init() {
	flag.Parse()
	if !isFlagPassed("v") {
		flag.Set("v", "2")
	}
	flag.Set("logtostderr", "1")

}

func main() {
	defer glog.Flush()
	if glog.V(3) {
		glog.Infoln("backend start process")
	}
	//controller := hwinterface.NewController()
	controller := &dummyinterface.DummyController{}

	server := server.MachineServer{}
	server.Init(controller)
	server.ListenAndServe()

}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
