package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/idalmasso/ovencontrol/backend/config"
)

func (s *MachineServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	var config config.Config
	s.logger.DebugContext(r.Context(), "updateConfig called")
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "updateConfig error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	s.configuration = &config
	if err := s.configuration.SaveToFile("configuration.yaml"); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "updateConfig error", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}

	s.updateMachineFromConfig()

	w.WriteHeader(http.StatusOK)

}

func (s *MachineServer) getConfig(w http.ResponseWriter, r *http.Request) {
	s.logger.DebugContext(r.Context(), "getConfig called")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.configuration)
}
