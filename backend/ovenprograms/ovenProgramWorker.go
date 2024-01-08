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

func (d OvenProgramWorker) GetRunningProgram() string {
	return d.programName
}
func (d OvenProgramWorker) GetTimeSeconds() float64 {
	return d.timeSeconds
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
	d.runName = time.Now().Format("2006-01-02T15-04-05") + "-" + program.Name
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
	lastNow := time.Now()
	step := 0.0
	for ovenTemperature < s.Temperature {
		step = time.Since(lastNow).Seconds()
		lastNow = time.Now()
		d.timeSeconds += step
		timeSave += step
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * step
		d.TargetTemperature += expectedVariance
		if d.TargetTemperature > s.Temperature {
			d.TargetTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*step
		if step != 0 {
			derivative = (errorValue - previousError) / step
		}
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
			timeSave = 0
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
	lastNow := time.Now()
	step := 0.0
	for ovenTemperature < s.Temperature {
		step = time.Since(lastNow).Seconds()
		d.timeSeconds += step
		timeSave += step
		errorValue := s.Temperature - ovenTemperature
		if first {
			first = false
			previousError = errorValue
		}
		integral = integral + errorValue*step
		if step != 0 {
			derivative = (errorValue - previousError) / step
		}
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
			timeSave = 0
		}
		lastNow = time.Now()
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
	lastNow := time.Now()
	step := 0.0
	for d.oven.GetTemperature() > s.Temperature {
		step = time.Since(lastNow).Seconds()
		d.timeSeconds += step
		timeSave += step
		newTemperature := d.oven.GetTemperature()
		temperatureVariance = newTemperature - ovenTemperature
		ovenTemperature = newTemperature
		expectedVariance := desiredVariance * step
		d.TargetTemperature += expectedVariance
		if d.TargetTemperature > s.Temperature {
			d.TargetTemperature = s.Temperature
		}
		errorValue := expectedVariance - temperatureVariance

		integral = integral + errorValue*step
		if step != 0 {
			derivative = (errorValue - previousError) / step
		}
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
			timeSave = 0
		}
		lastNow = time.Now()
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

func (d *OvenProgramWorker) GetAllDataActualWork(step int) ProgramDataPointArray {
	if step == 1 {
		return d.programHistory
	} else {
		programHistoryLn := len(d.programHistory)
		res := make([]ProgramDataPoint, programHistoryLn/10+1)

		for i := 0; i < programHistoryLn; i += step {
			res[i/step] = d.programHistory[i]
		}
		return res
	}
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
