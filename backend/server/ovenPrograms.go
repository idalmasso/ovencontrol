package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idalmasso/ovencontrol/backend/ovenprograms"
)

func (s *MachineServer) testRamp(w http.ResponseWriter, r *http.Request) {
	if ok := s.TryStartTestRamp(s.configuration.Server.TestRampTemperature, s.configuration.Server.TestRampTimeMinutes); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Machine is working"})
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func (s *MachineServer) TryStartTestRamp(temperature, timeMinutes float64) bool {

	stepPoint := ovenprograms.StepPoint{Temperature: temperature, TimeMinutes: timeMinutes}
	program := ovenprograms.OvenProgram{Name: "TestProgram", Points: []ovenprograms.StepPoint{stepPoint}}
	s.ovenProgramWorker.StartOvenProgram(program)
	return true
}

func (s *MachineServer) startProgram(w http.ResponseWriter, r *http.Request) {
	if s.ovenProgramWorker.IsWorking() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Machine is working"})
		return
	}
	programName := chi.URLParam(r, "programName")
	if program, ok := s.ovenProgramManager.Programs()[programName]; ok {
		s.ovenProgramWorker.StartOvenProgram(program)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Program not found"})
		return
	}

}
