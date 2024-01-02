package server

import (
	"encoding/json"
	"net/http"
)

type temperatureReader interface {
	GetTemperature() float64
	GetTemperatureExpected() float64
}

func (s *MachineServer) getTemperature(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ Temperature float64 }{Temperature: s.machine.GetTemperature()})
}

func (s *MachineServer) getTemperaturesProcess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ Temperature, ExpectedTemperature float64 }{Temperature: s.machine.GetTemperature(), ExpectedTemperature: s.machine.GetTemperatureExpected()})
}
