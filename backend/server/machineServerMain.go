package server

import (
	"log/slog"
	"net/http"
	"os"
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
}

// PiServer
type MachineServer struct {
	ovenProgramManager ovenprograms.OvenProgramManager
	configuration      *config.Config
	initialized        bool
	Router             chi.Router
	machine            controllerMachine
	ovenProgramWorker  *ovenprograms.OvenProgramWorker
}

// ListenAndServe is the main server procedure that only wraps http.ListenAndServe
func (s *MachineServer) ListenAndServe() {

	if !s.initialized {
		panic("Server not initialized")
	}

	if err := http.ListenAndServe(":"+strconv.Itoa(s.configuration.Server.Port), s.Router); err != nil {
		panic("Cannot listen on server: " + err.Error())
	}
}

// Init initialize the server router and set the controllerMachine needed to do the work
func (s *MachineServer) Init(machine controllerMachine) {
	// Logger
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
	s.configuration = &config.Config{}
	if err := s.configuration.ReadFromFile("configuration.yaml"); err != nil {
		panic("cannot read configuration file")
	}
	var err error
	s.ovenProgramManager, err = ovenprograms.NewOvenProgramManager(s.configuration.Server.OvenProgramFolder)

	if err != nil {
		slog.Error("Error", err)
		panic("Something wrong")
	}
	s.machine = machine

	s.updateMachineFromConfig()
	s.ovenProgramWorker = ovenprograms.NewOvenProgramWorker(s.machine, *s.configuration)
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
	s.Router.Use(httplog.RequestLogger(logger))
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	FileServer(s.Router.(*chi.Mux), s.configuration.Server.DistributionDirectory)
	s.Router.Route("/api", func(router chi.Router) {
		router.Route("/processes", func(processRouter chi.Router) {
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
				r.Get("/", s.startProgram)
			})
		})
		router.Route("/configuration", func(configRouter chi.Router) {
			configRouter.Route("/programs", func(r chi.Router) {
				r.Get("/", s.getPrograms)
				r.Post("/", s.addUpdateProgram)
			})
			//configRouter.Get("/", s.getConfig)
			//configRouter.Put("/", s.updateConfig)
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
