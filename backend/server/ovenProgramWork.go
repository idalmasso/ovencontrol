package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/idalmasso/ovencontrol/backend/ovenprograms"
)

func (s *MachineServer) getAllDataActualWork(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getAllDataActualWork called")
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
	s.logger.Debug("isWorking called")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		IsWorking   bool   `json:"is-working"`
		ProgramName string `json:"program-name"`
	}{
		IsWorking:   s.ovenProgramWorker.IsWorking(),
		ProgramName: s.ovenProgramWorker.GetRunningProgram(),
	})
}

func (s *MachineServer) stopProgram(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("stopProgram called")
	if s.ovenProgramWorker.IsWorking() {
		s.ovenProgramWorker.RequestStopProgram()
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Oven not in running mode"})
		return
	}

}

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

func (s *MachineServer) startProgram(w http.ResponseWriter, r *http.Request) {
	if s.ovenProgramWorker.IsWorking() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Machine is working"})
		return
	}
	programName := chi.URLParam(r, "programName")
	if program, ok := s.ovenProgramManager.Programs()[programName]; ok {
		s.ovenProgramWorker.StartOvenProgram(program)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Program not found"})
		return
	}

}

func (s *MachineServer) moveAllRunsToUsb(w http.ResponseWriter, r *http.Request) {
	saver := ovenprograms.NewOvenProgramSaver(s.configuration.Controller.UsbPath,
		s.configuration.Controller.UsbSaveFolderName,
		s.configuration.Controller.SavedRunFolder)
	if err := saver.MoveAllRuns(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: err.Error()})
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
