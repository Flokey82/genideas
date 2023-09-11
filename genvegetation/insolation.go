package genvegetation

import (
	"fmt"
	"math"
	"strconv"
)

type Insolation struct {
	controller      *Controller
	insolationImage *Image
}

func NewInsolation(controller *Controller) *Insolation {
	return &Insolation{
		controller:      controller,
		insolationImage: NewImage(len(controller.image_height_map), 0, DtUint16),
	}
}

// CalculateRawInsolation calculates if a given point on the terrain receives light at a given daylight hour.
// Atmospheric absorption etc. is not considered during this calculation.
// :param sun: Object of the sun class. Used for calculating from which direction the sun shines.
// :param x, y: Integer of the x and y position of the point on the terrain.
// :param map_size: Integer of the size of the terrain.
// :param pixel_size: Float. The size a pixel represents of the real terrain.
// :param heightmap_max_height: Integer. Maximal height of the terrain.
// :param height_conversion: Float. Conversion value of the height of the heightmap to calculate the real height.
// :return: amount of raw insolation, 0 if no light reaches the pixel.
func (i *Insolation) CalculateRawInsolation(sun *Sun, x, y, mapSize int, pixelSize, heightmapMaxHeight, heightConversion float64) float64 {
	// Transform the polar coordinates of the sun to cartesian coordinates.
	uvw := sun.ConvertToUvCoordinates()
	xStep := uvw.X
	yStep := uvw.Y
	zStep := uvw.Z

	// We start at the target pixel position (translated to the real world position).
	xStartWorldPos := float64(x) * pixelSize
	yStartWorldPos := float64(y) * pixelSize
	zStartWorldPos := i.controller.image_height_map[y][x] * heightConversion

	// Initialize the current position for walking along the direction vector.
	xRealWorldPos := xStartWorldPos
	yRealWorldPos := yStartWorldPos
	zRealWorldPos := zStartWorldPos

	t := 0.0 // used for the vector equation
	sunBeamReachesPixel := true

	// Calculate the length of one side of the map.
	mapXYBoundary := float64(mapSize) * pixelSize // boundary of the map

	// Walk along the sun direction vector until the sun beam leaves the map boundary or the height of the terrain
	// is higher than the current height of the sun beam.
	for 0 <= xRealWorldPos && xRealWorldPos < mapXYBoundary && 0 <= yRealWorldPos && yRealWorldPos < mapXYBoundary && zRealWorldPos < heightmapMaxHeight {
		var tStep float64
		if xStep == 0.0 && yStep == 0.0 {
			break // sun stands in zenith so every pixel will receive light
		}

		// This if statement decides how far the sun direction vector will be followed until a new
		// pixel in the pixel space will be reached. Only then a new height can be compared.
		// This accelerates the algorithm.

		// Get the remainder distance to the next pixel.
		//
		// NOTE: Depending on the sign of the vector components, we need to calculate the remainder differently
		// as we are either moving in a positive or a negatve direction, so we either need the distance to the
		// next pixel (in terms of increasing coordinates) or the distance to the previous pixel (in terms of
		// decreasing coordinates).
		//
		// Negative xStep and yStep:
		//
		//  pixel size
		// |__________|
		// |__________|
		// |          |
		// |          |remX = mod(x, pixel_size)
		// | previous |____|
		// | pixel    |    |
		// -------------------------
		//            |    |     | |
		//            |    |     | | remY = mod(y, pixel_size)
		//            |    P-----|---
		//            |          |
		//            ------------
		//
		// Positive xStep and yStep:
		//
		// ____________
		// |          |
		// |          |
		// |    P-----|---
		// |    |     | | remY = pixel_size - mod(y, pixel_size)
		// ------------------------
		//      |_____| Next pixel|
		//      |     |           |
		//       remX = pixel_size - mod(x, pixel_size)
		//            |___________|
		//
		remX := math.Mod(xRealWorldPos, pixelSize)
		if xStep >= 0 || remX == 0 {
			remX = pixelSize - remX
		}
		remY := math.Mod(yRealWorldPos, pixelSize)
		if yStep >= 0 || remY == 0 {
			remY = pixelSize - remY
		}

		// Calculate the step size for the current direction vector to reach the next pixel.
		if remX <= remY && xStep != 0 {
			tStep = remX / xStep
		} else {
			tStep = remY / yStep
		}

		// Adjust our current position along the direction vector.
		t += math.Abs(tStep)

		xRealWorldPos = xStartWorldPos + t*xStep
		yRealWorldPos = yStartWorldPos + t*yStep
		zRealWorldPos = zStartWorldPos + t*zStep

		// Calculate the pixel position in the pixel space.
		xPixelPos := int(xRealWorldPos / pixelSize)
		yPixelPos := int(yRealWorldPos / pixelSize)

		// Check if the pixel position is outside of the map.
		if xPixelPos < 0 || yPixelPos < 0 || xPixelPos > mapSize-1 || yPixelPos > mapSize-1 {
			break // sun beam leaves the map boundary
		}

		// Check if the height of the terrain is higher than the current height of the sun beam.
		terrainHeight := i.controller.image_height_map[yPixelPos][xPixelPos] * heightConversion
		if terrainHeight > zRealWorldPos {
			sunBeamReachesPixel = false
			break // something blocks the light from the sun for that pixel
		}
	}
	if sunBeamReachesPixel {
		return SOLAR_CONSTANT_K_CALORIES_PER_HOUR // * (map.pixel_size ** 2)
	}
	return 0.0
}

