package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		DistributionDirectory string  `yaml:"distributionDirectory" json:"distribution-directory"`
		Port                  int     `yaml:"port" json:"port,string"`
		OvenProgramFolder     string  `yaml:"ovenProgramFolder" json:"oven-program-folder"`
		TestRampTemperature   float64 `yaml:"testRampTemperature" json:"test-ramp-temperature"`
		TestRampTimeMinutes   float64 `yaml:"testRampTimeMinutes" json:"test-ramp-time-minutes"`
	} `yaml:"server" json:"server"`
	Oven struct {
		Length                float64   `yaml:"length" json:"length,string"`
		Height                float64   `yaml:"height" json:"height,string"`
		Width                 float64   `yaml:"width" json:"width,string"`
		InsultationWidths     []float64 `yaml:"insulationWidths" json:"insulation-widths"`
		ThermalConductivities []float64 `yaml:"thermalConductivities" json:"thermal-conductivities"`
		ThermalCapacity       float64   `yaml:"thermalCapacity" json:"thermal-capacity,string"`
		Weight                float64   `yaml:"weight" json:"weight,string"`
		MaxPower              float64   `yaml:"maxPower" json:"max-power,string"`
	} `yaml:"oven" json:"oven"`
	Controller struct {
		KpRamp         float64 `yaml:"kpRamp" json:"kp-ramp,string"`
		KiRamp         float64 `yaml:"kiRamp" json:"ki-ramp,string"`
		KdRamp         float64 `yaml:"kdRamp" json:"kd-ramp,string"`
		KpMaintain     float64 `yaml:"kpMaintain" json:"kp-maintain,string"`
		KiMaintain     float64 `yaml:"kiMaintain" json:"ki-maintain,string"`
		KdMaintain     float64 `yaml:"kdMaintain" json:"kd-maintain,string"`
		StepTime       float64 `yaml:"stepTime" json:"step-time,string"`
		StepSave       float64 `yaml:"saveTime" json:"save-time,string"`
		SavedRunFolder string  `yaml:"savedRunFolder" json:"saved-run-folder"`
	} `yaml:"controller" json:"controller"`
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
