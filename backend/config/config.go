package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hardware struct {
		MotorDegreePerStep float64 `yaml:"motorDegreePerStep" json:"motorDegreePerStep,string"`
		WaitForStep        int     `yaml:"waitForStep" json:"waitForStep,string"`
		GearRatio          float64 `yaml:"gearRatio" json:"gearRatio,string"`
	} `yaml:"hardware" json:"hardware"`
	Server struct {
		DistributionDirectory string `yaml:"distributionDirectory" json:"distributionDirectory"`
		Port                  int    `yaml:"port" json:"port,string"`
	} `yaml:"server" json:"server"`
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