// nppad is a go implementation of numpy.pad. It pads the image with the given number of pixels.
// :param image: 2D-List of the image.
// :param pad_width: Integer. Number of pixels that will be added to the image.
// :param mode: String. Defines how the image will be padded. 'edge' pads the image with the edge values of the image.
// :return: padded_image: 2D-List of the padded image.
func nppad(image [][]float64, padWidth int, mode string) [][]float64 {
	paddedImage := make([][]float64, len(image)+2*padWidth)
	for y := range paddedImage {
		paddedImage[y] = make([]float64, len(image)+2*padWidth)
	}
	for y := range image {
		for x := range image[y] {
			paddedImage[y+padWidth][x+padWidth] = image[y][x]
		}
	}
	return paddedImage
}

// AddReflectionInsolation adds the energy of the neighbours of a pixel and calculates the average. A fraction of this number will
// be added to the currently observed pixel.
// :param reflection_coefficient: Float. Fraction of the average energy of the neighbourspixels that the pixel
// will receive.
func (i *Insolation) AddReflectionInsolation(reflectionCoefficient float64) {
	paddedInsolationImage := nppad(i.insolationImage.image, 1, "edge")
	for y := 1; y < len(paddedInsolationImage)-1; y++ {
		for x := 1; x < len(paddedInsolationImage[y])-1; x++ {
			neighborInsolationSum := 0.0
			neighborInsolationSum += paddedInsolationImage[y][x+1]
			neighborInsolationSum += paddedInsolationImage[y+1][x+1]
			neighborInsolationSum += paddedInsolationImage[y+1][x]
			neighborInsolationSum += paddedInsolationImage[y+1][x-1]
			neighborInsolationSum += paddedInsolationImage[y][x-1]
			neighborInsolationSum += paddedInsolationImage[y-1][x-1]
			neighborInsolationSum += paddedInsolationImage[y-1][x]
			neighborInsolationSum += paddedInsolationImage[y-1][x+1]
			addedReflectionInsolation := neighborInsolationSum / 8 * reflectionCoefficient
			i.insolationImage.image[y-1][x-1] += addedReflectionInsolation
		}
	}
}

