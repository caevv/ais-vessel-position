package repository

import (
	"encoding/json"
	"fmt"
	"github.com/caevv/ais-vessel-position/pkg/aisvesselposition"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

const maxRoutines = 2

type Vessel interface {
	Positions(imo int) ([]*aisvesselposition.Position, error)
}

type VesselRepository struct {
	filesPath string
	filesName []string
}

func New(path string, filesName []string) Vessel {
	return &VesselRepository{
		filesPath: path,
		filesName: filesName,
	}
}

func (r VesselRepository) Positions(imo int) ([]*aisvesselposition.Position, error) {
	var (
		wg   sync.WaitGroup
		errs []error
	)

	sem := make(chan struct{}, maxRoutines)

	filesQty := len(r.filesName)
	positions := make([]*aisvesselposition.Position, filesQty)

	wg.Add(filesQty)

	for i, fileName := range r.filesName {
		sem <- struct{}{} // add to the goroutine and block if full

		go func(fileName string, i int) {
			defer func() { <-sem }() // decrement from goroutine

			position, err := r.readFile(fileName, imo)
			if err != nil {
				errs = append(errs, err)
			} else {
				positions[i] = position // ensure order
			}
			wg.Done()
		}(fileName, i)
	}

	wg.Wait()

	if len(errs) > 0 {
		var errorMessage string
		for _, err := range errs {
			errorMessage += " " + err.Error()
		}
		return nil, errors.New(errorMessage)
	}

	// Positions should be ordered, but we will assume worst case scenario in case json files were not ordered.
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].MovementDateTime.Before(positions[j].MovementDateTime)
	})

	return positions, nil
}

func (r VesselRepository) readFile(fileName string, imo int) (*aisvesselposition.Position, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s%s", r.filesPath, fileName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read json file")
	}

	var positions []aisvesselposition.Position

	err = json.Unmarshal(byteValue, &positions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize json")
	}

	for _, position := range positions {
		if position.Imo == imo {
			return &position, nil
		}
	}

	return nil, errors.New("not found")
}
