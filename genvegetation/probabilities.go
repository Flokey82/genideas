package genvegetation

import (
	"fmt"
)

type Probabilities struct {
	controller *Controller
	Map        *Map
	vegetation *Vegetation
}

func NewProbabilities(controller *Controller, vegetation *Vegetation, mapName string) *Probabilities {
	return &Probabilities{
		controller: controller,
		Map:        controller.maps[mapName],
		vegetation: vegetation,
	}
}

func (p *Probabilities) CalculateProbability(needed, available float64) float64 {
	if available <= needed {
		return available / needed
	} else if available <= needed*2 {
		return 1 - (available/needed - 1)
	} else {
		return 0
	}
}

// CalculateInsolationProbabilities calculates the insolation probabilities by comparing the available and needed insolation.
// :param map: Object of the map class.
// :param image_insolation_map: Result of the insolation calculation.
func (p *Probabilities) CalculateInsolationProbabilities() [][]float64 {
	insolationProbabilities := new2DArr(len(p.controller.image_height_map))
	for y := 0; y < len(p.controller.image_height_map); y++ {
		for x := 0; x < len(p.controller.image_height_map[y]); x++ {
			availableCalories := p.controller.image_insolation_map.image[y][x]
			neededCalories := p.vegetation.energy_demand
			insolationProbabilities[y][x] = p.CalculateProbability(neededCalories, availableCalories)
		}
	}
	return insolationProbabilities
}

// CalculateSoilDemandProbabilities calculates the soil demand probabilities by comparing the available and needed soil.
// :param map: Object of the map class.
// :param soil_ids_map: Map of the soil ids.
func (p *Probabilities) CalculateSoilDemandProbabilities() [][]float64 {
	soilDemandProbabilities := new2DArr(len(p.controller.image_height_map))
	for y := 0; y < len(p.controller.image_height_map); y++ {
		for x := 0; x < len(p.controller.image_height_map[y]); x++ {
			if p.vegetation.soil_demand.id == p.controller.soil_ids_map[y][x] {
				soilDemandProbabilities[y][x] = 1.0
			} else {
				soilDemandProbabilities[y][x] = 0.0
			}
		}
	}
	return soilDemandProbabilities
}

// CalculateSoilDepthProbabilities calculates the soil depth probabilities by comparing the available and needed soil depth.
// :param map: Object of the map class.
// :param edaphic_map: Object of the Edaphology class. Used to get the soil depth.
func (p *Probabilities) CalculateSoilDepthProbabilities() [][]float64 {
	soilDepthProbabilities := new2DArr(len(p.controller.image_height_map))
	for y := 0; y < len(p.controller.image_height_map); y++ {
		for x := 0; x < len(p.controller.image_height_map[y]); x++ {
			availableSoilDepth := p.controller.image_edaphic_map[y][x]
			neededSoilDepth := p.vegetation.soil_depth_demand
			if availableSoilDepth < neededSoilDepth {
				soilDepthProbabilities[y][x] = availableSoilDepth / neededSoilDepth
			} else {
				soilDepthProbabilities[y][x] = 1.0
			}
		}
	}
	return soilDepthProbabilities
}

// CalculateWaterDemandProbabilities calculates the water demand probabilities by comparing the available and needed water.
// :param map: Object of the map class.
// :param image_water_map: Result of the water calculation.
func (p *Probabilities) CalculateWaterDemandProbabilities() [][]float64 {
	waterDemandProbabilities := new2DArr(len(p.controller.image_height_map))
	for y := 0; y < len(p.controller.image_height_map); y++ {
		for x := 0; x < len(p.controller.image_height_map[y]); x++ {
			availableWater := p.controller.image_water_map[y][x]
			neededWater := p.vegetation.water_demand
			waterDemandProbabilities[y][x] = p.CalculateProbability(neededWater, availableWater)
		}
	}
	return waterDemandProbabilities
}

// CalculateProbabilities calculates the final probability by selecting the lowest probability at every pixel.
// :param map: Object of the map class.
// :param image_insolation_map: Result of the insolation calculation.
// :param image_water_map: Result of the water calculation.
// :param image_edaphic_map: Result of the edaphic calculation.
// :param soil_ids_map: Map of the soil ids.
// :return: final_probabilities: List of final probabilities at each pixel.
func (p *Probabilities) CalculateProbabilities() [][]float64 {
	allProbabilities := [][][]float64{
		p.CalculateInsolationProbabilities(),
		p.CalculateSoilDemandProbabilities(),
		p.CalculateSoilDepthProbabilities(),
		p.CalculateWaterDemandProbabilities(),
	}
	reasonsForNotGrowing := []int{0, 0, 0, 0}

	finalProbabilities := new2DArr(len(p.controller.image_height_map))
	for y := 0; y < len(p.controller.image_height_map); y++ {
		for x := 0; x < len(p.controller.image_height_map[y]); x++ {
			probability := 1.0
			for i := 0; i < len(allProbabilities); i++ {
				if allProbabilities[i][y][x] < probability {
					probability = allProbabilities[i][y][x]
					if probability == 0.0 {
						reasonsForNotGrowing[i] += 1
					}
				}
			}
			finalProbabilities[y][x] = probability
		}
	}
	locationFactorWithMaxReasonsForNotGrowing := 0
	for j := 0; j < len(reasonsForNotGrowing); j++ {
		if j >= 2 { // soil demand should be skipped because it is a obvious reason
			if reasonsForNotGrowing[j] > reasonsForNotGrowing[locationFactorWithMaxReasonsForNotGrowing] {
				locationFactorWithMaxReasonsForNotGrowing = j
			}
		}
	}
	locationFactors := []string{"insolation", "soil demand", "soil depth", "water demand"}
	fmt.Println("Main reason for not growing (except soil demand): " + locationFactors[locationFactorWithMaxReasonsForNotGrowing])
	fmt.Println("Reasons for not growing: ", reasonsForNotGrowing, " (insolation, soil demand, soil depth, water demand)")
	return finalProbabilities
}