// CalculateActualInsolation calculates the actual energy of each pixel based on the previously calculated raw energy. The atmosphere and
// reflection reduce the raw energy.
// :param map_name: String of the current map name.
// :param daylight_hours: Integer of the number of daylight hours.
// :param sun_start_elevation: Float of the start elevation of the sun.
// :param sun_start_azimuth: Float of the start azimuth of the sun.
// :param sun_max_elevation: Float of the maximal sun elevation (noon).
// :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
// pixel will be reflected.
//
// :return: insolation_image: Image of the calculated actual energy of each pixel.
func (i *Insolation) CalculateActualInsolation(m *Map, daylightHours int, sunStartElevation, sunStartAzimuth, sunMaxElevation, reflectionCoefficient float64) *Image {
	i.CalculateInsolationForDaylightHours(m, daylightHours, sunStartElevation, sunStartAzimuth, sunMaxElevation)
	for y := 0; y < len(i.controller.image_height_map); y++ {
		for x := 0; x < len(i.controller.image_height_map[y]); x++ {
			pixelRawInsolation := i.insolationImage.image[y][x]
			cloudReflectionLoss := pixelRawInsolation * m.biom.cloud_reflection / 100
			atmosphericAbsorptionLoss := pixelRawInsolation * m.biom.atmospheric_absorption / 100
			atmosphericDiffusionLoss := pixelRawInsolation * m.biom.atmospheric_diffusion / 100
			soilId := i.controller.soil_ids_map[y][x]
			soil := i.controller.SearchSoil(soilId)
			albedo := soil.albedo
			i.insolationImage.image[y][x] = (pixelRawInsolation - cloudReflectionLoss - atmosphericAbsorptionLoss - atmosphericDiffusionLoss) * (1.0 - albedo)
		}
	}
	i.AddReflectionInsolation(reflectionCoefficient)
	return i.insolationImage
}

/*

   def calculate_insolation_for_daylight_hours(self, map, daylight_hours, sun_start_elevation, sun_start_azimuth,
                                               sun_max_elevation):
       """
       Calculates the sun position for every day light hours. At each hour the raw energy for each pixel wil be
       calculated.
       :param map_name: String of the current map name.
       :param daylight_hours: Integer of the number of daylight hours.
       :param sun_start_elevation: Float of the start elevation of the sun.
       :param sun_start_azimuth: Float of the start azimuth of the sun.
       :param sun_max_elevation: Float of the maximal sun elevation (noon).
       """
       assert daylight_hours > 0, "Daylight hours must be at least one!"

       # elevation
       if daylight_hours == 1:
           elevation_per_hour = 0
       elif daylight_hours == 2:
           elevation_per_hour = sun_max_elevation - sun_start_elevation
       else:
           if daylight_hours % 2 == 1:
               elevation_per_hour = ((sun_max_elevation - sun_start_elevation) / math.ceil(
                   (daylight_hours - 2) / 2))  # the sun shall rise to 90 (or less) degrees till noon and then fall again
           else:
               elevation_per_hour = (sun_max_elevation - sun_start_elevation) / (daylight_hours / 2)

       # azimuth
       if daylight_hours == 1:
           azimuth_per_hour = 0
       elif daylight_hours == 2:
           azimuth_per_hour = 180 - 2 * sun_start_azimuth
       else:
           azimuth_per_hour = 180 / (daylight_hours - 1)  # the sun shall wander 180 degrees

       sun = Sun(elevation=sun_start_elevation, azimuth=sun_start_azimuth)

       for hour in range(daylight_hours):
           print("############ Hour: " + str(hour + 1) + " ############")
           print("Sun polar coordinates: Azimuth: " + str(round(sun.azimuth, 1)) + "° Elevation: " + str(
               round(sun.elevation, 1)) + "°")
           max_height = np.amax(self.controller.image_height_map.image) * map.height_conversion
           for y in range(self.controller.image_height_map.size):
               print("Raw Insolation: Row:" + str(y))
               for x in range(self.controller.image_height_map.size):
                   self.calculate_raw_insolation(sun, x, y, self.controller.image_height_map.size, map.pixel_size,
                                                 max_height, map.height_conversion)

           if daylight_hours % 2 == 1:
               if hour < int((daylight_hours / 2)):
                   sun.elevation += elevation_per_hour
               else:
                   sun.elevation -= elevation_per_hour
           else:
               if hour != daylight_hours / 2 - 1:
                   if hour < int((daylight_hours / 2)):
                       sun.elevation += elevation_per_hour
                   else:
                       sun.elevation -= elevation_per_hour
           sun.azimuth += azimuth_per_hour
*/

