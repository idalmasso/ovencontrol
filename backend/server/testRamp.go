package server

import (
	"encoding/json"
	"net/http"

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