func new2DArr(size int) [][]float64 {
	arr := make([][]float64, size)
	for y := 0; y < size; y++ {
		arr[y] = make([]float64, size)
	}
	return arr
}

/*
import numpy as np

from intern.src.python.Data.image import Image


class Probabilities:
    """
    The probabilities class calculates the probability of the growth of a vegetation for every pixel. First
    the probabilites for every location factor (insolation, soil, soil depth and water) is determined. The final
    probability is determined by selecting the lowest probability. This is done because of the law of minimum by
    Liebig, which states that grwoth is dictated by the scarcest resource.
    """
    def __init__(self, controller, vegetation_name, map_name):
        self.controller = controller
        self.map = self.controller.maps[map_name]
        self.vegetation = self.controller.vegetations[vegetation_name]

    @staticmethod
    def calculate_probability(needed, available):
        if available <= needed:
            probability = available / needed
        elif available <= needed * 2:
            probability = 1 - (available / needed - 1)
        else:
            probability = 0
        return probability

    def calculate_insolation_probabilities(self):
        """
        Calculates the insolation probabilities by comparing the available and needed insolation.
        :return: insolation_probabilities: List of calculated probabilities.
        """
        insolation_probabilities = []
        for y in range(self.controller.image_height_map.size):
            print("Calculating insolation probabilities: Row: " + str(y))
            row = []
            for x in range(self.controller.image_height_map.size):
                available_calories = self.controller.image_insolation_map.image[y][x]
                needed_calories = self.vegetation.energy_demand
                probability = self.calculate_probability(needed_calories, available_calories)
                row.append(probability)
            insolation_probabilities.append(row)
        return insolation_probabilities

    def calculate_soil_demand_probabilities(self):
        """
        Calculates the soil demand probabilities by comparing the available and needed soil.
        :return: soil_damand_probabilities: List of calculated probabilities.
        """
        soil_damand_probabilities = []
        for y in range(self.controller.image_height_map.size):
            print("Calculating soil demand probabilities: Row: " + str(y))
            row = []
            for x in range(self.controller.image_height_map.size):
                if self.vegetation.soil_demand.id == self.controller.soil_ids_map.image[y][x]:
                    probability = 1.0
                else:
                    probability = 0.0
                row.append(probability)
            soil_damand_probabilities.append(row)
        return soil_damand_probabilities

    def calculate_soil_depth_probabilities(self):
        """
        Calculates the soil depth probabilities by comparing the available and needed soil depth.
        :return: soil_depth_probabilities: List of calculated probabilities.
        """
        soil_depth_probabilities = []
        for y in range(self.controller.image_height_map.size):
            print("Calculating soil depth probabilities: Row: " + str(y))
            row = []
            for x in range(self.controller.image_height_map.size):
                available_soil_depth = self.controller.image_edaphic_map.image[y][x]
                needed_soil_depth = self.vegetation.soil_depth_demand
                if available_soil_depth < needed_soil_depth:
                    probability = available_soil_depth / needed_soil_depth
                else:
                    probability = 1.0
                row.append(probability)
            soil_depth_probabilities.append(row)
        return soil_depth_probabilities

    def calculate_water_demand_probabilities(self):
        """
        Calculates the water demand probabilities by comparing the available and needed water.
        :return: water_demand_probabilities: List of calculated probabilities.
        """
        water_demand_probabilities = []
        for y in range(self.controller.image_height_map.size):
            print("Calculating water demand probabilites: Row: " + str(y))
            row = []
            for x in range(self.controller.image_height_map.size):
                available_water = self.controller.image_water_map.image[y][x]
                needed_water = self.vegetation.water_demand
                probability = self.calculate_probability(needed_water, available_water)
                row.append(probability)
            water_demand_probabilities.append(row)
        return water_demand_probabilities

    def calculate_probabilities(self):
        """
        Calculates the final probability by selecting the lowest probability at every pixel.
        :return: final_probabilities: List of final probabilities at each pixel.
        """
        all_probabilities = [self.calculate_insolation_probabilities(),
                             self.calculate_soil_demand_probabilities(),
                             self.calculate_soil_depth_probabilities(),
                             self.calculate_water_demand_probabilities()]
        final_probabilities = Image(size=self.controller.image_height_map.size, dtype=np.float)
        reasons_for_not_growing = [0, 0, 0, 0]
        for y in range(self.controller.image_height_map.size):
            for x in range(self.controller.image_height_map.size):
                probability = 1.0
                for i in range(len(all_probabilities)):
                    if all_probabilities[i][y][x] < probability:
                        probability = all_probabilities[i][y][x]
                        if probability == 0.0:
                            reasons_for_not_growing[i] += 1
                final_probabilities.image[y][x] = probability
        location_factor_with_max_reasons_for_not_growing = 0
        for j in range(len(reasons_for_not_growing)):
            if j >= 2:  # soil demand should be skipped because it is a obvious reason
                if reasons_for_not_growing[j] > reasons_for_not_growing[location_factor_with_max_reasons_for_not_growing]:
                    location_factor_with_max_reasons_for_not_growing = j
        location_factors = ["insolation", "soil demand", "soil depth", "water demand"]
        print("Main reason for not growing (except soil demand): " + location_factors[location_factor_with_max_reasons_for_not_growing])
        return final_probabilities
*/