// npamax is a go implementation of numpy.amax. It returns the maximum value of the given image.
func npamax(image [][]float64) float64 {
	max := 0.0
	for y := range image {
		for x := range image[y] {
			if image[y][x] > max {
				max = image[y][x]
			}
		}
	}
	return max
}

// CalculateInsolationForDaylightHours calculates the sun position for every day light hours. At each hour the raw energy for each pixel wil be
// calculated.
// :param map_name: String of the current map name.
// :param daylight_hours: Integer of the number of daylight hours.
// :param sun_start_elevation: Float of the start elevation of the sun.
// :param sun_start_azimuth: Float of the start azimuth of the sun.
// :param sun_max_elevation: Float of the maximal sun elevation (noon).
func (i *Insolation) CalculateInsolationForDaylightHours(m *Map, daylightHours int, sunStartElevation, sunStartAzimuth, sunMaxElevation float64) {
	if daylightHours < 1 {
		panic("Daylight hours must be at least one!")
	}

	// elevation
	var elevationPerHour float64
	if daylightHours == 1 {
		elevationPerHour = 0
	} else if daylightHours == 2 {
		elevationPerHour = sunMaxElevation - sunStartElevation
	} else if daylightHours > 2 {
		if daylightHours%2 == 1 {
			elevationPerHour = (sunMaxElevation - sunStartElevation) / math.Ceil(float64(daylightHours-2)/2) // the sun shall rise to 90 (or less) degrees till noon and then fall again
		} else {
			elevationPerHour = (sunMaxElevation - sunStartElevation) / (float64(daylightHours) / 2)
		}
	}

	// azimuth
	var azimuthPerHour float64
	if daylightHours == 1 {
		azimuthPerHour = 0
	} else if daylightHours == 2 {
		azimuthPerHour = 180 - 2*sunStartAzimuth
	} else if daylightHours > 2 {
		azimuthPerHour = 180 / float64(daylightHours-1) // the sun shall wander 180 degrees
	}
	sun := NewSun(sunStartElevation, sunStartAzimuth)
	for hour := 0; hour < daylightHours; hour++ {
		fmt.Println("############ Hour: " + strconv.Itoa(hour+1) + " ############")
		fmt.Println("Sun polar coordinates: Azimuth: " + strconv.FormatFloat(sun.azimuth, 'f', 1, 64) + "° Elevation: " + strconv.FormatFloat(sun.elevation, 'f', 1, 64) + "°")
		maxHeight := npamax(i.controller.image_height_map) * m.height_conversion
		for y := 0; y < len(i.controller.image_height_map); y++ {
			fmt.Println("Raw Insolation: Row:" + strconv.Itoa(y))
			for x := 0; x < len(i.controller.image_height_map[y]); x++ {
				i.insolationImage.image[y][x] += i.CalculateRawInsolation(sun, x, y, len(i.controller.image_height_map), m.pixel_size, maxHeight, m.height_conversion)
			}
		}
		if daylightHours%2 == 1 {
			if hour < daylightHours/2 {
				sun.elevation += elevationPerHour
			} else {
				sun.elevation -= elevationPerHour
			}
		} else {
			if hour != daylightHours/2-1 {
				if hour < daylightHours/2 {
					sun.elevation += elevationPerHour
				} else {
					sun.elevation -= elevationPerHour
				}
			}
		}
		sun.azimuth += azimuthPerHour
	}
}

/*

import math
import numpy as np

from intern.src.python.Core.constants import *
from intern.src.python.Data.image import Image
from intern.src.python.Data.sun import Sun


class Insolation:
    """
    The insolation class is used for calculating the amount of energy that each point on the terrain receives by
    sunlight during one day. First it is calculated which points receive light. Then the actual energy will be
    calculated at every point. This depends on the number of daylight hours a point receives light, the solar constant,
    atmospheric absorption, cloud reflection, atmospheric diffusion and ground reflection (albedo).
    At last the reflection of surrounding pixels will be added to the energy of a pixel.
    """

    def __init__(self, controller):
        self.controller = controller
        self.insolation_image = Image(size=self.controller.image_height_map.size)
*/
