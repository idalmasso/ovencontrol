package server

import (
	"encoding/json"
	"net/http"
)

type temperatureReader interface {
	GetTemperature() float32
}

func (s *MachineServer) getTemperature(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ Temperature float32 }{Temperature: s.machine.GetTemperature()})
}
