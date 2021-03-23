package sqlstore

import "math"

func haversineDistanceBetweenTwoPointsInMeters(firstLon float64, firstLat float64, secondLon float64, secondLat float64) float64 {
	R := 6371000.0
	p := 0.017453292519943295
	latOneRad := firstLon * p
	latTwoRad := secondLon * p
	deltaLat := (secondLat - firstLat) * p
	deltaLon := (secondLon - firstLon) * p

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latOneRad)*math.Cos(latTwoRad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
