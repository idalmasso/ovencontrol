package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/glog"
	"github.com/idalmasso/ovencontrol/backend/config"
)

type controllerMachine interface {
	SetOnButtonPress(func())
}

// PiServer
type MachineServer struct {
	configuration *config.Config
	initialized   bool
	Router        chi.Router
	machine       controllerMachine
}

// ListenAndServe is the main server procedure that only wraps http.ListenAndServe
func (s *MachineServer) ListenAndServe() {
	if glog.V(3) {
		glog.Infoln("MachineServer -  MachineServer.ListenAndServe start")
	}
	if !s.initialized {
		panic("Server not initialized")
	}
	if glog.V(3) {
		glog.Infoln("MachineServer -  MachineServer.starting on port", s.configuration.Server.Port)
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(s.configuration.Server.Port), s.Router); err != nil {
		panic("Cannot listen on server: " + err.Error())
	}
}

// Init initialize the server router and set the controllerMachine needed to do the work
func (s *MachineServer) Init(machine controllerMachine) {
	if glog.V(3) {
		glog.Infoln("MachineServer -  MachineServer.Init start")
	}
	s.configuration = &config.Config{}
	if err := s.configuration.ReadFromFile("configuration.yaml"); err != nil {
		panic("cannot read configuration file")
	}
	s.machine = machine
	s.updateMachineFromConfig()
	s.Router = chi.NewRouter()
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	FileServer(s.Router.(*chi.Mux), s.configuration.Server.DistributionDirectory)
	s.Router.Route("/api", func(router chi.Router) {
		router.Route("/processes", func(processRouter chi.Router) {
		})
		router.Route("/configuration", func(configRouter chi.Router) {
			//configRouter.Get("/", s.getConfig)
			//configRouter.Put("/", s.updateConfig)
		})
	})

	s.initialized = true
}

func (s *MachineServer) updateMachineFromConfig() {

	s.machine.SetOnButtonPress(func() {

	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// FileServer is serving static files.
func FileServer(router *chi.Mux, root string) {
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
