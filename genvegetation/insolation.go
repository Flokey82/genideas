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

/*

   def calculate_raw_insolation(self, sun, x, y, map_size, pixel_size, heightmap_max_height, height_conversion):
       """
       Calculates if a given point on the terrain receives light at a given daylight hour. Atmospheric absorption etc.
       is not considered during this calculation.
       :param sun: Object of the sun class. Used for calculating from which direction the sun shines.
       :param x, y: Integer of the x and y position of the point on the terrain.
       :param map_size: Integer of the size of the terrain.
       :param pixel_size: Float. The size a pixel represents of the real terrain.
       :param heightmap_max_height: Integer. Maximal height of the terrain.
       :param height_conversion: Float. Conversion value of the height of the heightmap to calculate the real height.
       """
       uvw = sun.convert_to_uv_coordinates()  # transforms the polar coordinates of the sun to cartesian coordinates
       x_step = uvw[0]
       y_step = uvw[1]
       z_step = uvw[2]
       x_start_world_pos = x * pixel_size
       y_start_world_pos = y * pixel_size
       z_start_world_pos = self.controller.image_height_map.image[y][x] * height_conversion
       x_real_world_pos = x_start_world_pos  # current x position for walking along the direction vector
       y_real_world_pos = y_start_world_pos  # current y position for walking along the direction vector
       z_real_world_pos = z_start_world_pos  # current z position for walking along the direction vector
       t = 0  # used for the vector equation
       sun_beam_reaches_pixel = True
       map_x_y_boundary = map_size * pixel_size  # boundary of the map

       while 0 <= x_real_world_pos < map_x_y_boundary and 0 <= y_real_world_pos < map_x_y_boundary \
               and z_real_world_pos < heightmap_max_height:

           # this if statement decides how far the sun direction vector will be followed until a new
           # pixel in the pixel space will be reached. Only then a new height can be compared. This accelerates the
           # algorithm.
           if x_step != 0.0 or y_step != 0.0:
               if x_real_world_pos % pixel_size >= y_real_world_pos % pixel_size and x_step != 0:
                   t_step = (pixel_size - (x_real_world_pos % pixel_size)) / x_step
               else:
                   t_step = (pixel_size - (y_real_world_pos % pixel_size)) / y_step
           else:
               break  # sun stands in zenith so every pixel will receive light

           t_step = abs(t_step)
           t += t_step
           x_real_world_pos = x_start_world_pos + t * x_step
           y_real_world_pos = y_start_world_pos + t * y_step
           z_real_world_pos = z_start_world_pos + t * z_step
           x_pixel_pos = int(x_real_world_pos / pixel_size)
           y_pixel_pos = int(y_real_world_pos / pixel_size)

           if x_pixel_pos < 0 or y_pixel_pos < 0 or x_pixel_pos > map_size - 1 or y_pixel_pos > map_size - 1:
               break  # sun beam leaves the map boundary

           terrain_height = self.controller.image_height_map.image[y_pixel_pos][x_pixel_pos] * height_conversion
           line_height = z_real_world_pos
           if terrain_height > line_height:
               sun_beam_reaches_pixel = False
               break  # something blocks the light from the sun for that pixel
       if sun_beam_reaches_pixel:
           self.insolation_image.image[y][x] += SOLAR_CONSTANT_K_CALORIES_PER_HOUR  # * (map.pixel_size ** 2)


*/

