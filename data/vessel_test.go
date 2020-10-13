package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormula(t *testing.T) {
	assert.Equal(t,
		Distance{
			Kilometer:     48.70598703891164,
			NauticalMiles: 26.281689399277507,
			StatuteMiles:  30.264497235464656,
		},
		formula(38.643582, 15.12622, 38.554115, 14.57754),
	)
}

func TestCalculateDistance(t *testing.T) {
	var position []*Position
	position1 := &Position{
		Latitude:  1,
		Longitude: 1,
	}
	position2 := &Position{
		Latitude:  2,
		Longitude: 2,
	}
	position3 := &Position{
		Latitude:  10,
		Longitude: 10,
	}

	assert.Equal(t,
		Distance{
			Kilometer: 1411.2138654256103,
			NauticalMiles: 761.4892283660918,
			StatuteMiles:  876.8876420613678,
		},
		CalculateDistance(append(position, position1, position2, position3)),
	)
}
