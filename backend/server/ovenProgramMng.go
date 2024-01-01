package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type ovenProgramManager struct {
	Folder   string
	Programs map[string]config.OvenProgram
}

func NewOvenProgramManager(programsFolder string) (ovenProgramManager, error) {
	ovenProgramManager := ovenProgramManager{Folder: programsFolder, Programs: make(map[string]config.OvenProgram)}
	if _, err := os.Stat(programsFolder); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(programsFolder, os.ModePerm); err != nil {
				return ovenProgramManager, err
			}
		} else {
			return ovenProgramManager, err
		}
	}
	err := ovenProgramManager.readAllPrograms()
	return ovenProgramManager, err

}
func (o ovenProgramManager) readAllPrograms() error {
	err := filepath.Walk(o.Folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			prg, err := config.ReadOvenProgram(path)
			if err != nil {
				return err
			}
			o.Programs[prg.Name] = prg
		}
		return nil
	})

	return err
}

func (o ovenProgramManager) saveProgram(ovenProgram config.OvenProgram) error {
	err := ovenProgram.SaveToFile(o.Folder)
	if err != nil {
		return err
	}
	o.Programs[ovenProgram.Name] = ovenProgram
	return nil
}

func (s *MachineServer) getPrograms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs)
}

func (s *MachineServer) addUpdateProgram(w http.ResponseWriter, r *http.Request) {
	var program config.OvenProgram
	if err := json.NewDecoder(r.Body).Decode(&program); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	if err := s.ovenProgramManager.saveProgram(program); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{ Err error }{Err: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.ovenProgramManager.Programs)
}
