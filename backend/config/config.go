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
		DistributionDirectory string `yaml:"distributionDirectory" json:"distributionDirectory"`
		Port                  int    `yaml:"port" json:"port,string"`
		OvenProgramFolder     string `yaml:"ovenProgramFolder" json:"ovenProgramFolder"`
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
