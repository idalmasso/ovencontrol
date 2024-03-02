package server

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/idalmasso/ovencontrol/backend/config"
	"github.com/idalmasso/ovencontrol/backend/ovenprograms"
)

type controllerMachine interface {
	temperatureReader
	ovenprograms.Oven
	InitConfig(c config.Config)
	Terminate()
}

// PiServer
type MachineServer struct {
	ovenProgramManager ovenprograms.OvenProgramManager
	configuration      *config.Config
	initialized        bool
	Router             chi.Router
	machine            controllerMachine
	ovenProgramWorker  *ovenprograms.OvenProgramWorker
	logger             *httplog.Logger
}

// ListenAndServe is the main server procedure that only wraps http.ListenAndServe
func (s *MachineServer) ListenAndServe() {

	if !s.initialized {
		panic("Server not initialized")
	}
	defer s.machine.Terminate()
	if err := http.ListenAndServe(":"+strconv.Itoa(s.configuration.Server.Port), s.Router); err != nil {

		panic("Cannot listen on server: " + err.Error())
	}
}

// Init initialize the server router and set the controllerMachine needed to do the work
func (s *MachineServer) Init(machine controllerMachine, logger *httplog.Logger) {
	// Logger
	s.logger = logger
	s.configuration = &config.Config{}
	if err := s.configuration.ReadFromFile("configuration.yaml"); err != nil {
		s.logger.Error("Error", "error", err)
		panic("cannot read configuration file")
	}

	var err error
	s.ovenProgramManager, err = ovenprograms.NewOvenProgramManager(s.configuration.Server.OvenProgramFolder)

	if err != nil {
		s.logger.Error("Error", "error", err)
		panic("Something wrong")
	}
	s.machine = machine

	s.updateMachineFromConfig()
	s.ovenProgramWorker = ovenprograms.NewOvenProgramWorker(s.machine, *s.configuration, s.ovenProgramManager, s.logger)
	s.Router = chi.NewRouter()
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	s.Router.Use(httplog.RequestLogger(s.logger))
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	s.FileServer(s.Router.(*chi.Mux), s.configuration.Server.DistributionDirectory)
	s.Router.Route("/api", func(router chi.Router) {
		router.Route("/power-off", func(router chi.Router) {
			router.Post("/", s.powerOff)
		})
		router.Route("/processes", func(processRouter chi.Router) {
			processRouter.Route("/set-power-one-minute", func(r chi.Router) {
				r.Post("/", s.setPowerOneMinute)
			})
			processRouter.Route("/open-air", func(r chi.Router) {
				r.Post("/", s.openAir)
			})
			processRouter.Route("/close-air", func(r chi.Router) {
				r.Post("/", s.closeAir)
			})
			processRouter.Route("/get-temperature", func(r chi.Router) {
				r.Get("/", s.getTemperature)
			})
			processRouter.Route("/get-temperatures-process", func(r chi.Router) {
				r.Get("/", s.getTemperaturesProcess)
			})
			processRouter.Route("/is-working", func(r chi.Router) {
				r.Get("/", s.isWorking)
			})
			processRouter.Route("/test-ramp", func(r chi.Router) {
				r.Post("/", s.testRamp)
			})
			processRouter.Route("/get-actual-process-data", func(r chi.Router) {
				r.Get("/", s.getAllDataActualWork)
			})
			processRouter.Route("/start-process/{programName}", func(r chi.Router) {
				r.Post("/", s.startProgram)
			})
			processRouter.Route("/stop-process", func(r chi.Router) {
				r.Post("/", s.stopProgram)
			})

		})
		router.Route("/configuration", func(configRouter chi.Router) {
			configRouter.Route("/programs", func(r chi.Router) {
				r.Get("/", s.getPrograms)
				r.Post("/", s.addUpdateProgram)
				r.Get("/{programName}", s.getProgram)
				r.Delete("/{programName}", s.deleteProgram)
			})
			configRouter.Route("/oven-config", func(r chi.Router) {
				r.Get("/", s.getConfig)
				r.Post("/", s.updateConfig)
			})
			configRouter.Route("/move-runs-usb", func(r chi.Router) {
				r.Post("/", s.moveAllRunsToUsb)
			})

		})
	})

	s.initialized = true
}

func (s *MachineServer) updateMachineFromConfig() {
	s.machine.InitConfig(*s.configuration)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// FileServer is serving static files.
func (s MachineServer) FileServer(router *chi.Mux, root string) {
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		if _, err := os.Stat(filepath.Join(root, r.URL.Path)); os.IsNotExist(err) {
			s.logger.DebugContext(r.Context(), "not exist", "uri", r.URL.Path)
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func NewMachineServer(options ...func(*MachineServer)) *MachineServer {
	machineServer := &MachineServer{}
	for _, o := range options {
		o(machineServer)
	}
	return machineServer
}

func (s *MachineServer) powerOff(w http.ResponseWriter, r *http.Request) {
	if s.ovenProgramWorker.IsWorking() {
		s.ovenProgramWorker.RequestStopProgram()
		time.Sleep(time.Second)
	}
	err := exec.Command("shutdown", "-h", "now").Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: err.Error()})
	}

}
