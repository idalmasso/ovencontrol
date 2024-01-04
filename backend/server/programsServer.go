package server

import (
	"encoding/json"
	"net/http"

	"github.com/idalmasso/ovencontrol/backend/ovenprograms"
)

func (s *MachineServer) getPrograms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs())
}

func (s *MachineServer) addUpdateProgram(w http.ResponseWriter, r *http.Request) {
	var program ovenprograms.OvenProgram
	if err := json.NewDecoder(r.Body).Decode(&program); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	if err := s.ovenProgramManager.SaveProgram(program); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs)
}
