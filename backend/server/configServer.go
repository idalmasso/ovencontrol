package server

import (
	"encoding/json"
	"net/http"

	"github.com/idalmasso/ovencontrol/backend/config"
)

func (s *MachineServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	var config config.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	s.configuration = &config
	if err := s.configuration.SaveToFile("configuration.yaml"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}

	s.updateMachineFromConfig()

	w.WriteHeader(http.StatusOK)

}

func (s *MachineServer) getConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.configuration)
}
