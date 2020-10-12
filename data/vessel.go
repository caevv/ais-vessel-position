package data

import (
	"math"
)

type Position struct {
	Imo       int     `json:"IMO"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type Distance struct {
	Kilometer     float64
	NauticalMiles float64
	StatuteMiles  float64
}

func CalculateDistance(positions []*Position) Distance {
	var distance Distance

	for i := 0; i < len(positions); i++ {
		if i == len(positions)-1 {
			return distance
		}

		d := calculate(positions[i].Latitude, positions[i].Longitude, positions[i+1].Latitude, positions[i+1].Longitude)
		distance.Kilometer += d.Kilometer
		distance.NauticalMiles += d.NauticalMiles
		distance.StatuteMiles += d.StatuteMiles
	}

	return distance
}

// This routine calculates the distance between two points (given the latitude/longitude of those points), by GeoDataSource (TM) products.
func calculate(lat1 float64, lng1 float64, lat2 float64, lng2 float64) Distance {
	const PI float64 = 3.141592653589793

	radLat1 := PI * lat1 / 180
	radLat2 := PI * lat2 / 180

	theta := lng1 - lng2
	radTheta := PI * theta / 180

	dist := math.Sin(radLat1)*math.Sin(radLat2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Cos(radTheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	return Distance{
		Kilometer:     dist * 1.609344,
		NauticalMiles: dist * 0.8684,
		StatuteMiles:  dist,
	}
}
