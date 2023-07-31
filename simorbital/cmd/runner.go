package main

import (
	"fmt"

	"github.com/Flokey82/genideas/simorbital"
)

func main() {
	elements := simorbital.OrbitalElements{
		SemiMajorAxis:     7000000,
		Eccentricity:      0.1,
		Inclination:       0.1,
		LongitudeOfNode:   0.2,
		ArgumentOfPerigee: 0.3,
		MeanAnomaly:       0.4,
	}

	// Calculate position and velocity vectors
	radius := simorbital.CalculateRadius(elements.SemiMajorAxis, elements.Eccentricity, elements.MeanAnomaly)
	trueAnomaly := simorbital.CalculateTrueAnomaly(simorbital.CalculateEccentricAnomaly(elements.MeanAnomaly, elements.Eccentricity), elements.Eccentricity)
	x, y, z := simorbital.CalculatePosition(radius, elements.Inclination, elements.LongitudeOfNode, elements.ArgumentOfPerigee, trueAnomaly)
	vx, vy, vz := simorbital.CalculateVelocity(elements.SemiMajorAxis, elements.Eccentricity, trueAnomaly, elements.Inclination, elements.LongitudeOfNode, elements.ArgumentOfPerigee)

	// Print results
	fmt.Printf("Position: (%f, %f, %f)\n", x, y, z)
	fmt.Printf("Velocity: (%f, %f, %f)\n", vx, vy, vz)
}
