package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type OvenProgram struct {
	Name              string
	Points            []StepPoint
	AirCloseAtDegrees float64
}

type StepPoint struct {
	Temperature float64
	TimeMinutes float64
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
