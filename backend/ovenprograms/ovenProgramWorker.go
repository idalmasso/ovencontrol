package ovenprograms

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type Oven interface {
	GetTemperature() float64
	GetPercentual() float64
	GetMaxPower() float64
	SetPercentual(float64)
	InitStartProgram()
	EndProgram()
}

type OvenProgramWorker struct {
	programName                        string
	timeSeconds                        float64
	oven                               Oven
	mu                                 *sync.RWMutex
	isWorking                          bool
	kpRamp, kiRamp, kdRamp             float64
	kpMaintain, kiMaintain, kdMaintain float64
	stepTime, stepSave                 float64
	TargetTemperature                  float64
	SavedRunFolder                     string
	runName                            string
	programHistory                     ProgramDataPointArray
	lastPointsToBeWritten              int
}

func (d OvenProgramWorker) GetTargetTemperature() float64 {
	return math.Round(d.TargetTemperature*100) / 100
}
func (d *OvenProgramWorker) IsWorking() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.isWorking
}

func (d *OvenProgramWorker) StartOvenProgram(program OvenProgram) {
	d.mu.Lock()
	if d.isWorking {
		return
	}
	d.isWorking = true
	d.mu.Unlock()
	d.programName = program.Name
	d.runName = program.Name + time.Now().Format("2006-01-02T15:04:05")
	d.programHistory = make([]ProgramDataPoint, 0)
	d.lastPointsToBeWritten = 0
	go func(program OvenProgram) {

		defer func() {
			d.mu.Lock()
			d.isWorking = false
			d.mu.Unlock()
		}()

		d.timeSeconds = 0
		if len(program.Points) == 0 {
			return
		}
		d.oven.InitStartProgram()
		d.writeHeader()
		defer func() { d.oven.EndProgram() }()
		firstPoint := program.Points[0]
		if firstPoint.Temperature > d.oven.GetTemperature() {
			d.doRampUp(firstPoint)
		} else if firstPoint.Temperature == d.oven.GetTemperature() {
			d.maintainTemperature(firstPoint)
		} else {
			d.doRampDown(firstPoint)
		}
		lastTemp := firstPoint.Temperature
		for _, s := range program.Points[1:] {
			if s.Temperature > lastTemp {
				d.doRampUp(s)
			} else if s.Temperature == lastTemp {
				d.maintainTemperature(s)
			} else {
				d.doRampDown(s)
			}
			lastTemp = s.Temperature
		}
		d.Save()
		d.programName = ""
	}(program)
}

func (d *OvenProgramWorker) doRampUp(s StepPoint) {
	d.TargetTemperature = d.oven.GetTemperature()
	integral, previousError, derivative, temperatureVariance := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.oven.GetTemperature()) / s.TimeSeconds()
	ovenTemperature := d.oven.GetTemperature()
	timeSave := 0.0
	for ovenTemperature < s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * d.stepTime
		d.TargetTemperature += expectedVariance
		if d.TargetTemperature > s.Temperature {
			d.TargetTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := d.oven.GetPercentual()
		actualPercentual += d.kpRamp*errorValue + d.kiRamp*integral + d.kdRamp*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		d.programHistory = append(d.programHistory, createDataPoint(d.programName, d.timeSeconds, d.TargetTemperature, newTemperature, actualPercentual))
		d.lastPointsToBeWritten++
		if timeSave > d.stepSave {
			d.Save()
			d.lastPointsToBeWritten = 0
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
	}
}
func (d *OvenProgramWorker) maintainTemperature(s StepPoint) {
	d.TargetTemperature = d.oven.GetTemperature()
	integral, previousError, derivative := 0.0, 0.0, 0.0
	d.TargetTemperature = s.Temperature
	timeSave := 0.0
	first := true
	ovenTemperature := d.oven.GetTemperature()
	for ovenTemperature < s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		errorValue := s.Temperature - ovenTemperature
		if first {
			first = false
			previousError = errorValue
		}
		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := d.kpMaintain*errorValue + d.kiMaintain*integral + d.kdMaintain*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		d.programHistory = append(d.programHistory, createDataPoint(d.programName, d.timeSeconds, d.TargetTemperature, ovenTemperature, actualPercentual))
		d.lastPointsToBeWritten++
		if timeSave > d.stepSave {
			d.Save()
			d.lastPointsToBeWritten = 0
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
		ovenTemperature = d.oven.GetTemperature()
	}
}
func (d *OvenProgramWorker) doRampDown(s StepPoint) {
	d.TargetTemperature = d.oven.GetTemperature()
	integral, previousError, derivative, temperatureVariance := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.oven.GetTemperature()) / s.TimeSeconds()
	ovenTemperature := d.oven.GetTemperature()
	timeSave := 0.0
	for d.oven.GetTemperature() > s.Temperature {
		d.timeSeconds += d.stepTime
		timeSave += d.stepTime
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * d.stepTime
		d.TargetTemperature += expectedVariance
		if d.TargetTemperature > s.Temperature {
			d.TargetTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*d.stepTime
		derivative = (errorValue - previousError) / d.stepTime
		actualPercentual := d.oven.GetPercentual()
		actualPercentual += d.kpRamp*errorValue + d.kiRamp*integral + d.kdRamp*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		d.programHistory = append(d.programHistory, createDataPoint(d.programName, d.timeSeconds, d.TargetTemperature, newTemperature, actualPercentual))
		d.lastPointsToBeWritten++
		if timeSave > d.stepSave {
			d.Save()
			d.lastPointsToBeWritten = 0
		}
		time.Sleep(time.Duration(d.stepTime) * time.Second)
	}
}
func (d *OvenProgramWorker) Save() error {
	f, err := os.OpenFile(filepath.Join(d.SavedRunFolder, d.runName+".txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := csv.NewWriter(f)
	err = encoder.WriteAll(d.programHistory[len(d.programHistory)-d.lastPointsToBeWritten:].toStrings())
	return err
}
func (d *OvenProgramWorker) writeHeader() error {
	f, err := os.OpenFile(filepath.Join(d.SavedRunFolder, d.runName+".txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := csv.NewWriter(f)
	err = encoder.Write(programHistoryHeaders())
	encoder.Flush()
	return err
}
func NewOvenProgramWorker(oven Oven, config config.Config) *OvenProgramWorker {
	o := OvenProgramWorker{oven: oven}
	o.mu = &sync.RWMutex{}
	o.isWorking = false
	o.stepTime = config.Controller.StepTime
	o.stepSave = config.Controller.StepSave
	o.kdMaintain = config.Controller.KdMaintain
	o.kiMaintain = config.Controller.KiMaintain
	o.kpMaintain = config.Controller.KpMaintain
	o.kdRamp = config.Controller.KdRamp
	o.kiRamp = config.Controller.KiRamp
	o.kpRamp = config.Controller.KpRamp
	o.SavedRunFolder = config.Controller.SavedRunFolder
	if _, err := os.Stat(o.SavedRunFolder); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(o.SavedRunFolder, os.ModePerm); err != nil {
				return nil
			}
		}
	}
	return &o
}

func (d *OvenProgramWorker) GetAllDataActualWork() ProgramDataPointArray {
	return d.programHistory
}

func (d *OvenProgramWorker) GetEndedRunList() ([]string, error) {
	listRun := make([]string, 0)
	err := filepath.Walk(d.SavedRunFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			listRun = append(listRun, info.Name())
		}
		return err
	})

	return listRun, err
}