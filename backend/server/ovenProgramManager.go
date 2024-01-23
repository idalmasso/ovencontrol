package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idalmasso/ovencontrol/backend/ovenprograms"
)

func (s *MachineServer) getPrograms(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getPrograms called")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs())
}

func (s *MachineServer) addUpdateProgram(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("addUpdateProgram called")
	var program ovenprograms.OvenProgram
	if err := json.NewDecoder(r.Body).Decode(&program); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "addUpdateProgram error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	if err := s.ovenProgramManager.SaveProgram(program); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "addUpdateProgram error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs()[program.Name])
}

func (s *MachineServer) getProgram(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("getProgram called")
	programName := chi.URLParam(r, "programName")
	if program, ok := s.ovenProgramManager.Programs()[programName]; ok {

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(program)
		return
	} else {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "addUpdateProgram error", slog.String("error", "program not found"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Program not found"})
		return
	}

}

func (s *MachineServer) deleteProgram(w http.ResponseWriter, r *http.Request) {
	programName := chi.URLParam(r, "programName")
	if err := s.ovenProgramManager.DeleteProgram(programName); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "addUpdateProgram error", slog.String("error", "program not found"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Error string }{Error: "Program not found"})
	}
	w.WriteHeader(http.StatusOK)

}
