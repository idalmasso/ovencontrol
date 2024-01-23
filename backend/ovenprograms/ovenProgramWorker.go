package ovenprograms

import (
	"encoding/csv"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/idalmasso/ovencontrol/backend/config"
)

type Oven interface {
	GetTemperature() (float64, error)
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
	ticker                             *time.Ticker
	isWorking                          bool
	kpRamp, kiRamp, kdRamp             float64
	kpMaintain, kiMaintain, kdMaintain float64
	stepTime, stepSave                 float64
	TargetTemperature                  float64
	SavedRunFolder                     string
	runName                            string
	endRequest                         bool
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

func (d *OvenProgramWorker) RequestStopProgram() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.endRequest = true
}

func (d OvenProgramWorker) shouldStopProgram() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.endRequest
}
func (d *OvenProgramWorker) startedProgram() error {
	f, err := os.OpenFile(filepath.Join(d.SavedRunFolder, "work.txt"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := csv.NewWriter(f)
	err = encoder.Write([]string{d.programName, ""})
	encoder.Flush()
	return err
}

func (d *OvenProgramWorker) endedProgram() {
	d.Save()
	d.programName = ""
	os.Remove(filepath.Join(d.SavedRunFolder, "work.txt"))
	d.oven.SetPercentual(0)
	d.oven.EndProgram()
}
func (d *OvenProgramWorker) changedStepPoint(s StepPoint) error {
	f, err := os.OpenFile(filepath.Join(d.SavedRunFolder, "work.txt"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := csv.NewWriter(f)
	err = encoder.Write([]string{d.programName, s.SegmentName})
	encoder.Flush()
	return err
}
func (d *OvenProgramWorker) StartOvenProgram(program OvenProgram) {
	d.mu.Lock()
	if d.isWorking {
		return
	}
	d.isWorking = true
	d.endRequest = false
	d.mu.Unlock()
	d.programName = program.Name
	d.runName = time.Now().Format("2006-01-02T15-04-05") + "-" + program.Name
	d.programHistory = make([]ProgramDataPoint, 0)
	d.lastPointsToBeWritten = 0
	d.startedProgram()
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
		defer func() {
			d.endedProgram()
		}()
		firstPoint := program.Points[0]
		d.changedStepPoint(firstPoint)
		temperature, err := d.oven.GetTemperature()
		if err != nil {
			return
		}
		if firstPoint.Temperature > temperature {
			d.doRamp(firstPoint, true)
		} else if firstPoint.Temperature == temperature {
			d.maintainTemperature(firstPoint)
		} else {
			d.doRamp(firstPoint, false)
		}
		if d.shouldStopProgram() {
			return
		}
		lastTemp := firstPoint.Temperature
		for _, s := range program.Points[1:] {
			d.changedStepPoint(s)
			if s.Temperature > lastTemp {
				d.doRamp(s, true)
			} else if s.Temperature == lastTemp {
				d.maintainTemperature(s)
			} else {
				d.doRamp(s, false)
			}
			lastTemp = s.Temperature
			if d.shouldStopProgram() {
				d.Save()
				d.programName = ""
				return
			}
		}
	}(program)
}

func (d *OvenProgramWorker) doRamp(s StepPoint, isUpRamp bool) error {
	var err error

	d.TargetTemperature, err = d.oven.GetTemperature()
	if err != nil {
		return err
	}
	integral, previousError, derivative, temperatureVariance := 0.0, 0.0, 0.0, 0.0
	desiredVariance := (s.Temperature - d.TargetTemperature) / s.TimeSeconds()
	ovenTemperature := d.TargetTemperature
	timeSave := 0.0
	lastNow := time.Now()
	step, newTemperature := 0.0, 0.0
	d.ticker = time.NewTicker(time.Duration(d.stepTime) * time.Second)
	defer d.ticker.Stop()
	for now := range d.ticker.C {
		if d.shouldStopProgram() {
			break
		}
		if ((ovenTemperature >= s.Temperature) && isUpRamp) || ((ovenTemperature <= s.Temperature) && !isUpRamp) {
			break
		}
		step = (now.Sub(lastNow)).Seconds()
		lastNow = now
		d.timeSeconds += step
		timeSave += step
		newTemperature, err = d.oven.GetTemperature()
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
		actualPercentual := d.kpRamp*errorValue + d.kiRamp*integral + d.kdRamp*derivative
		actualPercentual = min(actualPercentual, 1)
		actualPercentual = max(actualPercentual, 0)
		d.oven.SetPercentual(actualPercentual)
		previousError = errorValue
		d.programHistory = append(d.programHistory, createDataPoint(d.programName, s.SegmentName, d.timeSeconds, d.TargetTemperature, newTemperature, actualPercentual))
		d.lastPointsToBeWritten++
		if timeSave > d.stepSave {
			d.Save()
			d.lastPointsToBeWritten = 0
			timeSave = 0
		}
	}
	d.Save()
	return nil
}
func (d *OvenProgramWorker) maintainTemperature(s StepPoint) error {
	var err error
	d.TargetTemperature, err = d.oven.GetTemperature()
	if err != nil {
		return err
	}
	integral, previousError, derivative := 0.0, 0.0, 0.0
	d.TargetTemperature = s.Temperature
	first := true
	ovenTemperature := 0.0
	timeSave := 0.0
	totalTime := 0.0
	lastNow := time.Now()
	step := 0.0
	d.ticker = time.NewTicker(time.Duration(d.stepTime) * time.Second)
	defer d.ticker.Stop()
	for now := range d.ticker.C {
		if d.shouldStopProgram() {
			break
		}
		if totalTime >= s.TimeSeconds() {
			break
		}
		totalTime += step
		step = (now.Sub(lastNow)).Seconds()
		lastNow = now
		d.timeSeconds += step
		timeSave += step
		ovenTemperature, err = d.oven.GetTemperature()
		if err != nil {
			return err
		}
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
		d.programHistory = append(d.programHistory, createDataPoint(d.programName, s.SegmentName, d.timeSeconds, d.TargetTemperature, ovenTemperature, actualPercentual))
		d.lastPointsToBeWritten++
		if timeSave > d.stepSave {
			d.Save()
			d.lastPointsToBeWritten = 0
			timeSave = 0
		}
	}
	d.Save()
	return nil
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

func NewOvenProgramWorker(oven Oven, config config.Config, ovenProgramManager OvenProgramManager) *OvenProgramWorker {
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
	if _, err := os.Stat(filepath.Join(o.SavedRunFolder, "work.txt")); err != nil {
		if !os.IsNotExist(err) {
			//maybe we need to restart the program!
			f, err := os.OpenFile(filepath.Join(o.SavedRunFolder, "work.txt"), os.O_RDONLY, 0644)
			if err != nil {
				slog.Debug("Cannot open work file", "err", err)
				o.endedProgram()
			}
			defer f.Close()
			rec, err := csv.NewReader(f).Read()
			if err != nil {
				slog.Debug("Cannot read work file", "err", err)
				o.endedProgram()
			}
			program, ok := ovenProgramManager.Programs()[rec[0]]
			if !ok {
				slog.Debug("Cannot find program")
				o.endedProgram()
			}
			newProgram := OvenProgram{Name: program.Name, AirCloseAtDegrees: program.AirCloseAtDegrees}
			found := false

			if found {
				o.StartOvenProgram(newProgram)
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
