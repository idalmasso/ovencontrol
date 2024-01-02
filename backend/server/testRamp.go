package server

import (
	"encoding/json"
	"net/http"
)

type machineRampTester interface {
	TryStartTestRamp(temperature, timeMinutes float64) bool
	IsWorking() bool
}

func (s *MachineServer) testRamp(w http.ResponseWriter, r *http.Request) {
	if ok := s.machine.TryStartTestRamp(s.configuration.Server.TestRampTemperature, s.configuration.Server.TestRampTimeMinutes); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Machine is working"})
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
