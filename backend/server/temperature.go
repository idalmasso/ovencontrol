package server

import (
	"encoding/json"
	"net/http"
)

type temperatureReader interface {
	GetTemperature() (float64, error)
}

func (s *MachineServer) getTemperature(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getTemperature called")
	temperature, err := s.machine.GetTemperature()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{Error: err.Error()})

		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Temperature float64 `json:"oven-temperature"`
	}{Temperature: temperature})
}

func (s *MachineServer) getTemperaturesProcess(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getTemperaturesProcess called")
	temperature, err := s.machine.GetTemperature()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{Error: err.Error()})

		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Temperature         float64 `json:"oven-temperature"`
		ExpectedTemperature float64 `json:"expected-temperature"`
		TimeSeconds         float64 `json:"time-seconds"`
	}{Temperature: temperature, ExpectedTemperature: s.ovenProgramWorker.GetTargetTemperature(), TimeSeconds: s.ovenProgramWorker.GetTimeSeconds()})
}
