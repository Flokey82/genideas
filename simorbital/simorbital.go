package simorbital

import (
	"math"
)

// Gravitational constant (units: m^3 kg^-1 s^-2)
const G = 6.67430e-11

// OrbitalElements represents the Keplerian elements of an orbit.
type OrbitalElements struct {
	SemiMajorAxis     float64 // Semi-major axis (units: meters)
	Eccentricity      float64 // Eccentricity
	Inclination       float64 // Inclination (units: radians)
	LongitudeOfNode   float64 // Longitude of the ascending node (units: radians)
	ArgumentOfPerigee float64 // Argument of perigee (units: radians)
	MeanAnomaly       float64 // Mean anomaly (units: radians)
}

// CalculateMeanMotion calculates the mean motion (mean angular velocity) of an orbit.
func CalculateMeanMotion(semiMajorAxis float64) float64 {
	return math.Sqrt(G / semiMajorAxis)
}

// CalculateEccentricAnomaly calculates the eccentric anomaly from the mean anomaly and eccentricity.
func CalculateEccentricAnomaly(meanAnomaly, eccentricity float64) float64 {
	var (
		eccentricAnomaly = meanAnomaly
		delta            = 1.0
	)

	for delta > 1e-12 {
		eccentricAnomalyNew := eccentricAnomaly - (eccentricAnomaly-eccentricity*math.Sin(eccentricAnomaly)-meanAnomaly)/(1.0-eccentricity*math.Cos(eccentricAnomaly))
		delta = math.Abs(eccentricAnomalyNew - eccentricAnomaly)
		eccentricAnomaly = eccentricAnomalyNew
	}

	return eccentricAnomaly
}

// CalculateTrueAnomaly calculates the true anomaly from the eccentric anomaly and eccentricity.
func CalculateTrueAnomaly(eccentricAnomaly, eccentricity float64) float64 {
	return 2 * math.Atan(math.Sqrt((1+eccentricity)/(1-eccentricity))*math.Tan(eccentricAnomaly/2))
}

// CalculateRadius calculates the distance from the center of attraction (e.g., Earth) to the orbiting body.
func CalculateRadius(semiMajorAxis, eccentricity, trueAnomaly float64) float64 {
	return semiMajorAxis * (1 - math.Pow(eccentricity, 2)) / (1 + eccentricity*math.Cos(trueAnomaly))
}

// CalculatePosition calculates the position vector (x, y, z) in the orbital plane.
func CalculatePosition(radius, inclination, longitudeOfNode, argumentOfPerigee, trueAnomaly float64) (x, y, z float64) {
	xOrb := radius * (math.Cos(longitudeOfNode)*math.Cos(trueAnomaly+argumentOfPerigee) - math.Sin(longitudeOfNode)*math.Sin(trueAnomaly+argumentOfPerigee)*math.Cos(inclination))
	yOrb := radius * (math.Sin(longitudeOfNode)*math.Cos(trueAnomaly+argumentOfPerigee) + math.Cos(longitudeOfNode)*math.Sin(trueAnomaly+argumentOfPerigee)*math.Cos(inclination))
	zOrb := radius * (math.Sin(trueAnomaly+argumentOfPerigee) * math.Sin(inclination))
	return xOrb, yOrb, zOrb
}

// CalculateVelocity calculates the velocity vector (vx, vy, vz) in the orbital plane.
func CalculateVelocity(semiMajorAxis, eccentricity, trueAnomaly, inclination, longitudeOfNode, argumentOfPerigee float64) (vx, vy, vz float64) {
	n := CalculateMeanMotion(semiMajorAxis)
	radius := CalculateRadius(semiMajorAxis, eccentricity, trueAnomaly)
	xOrb, yOrb, _ := CalculatePosition(radius, inclination, longitudeOfNode, argumentOfPerigee, trueAnomaly)
	velocityMagnitude := math.Sqrt(G*semiMajorAxis*(1-math.Pow(eccentricity, 2))) / radius
	vxOrb := velocityMagnitude * (-math.Sin(longitudeOfNode)*math.Cos(trueAnomaly+argumentOfPerigee) - math.Cos(longitudeOfNode)*math.Sin(trueAnomaly+argumentOfPerigee)*math.Cos(inclination))
	vyOrb := velocityMagnitude * (math.Cos(longitudeOfNode)*math.Cos(trueAnomaly+argumentOfPerigee) - math.Sin(longitudeOfNode)*math.Sin(trueAnomaly+argumentOfPerigee)*math.Cos(inclination))
	vzOrb := velocityMagnitude * (math.Sin(trueAnomaly+argumentOfPerigee) * math.Sin(inclination))
	vx = vxOrb - n*yOrb
	vy = vyOrb + n*xOrb
	vz = vzOrb
	return vx, vy, vz
}
