package server

import (
	"encoding/json"
	"net/http"

	"github.com/idalmasso/ovencontrol/backend/config"
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

	stepPoint := config.StepPoint{Temperature: temperature, TimeMinutes: timeMinutes}
	program := config.OvenProgram{Name: "TestProgram", Points: []config.StepPoint{stepPoint}}
	s.ovenProgramWorker.followOvenProgram(program)
	return true
}
func (s *MachineServer) isWorking(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ IsWorking bool }{IsWorking: s.ovenProgramWorker.isWorking})
}
