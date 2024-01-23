package server

import (
	"encoding/json"
	"net/http"
)

type temperatureReader interface {
	GetTemperature() float64
}

func (s *MachineServer) getTemperature(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getTemperature called")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Temperature float64 `json:"oven-temperature"`
	}{Temperature: s.machine.GetTemperature()})
}

func (s *MachineServer) getTemperaturesProcess(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getTemperaturesProcess called")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Temperature         float64 `json:"oven-temperature"`
		ExpectedTemperature float64 `json:"expected-temperature"`
		TimeSeconds         float64 `json:"time-seconds"`
	}{Temperature: s.machine.GetTemperature(), ExpectedTemperature: s.ovenProgramWorker.GetTargetTemperature(), TimeSeconds: s.ovenProgramWorker.GetTimeSeconds()})
}