// CalculateRawInsolation calculates if a given point on the terrain receives light at a given daylight hour.
// Atmospheric absorption etc. is not considered during this calculation.
// :param sun: Object of the sun class. Used for calculating from which direction the sun shines.
// :param x, y: Integer of the x and y position of the point on the terrain.
// :param map_size: Integer of the size of the terrain.
// :param pixel_size: Float. The size a pixel represents of the real terrain.
// :param heightmap_max_height: Integer. Maximal height of the terrain.
// :param height_conversion: Float. Conversion value of the height of the heightmap to calculate the real height.
func (i *Insolation) CalculateRawInsolation(sun *Sun, x, y, mapSize int, pixelSize, heightmapMaxHeight, heightConversion float64) {
	uvw := sun.ConvertToUvCoordinates() // transforms the polar coordinates of the sun to cartesian coordinates
	xStep := uvw.X
	yStep := uvw.Y
	zStep := uvw.Z

	xStartWorldPos := float64(x) * pixelSize
	yStartWorldPos := float64(y) * pixelSize
	zStartWorldPos := i.controller.image_height_map[y][x] * heightConversion
	xRealWorldPos := xStartWorldPos // current x position for walking along the direction vector
	yRealWorldPos := yStartWorldPos // current y position for walking along the direction vector
	zRealWorldPos := zStartWorldPos // current z position for walking along the direction vector
	t := 0.0                        // used for the vector equation
	sunBeamReachesPixel := true
	mapXYBoundary := float64(mapSize) * pixelSize // boundary of the map

	for 0 <= xRealWorldPos && xRealWorldPos < mapXYBoundary && 0 <= yRealWorldPos && yRealWorldPos < mapXYBoundary && zRealWorldPos < heightmapMaxHeight {
		// this if statement decides how far the sun direction vector will be followed until a new
		// pixel in the pixel space will be reached. Only then a new height can be compared. This accelerates the
		// algorithm.
		var tStep float64
		if xStep != 0.0 || yStep != 0.0 {
			if math.Mod(xRealWorldPos, pixelSize) >= math.Mod(yRealWorldPos, pixelSize) && xStep != 0 {
				tStep = (pixelSize - math.Mod(xRealWorldPos, pixelSize)) / xStep
			} else {
				tStep = (pixelSize - math.Mod(yRealWorldPos, pixelSize)) / yStep
			}
		} else {
			break // sun stands in zenith so every pixel will receive light
		}

		tStep = math.Abs(tStep)
		t += tStep
		xRealWorldPos = xStartWorldPos + t*xStep
		yRealWorldPos = yStartWorldPos + t*yStep
		zRealWorldPos = zStartWorldPos + t*zStep
		xPixelPos := int(xRealWorldPos / pixelSize)
		yPixelPos := int(yRealWorldPos / pixelSize)

		if xPixelPos < 0 || yPixelPos < 0 || xPixelPos > mapSize-1 || yPixelPos > mapSize-1 {
			break // sun beam leaves the map boundary
		}

		terrainHeight := i.controller.image_height_map[yPixelPos][xPixelPos] * heightConversion
		lineHeight := zRealWorldPos
		if terrainHeight > lineHeight {
			sunBeamReachesPixel = false
			break // something blocks the light from the sun for that pixel
		}
	}
	if sunBeamReachesPixel {
		i.insolationImage.image[y][x] += SOLAR_CONSTANT_K_CALORIES_PER_HOUR // * (map.pixel_size ** 2)
	}
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

/*

	def add_reflection_insolation(self, reflection_coefficient):
        """
        Adds the energy of the neighbours of a pixel and calculates the average. A fraction of this number will
        be added to the currently observed pixel.
        :param reflection_coefficient: Float. Fraction of the average energy of the neighbourspixels that the pixel
            will receive.
        """
        padded_insolation_image = np.pad(self.insolation_image.image, 1, 'edge')
        for y in range(1, self.controller.image_height_map.size + 1):
            print("Reflection: Row: " + str(y))
            for x in range(1, self.controller.image_height_map.size + 1):
                neighbor_insolation_sum = 0
                neighbor_insolation_sum += padded_insolation_image[y][x + 1]
                neighbor_insolation_sum += padded_insolation_image[y + 1][x + 1]
                neighbor_insolation_sum += padded_insolation_image[y + 1][x]
                neighbor_insolation_sum += padded_insolation_image[y + 1][x - 1]
                neighbor_insolation_sum += padded_insolation_image[y][x - 1]
                neighbor_insolation_sum += padded_insolation_image[y - 1][x - 1]
                neighbor_insolation_sum += padded_insolation_image[y - 1][x]
                neighbor_insolation_sum += padded_insolation_image[y - 1][x + 1]
                added_reflection_insolation = neighbor_insolation_sum / 8 * reflection_coefficient
                self.insolation_image.image[y - 1][x - 1] += added_reflection_insolation

*/

// AddReflectionInsolation adds the energy of the neighbours of a pixel and calculates the average. A fraction of this number will
// be added to the currently observed pixel.
// :param reflection_coefficient: Float. Fraction of the average energy of the neighbourspixels that the pixel
//
//	will receive.
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

/*

   def calculate_actual_insolation(self, map, daylight_hours, sun_start_elevation, sun_start_azimuth,
                                   sun_max_elevation, reflection_coefficient):
       """
       Calculates the actual energy of each pixel based on the previously calculated raw energy. The atmosphere and
       reflection reduce the raw energy.
       :param map_name: String of the current map name.
       :param daylight_hours: Integer of the number of daylight hours.
       :param sun_start_elevation: Float of the start elevation of the sun.
       :param sun_start_azimuth: Float of the start azimuth of the sun.
       :param sun_max_elevation: Float of the maximal sun elevation (noon).
       :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
               pixel will be reflected.
       :return: insolation_image: Image of the calculated actual energy of each pixel.
       """
       self.calculate_insolation_for_daylight_hours(map, daylight_hours, sun_start_elevation, sun_start_azimuth,
                                                    sun_max_elevation)
       for y in range(self.controller.image_height_map.size):
           for x in range(self.controller.image_height_map.size):
               pixel_raw_insolation = self.insolation_image.image[y][x]
               cloud_reflection_loss = pixel_raw_insolation * map.biom.cloud_reflection / 100
               atmospheric_absorption_loss = pixel_raw_insolation * map.biom.atmospheric_absorption / 100
               atmospheric_diffusion_loss = pixel_raw_insolation * map.biom.atmospheric_diffusion / 100
               soil_id = self.controller.soil_ids_map.image[y][x]
               soil = self.controller.search_soil(soil_id)
               albedo = soil.albedo
               self.insolation_image.image[y][x] = (pixel_raw_insolation - cloud_reflection_loss -
                                                    atmospheric_absorption_loss -
                                                    atmospheric_diffusion_loss) * (1.0 - albedo)
       self.add_reflection_insolation(reflection_coefficient)
       return self.insolation_image

*/

// CalculateActualInsolation calculates the actual energy of each pixel based on the previously calculated raw energy. The atmosphere and
// reflection reduce the raw energy.
// :param map_name: String of the current map name.
// :param daylight_hours: Integer of the number of daylight hours.
// :param sun_start_elevation: Float of the start elevation of the sun.
// :param sun_start_azimuth: Float of the start azimuth of the sun.
// :param sun_max_elevation: Float of the maximal sun elevation (noon).
// :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
//
//	pixel will be reflected.
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
				i.CalculateRawInsolation(sun, x, y, len(i.controller.image_height_map), m.pixel_size, maxHeight, m.height_conversion)
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
