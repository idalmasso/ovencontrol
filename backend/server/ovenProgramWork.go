package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *MachineServer) getAllDataActualWork(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	step := 1
	var err error
	if r.URL.Query().Get("step") != "" {
		step, err = strconv.Atoi(r.URL.Query().Get("step"))
		if err != nil || step < 1 {
			step = 1
		}
	}
	json.NewEncoder(w).Encode(s.ovenProgramWorker.GetAllDataActualWork(step))
}
func (s *MachineServer) isWorking(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		IsWorking   bool   `json:"is-working"`
		ProgramName string `json:"program-name"`
	}{
		IsWorking:   s.ovenProgramWorker.IsWorking(),
		ProgramName: s.ovenProgramWorker.GetRunningProgram(),
	})
}
