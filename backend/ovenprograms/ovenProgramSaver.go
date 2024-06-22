package ovenprograms

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type OvenProgramSaver struct {
	usbPath, usbSaveFolderName string
	savedRunFolder             string
}

func NewOvenProgramSaver(usbPath, usbSaveFolderName, savedRunFolder string) *OvenProgramSaver {
	return &OvenProgramSaver{usbPath: usbPath, usbSaveFolderName: usbSaveFolderName, savedRunFolder: savedRunFolder}
}

func (s OvenProgramSaver) MoveAllRuns() error {
	if _, err := os.Stat(s.usbPath); os.IsNotExist(err) {
		return fmt.Errorf("cannot find usb drive %w", err)
	}
	if _, err := os.Stat(s.savedRunFolder); os.IsNotExist(err) {
		return fmt.Errorf("cannot find saved run folder %w", err)
	}
	usbFilePath := s.usbPath
	//look for the usb in the mount folder (normally in raspy it is automounted in /media/pi/xxxxx, so I just search for first folder in this folder)
	if directories, err := os.ReadDir(s.usbPath); err != nil {
		return fmt.Errorf("error while selecting subfolders: %w", err)
	} else {
		found := false
		for _, d := range directories {
			if d.IsDir() {
				usbFilePath = filepath.Join(s.usbPath, d.Name())
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("error: cannot find subfolder in " + s.usbPath + " directory, unexpected")
		}
	}
	if _, err := os.Stat(filepath.Join(usbFilePath, s.usbSaveFolderName)); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(filepath.Join(usbFilePath, s.usbSaveFolderName), os.ModePerm); err != nil {
				return err
			}
		}
	}
	err := filepath.Walk(s.savedRunFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {

			err := moveFile(path, filepath.Join(usbFilePath, s.usbSaveFolderName, filepath.Base(path)))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}
