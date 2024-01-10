package ovenprograms

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type OvenProgram struct {
	Name              string      `json:"name"`
	IconColor         string      `json:"icon-color"`
	Points            []StepPoint `json:"points"`
	AirCloseAtDegrees float64     `json:"air-closed-at-degrees,string"`
}

type StepPoint struct {
	Temperature float64 `json:"temperature,string"`
	TimeMinutes float64 `json:"time-minutes,string"`
}

func (s StepPoint) TimeSeconds() float64 {
	return s.TimeMinutes * 60
}

func (p OvenProgram) SaveToFile(folderName string) error {
	file, err := os.Create(filepath.Join(folderName, p.Name) + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(p)
	return err
}

func (p *OvenProgram) ReadFromFile(fileName string) error {
	if !strings.HasSuffix(fileName, ".json") {
		fileName = fileName + ".json"
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(p)
	return err
}

func ReadOvenProgram(filename string) (OvenProgram, error) {
	p := OvenProgram{}
	err := p.ReadFromFile(filename)
	return p, err
}
