package ovenprograms

import (
	"fmt"
	"math"
	"time"
)

type ProgramDataPoint struct {
	ProgramName        string  `json:"program-name"`
	SegmentName        string  `json:"segment-name"`
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
		s := make([]string, 7)
		s[0] = d.ProgramName
		s[1] = d.SegmentName
		s[2] = fmt.Sprintf("%.1f", d.SecondsFromStart)
		s[3] = d.DateTime
		s[4] = fmt.Sprintf("%.1f", d.DesiredTemperature)
		s[5] = fmt.Sprintf("%.1f", d.OvenTemperature)
		s[6] = fmt.Sprintf("%.2f", d.OvenPercentage)
		res[idx] = s
	}
	return res
}
func programHistoryHeaders() []string {
	s := make([]string, 7)
	s[0] = "Program name"
	s[1] = "Segment name"
	s[2] = "Seconds from start"
	s[3] = "Datetime"
	s[4] = "Target temperature"
	s[5] = "Oven temperature"
	s[6] = "Power percentage"
	return s
}

func createDataPoint(programName string, segmentName string, secondsFromStart float64, desiredTemperature float64, ovenTemperature float64, ovenPercentage float64) ProgramDataPoint {
	now := time.Now()
	return ProgramDataPoint{ProgramName: programName,
		SegmentName:        segmentName,
		SecondsFromStart:   math.Round(secondsFromStart*100) / 100,
		DesiredTemperature: math.Round(desiredTemperature*100) / 100,
		OvenTemperature:    math.Round(ovenTemperature*100) / 100,
		OvenPercentage:     math.Round(ovenPercentage*10000) / 10000,
		DateTime:           now.Format("2006-01-02T15:04:05"),
	}
}
