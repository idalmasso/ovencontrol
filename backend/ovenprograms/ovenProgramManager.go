package ovenprograms

import (
	"os"
	"path/filepath"
)

type OvenProgramManager struct {
	folder   string
	programs map[string]OvenProgram
}

func (d OvenProgramManager) Programs() map[string]OvenProgram {
	return d.programs
}
func NewOvenProgramManager(programsFolder string) (OvenProgramManager, error) {
	ovenProgramManager := OvenProgramManager{folder: programsFolder, programs: make(map[string]OvenProgram)}
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
func (o OvenProgramManager) readAllPrograms() error {
	err := filepath.Walk(o.folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			prg, err := ReadOvenProgram(path)
			if err != nil {
				return err
			}
			o.programs[prg.Name] = prg
		}
		return nil
	})

	return err
}

func (o OvenProgramManager) SaveProgram(ovenProgram OvenProgram) error {
	err := ovenProgram.SaveToFile(o.folder)
	if err != nil {
		return err
	}
	o.programs[ovenProgram.Name] = ovenProgram
	return nil
}

func (o OvenProgramManager) DeleteProgram(ovenProgramName string) error {

	if _, ok := o.programs[ovenProgramName]; ok {
		err := o.programs[ovenProgramName].DeleteOvenProgramFile(o.folder)
		delete(o.programs, ovenProgramName)
		return err
	}

	return nil
}
