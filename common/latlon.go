package common

import (
	"github.com/goccy/go-json"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"math"
)

type SafeLatitudeLongitude []int

func convertToPrecisionGrid(latitude, longitude float64, precision float64) (int, int) {
	// Constants for conversion
	const latDegreeTo5km = 1 / (111.0 / 5.0)
	const earthCircumference = 40075.0 // in km

	// Convert latitude to a discrete value
	latIndex := int(latitude / latDegreeTo5km)

	// Calculate the distance per degree of longitude at the given latitude
	longitudeDegreeDistance := math.Cos(latitude*math.Pi/180) * earthCircumference / 360.0
	longDegreeToPrecision := 1 / (longitudeDegreeDistance / precision)

	// Convert longitude to a discrete value
	longIndex := int(longitude / longDegreeToPrecision)

	return latIndex, longIndex
}

func (s SafeLatitudeLongitude) Set(lat, lon, precision float64) {

	l1, l2 := convertToPrecisionGrid(lat, lon, precision)
	s[0] = l1
	s[1] = l2

}

func CommitLocation(secret kyber.Scalar, location []byte) kyber.Point {
	// Hash the location to a scalar
	suite := edwards25519.NewBlakeSHA256Ed25519()

	locationHash := suite.Hash().Sum(location)
	locationScalar := suite.Scalar().SetBytes(locationHash)

	// Generate the commitment as C = secret * G + locationScalar * G
	G := suite.Point().Base() // Generator point
	secretPart := suite.Point().Mul(secret, G)
	locationPart := suite.Point().Mul(locationScalar, G)
	commitment := suite.Point().Add(secretPart, locationPart)
	return commitment
}

func (s SafeLatitudeLongitude) Bytes() []byte {
	joinedLatLon, _ := json.Marshal(s)
	return joinedLatLon
}
