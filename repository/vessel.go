package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/caevv/ais-vessel-position/data"
)

type Vessel interface {
	Positions(imo int) ([]*data.Position, error)
}

type VesselRepository struct{
	filesPath string
}

func New(path string) Vessel {
	return &VesselRepository{
		filesPath: path,
	}
}

func (r VesselRepository) Positions(imo int) ([]*data.Position, error) {
	var (
		wg                   sync.WaitGroup
		positions            []*data.Position
		position202007291231 *data.Position
		position202007291931 *data.Position
		position202007292331 *data.Position
		errs                 []error
	)

	wg.Add(1)
	go func() {
		var err error
		position202007291231, err = r.readFile("CombinedPositionsData_20200729_202007291231_CombinedPositionsData.json", imo)
		if err != nil {
			errs = append(errs, err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var err error
		position202007291931, err = r.readFile("CombinedPositionsData_20200729_202007291931_CombinedPositionsData.json", imo)
		if err != nil {
			errs = append(errs, err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var err error
		position202007292331, err = r.readFile("CombinedPositionsData_20200729_202007292331_CombinedPositionsData.json", imo)
		if err != nil {
			errs = append(errs, err)
		}
		wg.Done()
	}()

	wg.Wait()

	if len(errs) > 0 {
		var errorMessage string
		for _, err := range errs {
			errorMessage += " " + err.Error()
		}
		return nil, errors.New(errorMessage)
	}

	return append(positions, position202007291231, position202007291931, position202007292331), nil
}

func (r VesselRepository) readFile(fileName string, imo int) (*data.Position, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s%s", r.filesPath, fileName))
	if err != nil {
		return nil, err
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var positions []data.Position

	err = json.Unmarshal(byteValue, &positions)
	if err != nil {
		return nil, err
	}

	for _, position := range positions {
		if position.Imo == imo {
			return &position, nil
		}
	}

	return nil, errors.New("not found")
}
