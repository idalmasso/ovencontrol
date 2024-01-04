package ovenprograms

import (
	"fmt"
	"math"
	"time"
)

type ProgramDataPoint struct {
	ProgramName        string  `json:"program-name"`
	SecondsFromStart   float64 `json:"seconds-from-start"`
	DateTime           string  `json:"datetime"`
	DesiredTemperature float64 `json:"desired-temperature"`
	OvenTemperature    float64 `json:"oven-temperature"`
	OvenPercentage     float64 `json:"oven-percentage"`
}

type ProgramDataPointArray []ProgramDataPoint

func (history ProgramDataPointArray) toStrings() [][]string {
	res := make([][]string, len(history))
	for idx, d := range history {
		s := make([]string, 6)
		s[0] = d.ProgramName
		s[1] = fmt.Sprintf("%.1f", d.SecondsFromStart)
		s[2] = d.DateTime
		s[3] = fmt.Sprintf("%.1f", d.DesiredTemperature)
		s[4] = fmt.Sprintf("%.1f", d.OvenTemperature)
		s[5] = fmt.Sprintf("%.1f", d.OvenPercentage)
		res[idx] = s
	}
	return res
}
func programHistoryHeaders() []string {
	s := make([]string, 6)
	s[0] = "Program name"
	s[1] = "Seconds from start"
	s[2] = "Datetime"
	s[3] = "Target temperature"
	s[4] = "Oven temperature"
	s[5] = "Power percentage"
	return s
}

func createDataPoint(programName string, secondsFromStart float64, desiredTemperature float64, ovenTemperature float64, ovenPercentage float64) ProgramDataPoint {
	now := time.Now()
	return ProgramDataPoint{ProgramName: programName,
		SecondsFromStart:   math.Round(secondsFromStart*100) / 100,
		DesiredTemperature: math.Round(desiredTemperature*100) / 100,
		OvenTemperature:    math.Round(ovenTemperature*100) / 100,
		OvenPercentage:     math.Round(ovenPercentage*10000) / 10000,
		DateTime:           now.Format("2006-01-02T15:04:05"),
	}
}
