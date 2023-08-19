package genvegetation

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

type Edaphology struct {
}

func NewEdaphology() Edaphology {
	return Edaphology{}
}

// CalculateAngles calculates the angles (steepness) for every pixel. This is done by calculating the dot product of the normal
// vector and the z-vector on each point of the terrain.
// :param normal_map: 3D-List of all normal vectors.
// :return: angles: List all angles on the terrain.
func (e *Edaphology) CalculateAngles(normalMap [][]vectors.Vec3) [][]float64 {
	angles := make([][]float64, len(normalMap))
	for y := range normalMap {
		angles[y] = make([]float64, len(normalMap[y]))
		for x := range normalMap[y] {
			standardNormal := vectors.Vec3{X: 0.0, Y: 0.0, Z: 1.0}
			dotProduct := standardNormal.Dot(normalMap[y][x])
			standardNormalLength := standardNormal.Len()
			normalLength := normalMap[y][x].Len()
			angle := dotProduct / (standardNormalLength * normalLength)
			angle = math.Acos(angle)
			angle = radToDeg(angle)
			angles[y][x] = angle
		}
	}
	return angles
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// CalculateSoilDepth calculates the soil depth at every point of the terrain. It uses the previously calculated angles (steepness).
// :param map: Object of Map class. Used to get the max soil depth.
// :param size: Integer. Size of the terrain. Used to initialize the image.
// :param angles: 3D-List of all normal vectors.
// :return: angles: List all angles on the terrain.
func (e *Edaphology) CalculateSoilDepth(m *Map, size int, angles [][]float64) [][]float64 {
	soilDepths := make([][]float64, size)
	for y := range angles {
		soilDepths[y] = make([]float64, size)
		for x := range angles[y] {
			depth := (1 - (angles[y][x] / 90)) * m.max_soil_depth
			soilDepths[y][x] = depth
		}
	}
	return soilDepths
}

/*
import math
import numpy as np

from intern.src.python.Data.image import Image


class Edaphology:
    """
    The Edaphology class is used to calculate the soil depth for every pixel. It uses the previously calculated
    normal map. The dot product of every normal vector and the z-vector (0,0,1) will be calculated. The result are
    the angles of the terrain at each point. With the help of the angles (steepness) the soil depth will be calculated.
    This depends on the maximal soil depth stated in the map.
    """

    @staticmethod
    def calculate_angles(normal_map):
        """
        Calculates the angles (steepness) for every pixel. This is done by calculating the dot product of the normal
        vector and the z-vector on each point of the terrain.
        :param normal_map: 3D-List of all normal vectors.
        :return: angles: List all angles on the terrain.
        """
        angles = []
        for rows in normal_map:
            row = []
            for normal in rows:
                standard_normal = [0.0, 0.0, 1.0]
                dot_product = np.dot(standard_normal, normal)
                standard_normal_length = np.linalg.norm(standard_normal)
                normal_length = np.linalg.norm(normal)
                angle = dot_product / (standard_normal_length * normal_length)
                angle = math.acos(angle)
                angle = math.degrees(angle)
                row.append(angle)
            angles.append(row)
        return angles

    @staticmethod
    def calculate_soil_depth(map, size, angles):
        """
        It calculates the soil depth at every point of the terrain. It uses the previously calculated angles (steepness).
        :param map: Object of Map class. Used to get the max soil depth.
        :param size: Integer. Size of the terrain. Used to initialize the image.
        :param angles: 3D-List of all normal vectors.
        :return: angles: List all angles on the terrain.
        """
        soil_depths = Image(size=size)
        x = 0
        y = 0
        for rows in angles:
            print("Calculating edaphology: Row: " + str(y))
            for angle in rows:
                depth = (1 - angle/90) * map.max_soil_depth
                soil_depths.image[y][x] = depth
                x += 1
            y += 1
            x = 0
        return soil_depths
*/
