package server

import (
	"encoding/json"
	"net/http"
)

func (s *MachineServer) openAir(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("openAir called")
	err := s.machine.OpenAir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{Error: err.Error()})

		return
	}
	w.WriteHeader(http.StatusOK)

}

func (s *MachineServer) closeAir(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("closeAir called")
	err := s.machine.CloseAir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{Error: err.Error()})

		return
	}
	w.WriteHeader(http.StatusOK)

}
