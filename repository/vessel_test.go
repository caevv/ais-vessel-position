package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/caevv/ais-vessel-position/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestVesselRepository_Positions(t *testing.T) {
	// prepare
	filePath := "./mockfile/"
	file := []string{"a.json", "b.json"}
	appFs := afero.NewOsFs()
	err := appFs.Mkdir(filePath, 0755)
	assert.NoError(t, err)

	movementTime, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.NoError(t, err)
	writeFile(t, appFs, filePath, file[0], movementTime.Add(time.Hour)) // This will test the ordering

	movementTime2, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.NoError(t, err)
	writeFile(t, appFs, filePath, file[1], movementTime2)

	// test
	r := New(filePath, file)
	actualPositions, err := r.Positions(1)
	assert.NoError(t, err)

	assert.Equal(
		t,
		[]*data.Position{
			{
				Imo:              1,
				Latitude:         1,
				Longitude:        1,
				MovementDateTime: movementTime,
			},
			{
				Imo:              1,
				Latitude:         1,
				Longitude:        1,
				MovementDateTime: movementTime2.Add(time.Hour),
			},
		},
		actualPositions,
	)

	err = appFs.RemoveAll(filePath)
	assert.NoError(t, err)
}

func writeFile(t *testing.T, appFs afero.Fs, filePath string, fileName string, movementTime time.Time) {
	jsonString, err := json.Marshal([]*data.Position{
		{
			Imo:              1,
			Latitude:         1,
			Longitude:        1,
			MovementDateTime: movementTime,
		},
	})
	assert.NoError(t, err)

	err = afero.WriteFile(appFs, filePath+fileName, jsonString, 0644)
	assert.NoError(t, err)
}
