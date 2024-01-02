package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hardware struct {
		GearRatio float64 `yaml:"gearRatio" json:"gearRatio,string"`
	} `yaml:"hardware" json:"hardware"`
	Server struct {
		DistributionDirectory string  `yaml:"distributionDirectory" json:"distributionDirectory"`
		Port                  int     `yaml:"port" json:"port,string"`
		OvenProgramFolder     string  `yaml:"ovenProgramFolder" json:"ovenProgramFolder"`
		TestRampTemperature   float64 `yaml:"testRampTemperature" json:"testRampTemperature"`
		TestRampTimeMinutes   float64 `yaml:"testRampTimeMinutes" json:"testRampTimeMinutes"`
	} `yaml:"server" json:"server"`
	Oven struct {
		Length                float64   `yaml:"length" json:"length"`
		Height                float64   `yaml:"height" json:"height"`
		Width                 float64   `yaml:"width" json:"width"`
		InsultationWidths     []float64 `yaml:"insulationWidths" json:"insulationWidths"`
		ThermalConductivities []float64 `yaml:"thermalConductivities" json:"thermalConductivities"`
		ThermalCapacity       float64   `yaml:"thermalCapacity" json:"thermalCapacity"`
		Weight                float64   `yaml:"weight" json:"weight"`
		MaxPower              float64   `yaml:"maxPower" json:"maxPower"`
	} `yaml:"oven" json:"oven"`
}

func (c *Config) ReadFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	return
}

func (c *Config) SaveToFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(c)
	return
}
