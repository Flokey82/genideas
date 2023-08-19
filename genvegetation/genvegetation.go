package genvegetation

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

type Controller struct {
	image_height_map     [][]float64
	soil_ids_map         [][]int
	image_insolation_map *Image
	image_orographic_map [][]vectors.Vec3
	image_edaphic_map    [][]float64 // soil depth
	image_water_map      [][]float64
	image_probabilities  [][]float64
	bioms                map[string]*Biom
	soils                map[string]*Soil
	vegetations          map[string]*Vegetation
	maps                 map[string]*Map
	//main_window          *MainWindow
}

func NewController() *Controller {
	c := &Controller{
		bioms:       make(map[string]*Biom),
		soils:       make(map[string]*Soil),
		vegetations: make(map[string]*Vegetation),
		maps:        make(map[string]*Map),
	}
	c.loadBioms()
	c.loadSoils()
	c.loadVegetations()
	c.loadMaps()
	//c.main_window = NewMainWindow(c)
	return c
}

func (c *Controller) getMap(name string) *Map {
	return c.maps[name]
}

// LoadHeightAndSoilMap creates two images and loads the height- and soil- map into it.
// :param map_name: String of the map name. It is used to find the images on disk.
func (c *Controller) LoadHeightAndSoilMap(map_name string) {
	m := c.getMap(map_name)
	c.image_height_map = loadFloatPNG(m.height_map_path)
	c.soil_ids_map = loadUint16AsIntPNG(m.texture_map_path)

	log.Println("Hack: Transforming soil ids to valid ids.")
	unique := FilterUniqueNumbersFrom2DArray(c.soil_ids_map)
	// Now map it as best as we can to:
	// 0: NotFound
	// 40: Dirt
	// 80: Loam
	// 120: Silt
	// 160: Clay
	// 200: Stone
	// 240: Snow
	replace := []int{0, 40, 80, 120, 160, 200, 240}
	res := make(map[int]int)
	for i, v := range unique {
		res[v] = replace[i]
	}
	if len(res) != len(replace) {
		log.Fatalf("Error: Could not map all soil ids to valid ids. %d != %d", len(res), len(replace))
	}
	TransformImageToValidSoils(c.soil_ids_map, res)
	// c.transformAndSaveSoilIdMap(map.texture_map_path)
	// c.saveImageAsCsv(c.image_height_map.image)
}

// TransformAndSaveSoilIdMap is used to transform all occuring IDs in a soil map to wanted IDs (usually the IDs of the
// created soils). To transform the IDs you have to change the transformation list. This function is used
// as a workaround if you can't find a valid soil-map.
// :param path: String of the path to the map that shall be transformed.
func (c *Controller) TransformAndSaveSoilIdMap(path string) {
	FilterUniqueNumbersFrom2DArray(c.soil_ids_map)
	transformation_list := map[int]int{0: 40, 9362: 80, 18724: 120, 28086: 160, 37449: 200, 46811: 240}
	TransformImageToValidSoils(c.soil_ids_map, transformation_list)
	FilterUniqueNumbersFrom2DArray(c.soil_ids_map) // check if the transformation was successfull
	saveUint16AsIntPNG(path, c.soil_ids_map)
	//c.soil_ids_map.SaveImage(path)
}

// PrepareInsolationCalculation finds the correct map obect and creates an image and an object of insolation class. Then the heightmap
// gets loaded and the calculation of the insolation gets started. The results will be shown in the UI and saved.
// :param map_name: String of the current map name
// :param daylight_hours: Integer of the number of daylight hours
// :param sun_start_elevation: Float of the start elevation of the sun
// :param sun_start_azimuth: Float of the start azimuth of the sun
// :param sun_max_elevation: Float of the maximal sun elevation (noon)
// :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
//
//	pixel will be reflected.
func (c *Controller) PrepareInsolationCalculation(map_name string, daylight_hours int, sun_start_elevation, sun_start_azimuth, sun_max_elevation, reflection_coefficient float64) {
	m := c.getMap(map_name)
	insolation := NewInsolation(c)
	c.image_height_map = loadFloatPNG(m.height_map_path)
	c.image_insolation_map = insolation.CalculateActualInsolation(m, daylight_hours, sun_start_elevation, sun_start_azimuth, sun_max_elevation, reflection_coefficient)
	//c.main_window.frames["ProbabilityCloudWindow"].DrawInsolationImage(c.image_insolation_map)
	save_path := "resources/results/" + map_name + "/" + map_name + "_" + string(daylight_hours) + "daylight_hours_insolation_image.png"
	c.image_insolation_map.SaveImage(save_path)
	saveFloatPNG(map_name+"_insolation_image.png", c.image_insolation_map.image)
}

// PrepareOrographicCalculation loads the correct map object, starts the orograhic calculation, displays the result in the UI and
// saves it as an image.
// :param map_name: name of the current map
func (c *Controller) PrepareOrographicCalculation(map_name string) {
	m := c.getMap(map_name)
	o := new(Orgraphy)
	c.image_height_map = loadFloatPNG(m.height_map_path)
	c.image_orographic_map = o.CalculateNormalMap(m, c.image_height_map)
	saveVec3ToPNG(map_name+"_orographic_normals.png", c.image_orographic_map)
	//c.main_window.frames["ProbabilityCloudWindow"].DrawOrographicImage(c.image_orographic_map)
	c.save3DList(c.image_orographic_map, "resources/results/"+map_name+"/"+map_name+"_orographic_normals")
}

// PrepareEdaphicCalculation loads the correct map object, calculates all angles on the map (between the normal vector and the z-vector),
// starts the edaphic calculation, displays the result in the UI and saves it as an image.
// :param map_name: name of the current map
func (c *Controller) PrepareEdaphicCalculation(map_name string) {
	m := c.getMap(map_name)
	e := new(Edaphology)
	angles := e.CalculateAngles(c.image_orographic_map)
	saveFloatPNG(map_name+"_edaphic_image_angles.png", angles)
	c.image_edaphic_map = e.CalculateSoilDepth(m, len(c.image_height_map), angles)
	//c.main_window.frames["ProbabilityCloudWindow"].DrawEdaphicImage(c.image_edaphic_map)
	//c.image_edaphic_map.SaveImage("resources/results/" + map_name + "/" + map_name + "_edaphic_image.png")
	saveFloatPNG(map_name+"_edaphic_image.png", c.image_edaphic_map)
}

// PrepareWaterCalculation loads the correct map and biom object, starts the hydrologic calculation, displays the result in the UI and
// saves it as an image.
// :param map_name: name of the current map
func (c *Controller) PrepareWaterCalculation(map_name string) {
	m := c.getMap(map_name)
	biom := m.biom
	hydrology := NewHydrology(c)
	c.image_water_map = hydrology.CalculateHydrologyMap(m, c.image_edaphic_map, c.soil_ids_map, c.image_insolation_map, biom)
	//c.main_window.frames["ProbabilityCloudWindow"].DrawHydrologyImage(c.image_water_map)
	// c.image_water_map.SaveImage("resources/results/" + map_name + "/" + map_name + "_water_image.png")
	saveFloatPNG(map_name+"_water_image.png", c.image_water_map)
}

// CalculateAll starts all calculation. This function is used for very large maps. So the user does not have to check the
// state of the application all the time. You can start all calculations and come back a few hours later.
// :param map_name: String of the current map name.
// :param daylight_hours: Integer of the number of daylight hours.
// :param sun_start_elevation: Float of the start elevation of the sun.
// :param sun_start_azimuth: Float of the start azimuth of the sun.
// :param sun_max_elevation: Float of the maximal sun elevation (noon).
// :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
//
//	pixel will be reflected.
func (c *Controller) CalculateAll(map_name string, daylight_hours int, sun_start_elevation, sun_start_azimuth, sun_max_elevation, reflection_coefficient float64) {
	c.PrepareInsolationCalculation(map_name, daylight_hours, sun_start_elevation, sun_start_azimuth, sun_max_elevation, reflection_coefficient)
	c.PrepareOrographicCalculation(map_name)
	c.PrepareEdaphicCalculation(map_name)
	c.PrepareWaterCalculation(map_name)
}

// PrepareProbabilitesCalculation creates a probability object, starts the calculation of the probabilites and displays the result in the UI.
// :param vegetation_name: String. Name of the vegetation for which the probabilites will be calculated.
// :param map_name: String. Name of the current map.
func (c *Controller) PrepareProbabilitesCalculation(vegetation_name, map_name string) {
	probability_calculator := NewProbabilities(c, c.SearchVegetation(vegetation_name), map_name)
	c.image_probabilities = probability_calculator.CalculateProbabilities()
	//c.main_window.frames["ProbabilityCloudWindow"].DrawProbabilityImage(c.image_probabilities)
	saveFloatPNG(map_name+"_probability_image.png", c.image_probabilities)
}

// SearchSoil searches for the soil object in the soil list with the help of the soil ID.
// :param soil_id: ID of the soil for which the object from the list shall be
// :return: soil: Object of the soil class.
func (c *Controller) SearchSoil(soil_id int) *Soil {
	for _, soil := range c.soils {
		if soil.id == soil_id {
			return soil
		}
	}
	return c.soils["NotFound"]
}

// SearchVegetation searches for the vegetation object in the vegetation list with the help of the vegetation name.
// :param vegetation_name: Name of the vegetation for which the object from the list shall be
// :return: vegetation: Object of the vegetation class.
func (c *Controller) SearchVegetation(vegetation_name string) *Vegetation {
	for _, vegetation := range c.vegetations {
		if vegetation.name == vegetation_name {
			return vegetation
		}
	}
	return nil
}

// SearchMap searches for the map object in the map list with the help of the map name.
// :param map_name: Name of the map for which the object from the list shall be
// :return: map: Object of the map class.
func (c *Controller) SearchMap(map_name string) *Map {
	for _, m := range c.maps {
		if m.name == map_name {
			return m
		}
	}
	return nil
}

// LoadBioms loads all bioms from the biom file and saves them in the biom list.
func (c *Controller) loadBioms() {
	/*
			PolarZone:
		  atmospheric_absorption: '45'
		  atmospheric_diffusion: '8'
		  avg_rainfall_per_day: '1.5'
		  cloud_reflection: '40'
		  groundwater: '1.0'
		TemperateZone:
		  atmospheric_absorption: '35'
		  atmospheric_diffusion: '9'
		  avg_rainfall_per_day: '2.2'
		  cloud_reflection: '33'
		  groundwater: '2.0'
		TropicalZone:
		  atmospheric_absorption: '30'
		  atmospheric_diffusion: '10'
		  avg_rainfall_per_day: '3.0'
		  cloud_reflection: '28'
		  groundwater: '2.5'
		SubtropicalZone:
		  atmospheric_absorption: '25'
		  atmospheric_diffusion: '11'
		  avg_rainfall_per_day: '4.0'
		  cloud_reflection: '10'
		  groundwater: '2.5'
	*/
	c.bioms["PolarZone"] = NewBiom("PolarZone", 8, 45, 40, 1.5, 1.0)
	c.bioms["TemperateZone"] = NewBiom("TemperateZone", 9, 35, 33, 2.2, 2.0)
	c.bioms["TropicalZone"] = NewBiom("TropicalZone", 10, 30, 28, 3.0, 2.5)
	c.bioms["SubtropicalZone"] = NewBiom("SubtropicalZone", 11, 25, 10, 4.0, 2.5)
}

// LoadSoils loads all soils from the soil file and saves them in the soil list.
func (c *Controller) loadSoils() {
	/*
			Dirt:
		  id: '40'
		  albedo: '0.25'
		  water_absorption: '0.38'
		Loam:
		  id: '80'
		  albedo: '0.25'
		  water_absorption: '0.43'
		Silt:
		  id: '120'
		  albedo: '0.25'
		  water_absorption: '0.4'
		Clay:
		  id: '160'
		  albedo: '0.3'
		  water_absorption: '0.42'
		Stone:
		  id: '200'
		  albedo: '0.275'
		  water_absorption: '0.28'
		Snow:
		  id: '240'
		  albedo: '0.84'
		  water_absorption: '0.01'
		NotFound:
		  id: '0'
		  albedo: '0.0'
		  water_absorption: '0.0'
	*/
	c.soils["Dirt"] = NewSoil(40, "Dirt", 0.25, 0.38)
	c.soils["Loam"] = NewSoil(80, "Loam", 0.25, 0.43)
	c.soils["Silt"] = NewSoil(120, "Silt", 0.25, 0.4)
	c.soils["Clay"] = NewSoil(160, "Clay", 0.3, 0.42)
	c.soils["Stone"] = NewSoil(200, "Stone", 0.275, 0.28)
	c.soils["Snow"] = NewSoil(240, "Snow", 0.84, 0.01)
	c.soils["NotFound"] = NewSoil(0, "NotFound", 0.0, 0.0)
}

// LoadVegetations loads all vegetations from the vegetation file and saves them in the vegetation list.
func (c *Controller) loadVegetations() {
	/*
			Appletree:
		  energy_demand: '2000'
		  soil_demand: Dirt
		  soil_depth_demand: '80'
		  water_demand: '0.7'
		Shrub:
		  energy_demand: '1200'
		  soil_demand: Dirt
		  soil_depth_demand: '20'
		  water_demand: '0.3'
		Moss:
		  energy_demand: '1000'
		  soil_demand: Stone
		  soil_depth_demand: '1'
		  water_demand: '0.3'
		Beech:
		  energy_demand: '2200'
		  soil_demand: Loam
		  soil_depth_demand: '50'
		  water_demand: '0.5'
		Fern:
		  energy_demand: '1100'
		  soil_demand: Silt
		  soil_depth_demand: '10'
		  water_demand: '0.35'
		Fir:
		  energy_demand: '2200'
		  soil_demand: Clay
		  soil_depth_demand: '45'
		  water_demand: '0.6'
		Cactus:
		  energy_demand: '4000'
		  soil_demand: Dirt
		  soil_depth_demand: '30'
		  water_demand: '0.2'
	*/
	c.vegetations["Appletree"] = NewVegetation("Appletree", 2000, 0.7, 80, c.soils["Dirt"])
	c.vegetations["Shrub"] = NewVegetation("Shrub", 1200, 0.3, 20, c.soils["Dirt"])
	c.vegetations["Moss"] = NewVegetation("Moss", 1000, 0.3, 1, c.soils["Stone"])
	c.vegetations["Beech"] = NewVegetation("Beech", 2200, 0.5, 50, c.soils["Loam"])
	c.vegetations["Fern"] = NewVegetation("Fern", 1100, 0.35, 10, c.soils["Silt"])
	c.vegetations["Fir"] = NewVegetation("Fir", 2200, 0.6, 45, c.soils["Clay"])
	c.vegetations["Cactus"] = NewVegetation("Cactus", 4000, 0.2, 30, c.soils["Dirt"])
}

// LoadMaps loads all maps from the map file and saves them in the map list.
func (c *Controller) loadMaps() {

	//fileHeightMap := "../maps/16bitHeightmap64x64.png"
	//fileSoilMap := "../maps/16bitHeightmap64x64_soil_ids.png"
	fileHeightMap := "../maps/b_image3.png"
	fileSoilMap := "../maps/b_image3_gray.png"

	c.maps["default"] = &Map{
		name:              "default",
		biom:              c.bioms["TemperateZone"],
		height_map_path:   fileHeightMap,
		texture_map_path:  fileSoilMap,
		height_conversion: 5.0,
		max_soil_depth:    140.0,
		pixel_size:        16.0,
	}
}

func (c *Controller) save3DList(image [][]vectors.Vec3, path string) {
	/*
	   It saves the result of the orographic class (normal map) as a file.
	   :param list: List of the normal vectors of the current map.
	   :param path: String of the path where the file should be saved.
	*/
	//panic("not implemented")
}

func (c *Controller) load3DList(path string) {
	/*
		It loads the result of the orographic class (normal map) from a file.
		:param path: String of the path where the file should be loaded from.
		:return: List of the normal vectors of the current map.
	*/
	panic("not implemented")
}

func (c *Controller) saveImageAsCsv(image [][]float64) {
	/*
		It saves the result of the insolation class as a file.
		:param image: List of the insolation values of the current map.
	*/
	panic("not implemented")
}

/*

import numpy as np
from pathlib import Path
import pickle
import os
import yaml

from intern.src.python.Data.biom import Biom
from intern.src.python.Data.image import Image
from intern.src.python.Data.map import Map
from intern.src.python.Data.soil import Soil
from intern.src.python.Data.vegetation import Vegetation
from intern.src.python.Gui.main_window import MainWindow
from intern.src.python.Logic.edaphology import Edaphology
from intern.src.python.Logic.hydrology import Hydrology
from intern.src.python.Logic.insolation import Insolation
from intern.src.python.Logic.probabilities import Probabilities
from intern.src.python.Logic.orography import Orography


class Controller:
    """
    The Controller class controls the flow of the application. It starts the UI, loads all data and controlls
    the calculation of the maps (insolation, orography, edaphology and hydrology) by starting the other logic
    classes. It also starts the calculation of the probability map. All results will be saved and can be loaded
    later on.
    """

    def __init__(self):
        self.image_height_map = None
        self.soil_ids_map = None
        self.image_insolation_map = None
        self.image_orographic_map = None
        self.image_edaphic_map = None
        self.image_water_map = None
        self.image_probabilities = None
        self.bioms = {}
        self.load_bioms()
        self.soils = {}
        self.load_soils()
        self.vegetations = {}
        self.load_vegetations()
        self.maps = {}
        self.load_maps()
        self.main_window = MainWindow(self)

    def load_height_and_soil_map(self, map_name):
        """
        It creates two images and loads the height- and soil- map into it.
        :param map_name: String of the map name. It is used to find the images on disk.
        """
        map = self.maps[map_name]
        self.image_height_map = Image()
        self.image_height_map.load_image(map.height_map_path)
        self.soil_ids_map = Image()
        self.soil_ids_map.load_image(map.texture_map_path)
        # self.transform_and_save_soil_id_map(map.texture_map_path)
        # self.save_image_as_csv(self.image_height_map.image)

    def transform_and_save_soil_id_map(self, path):
        """
        This function is used to transform all occuring IDs in a soil map to wanted IDs (usually the IDs of the
        created soils). To transform the IDs you have to change the transformation list. This function is used
        as a workaround if you can't find a valid soil-map.
        :param path: String of the path to the map that shall be transformed.
        """
        self.soil_ids_map.filter_unique_numbers_from_2d_array()
        transformation_list = {0: 40, 9362: 80, 18724: 120, 28086: 160, 37449: 200, 46811: 240}
        self.soil_ids_map.transform_image_to_valid_soils(transformation_list)
        self.soil_ids_map.filter_unique_numbers_from_2d_array()  # check if the transformation was successfull
        self.soil_ids_map.save_image(path)

    def prepare_insolation_calculation(self, map_name, daylight_hours, sun_start_elevation, sun_start_azimuth,
                                       sun_max_elevation, reflection_coefficient):
        """
        It finds the correct map obect and creates an image and an object of insolation class. Then the heightmap
        gets loaded and the calculation of the insolation gets started. The results will be shown in the UI and saved.
        :param map_name: String of the current map name
        :param daylight_hours: Integer of the number of daylight hours
        :param sun_start_elevation: Float of the start elevation of the sun
        :param sun_start_azimuth: Float of the start azimuth of the sun
        :param sun_max_elevation: Float of the maximal sun elevation (noon)
        :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
                pixel will be reflected.
        """
        map = self.maps[map_name]
        self.image_insolation_map = Image(size=self.image_height_map.size)
        insolation = Insolation(self)
        self.image_height_map.load_image(map.height_map_path)
        self.image_insolation_map = insolation.calculate_actual_insolation(map, daylight_hours, sun_start_elevation,
                                                                           sun_start_azimuth, sun_max_elevation,
                                                                           reflection_coefficient)
        self.main_window.frames['ProbabilityCloudWindow'].draw_insolation_image(self.image_insolation_map)
        save_path = "resources/results/" + map_name + "/" + map_name + "_" + str(daylight_hours) + "daylight_hours_insolation_image.png"
        self.image_insolation_map.save_image(save_path)

    def prepare_orographic_calculation(self, map_name):
        """
        It loads the correct map object, starts the orograhic calculation, displays the result in the UI and
        saves it as an image.
        :param map_name: name of the current map
        """
        map = self.maps[map_name]
        self.image_orographic_map = Orography.calculate_normal_map(map, self.image_height_map)
        self.main_window.frames['ProbabilityCloudWindow'].draw_orographic_image(self.image_orographic_map)
        self.save_3d_list(self.image_orographic_map, "resources/results/" + map_name + "/" + map_name + "_orographic_normals")

    def prepare_edaphic_calculation(self, map_name):
        """
        It loads the correct map object, calculates all angles on the map (between the normal vector and the z-vector),
        starts the edaphic calculation, displays the result in the UI and saves it as an image.
        :param map_name: name of the current map
        """
        map = self.maps[map_name]
        angles = Edaphology.calculate_angles(self.image_orographic_map)
        self.image_edaphic_map = Edaphology.calculate_soil_depth(map, self.image_height_map.size, angles)
        self.main_window.frames['ProbabilityCloudWindow'].draw_edaphic_image(self.image_edaphic_map)
        self.image_edaphic_map.save_image("resources/results/" + map_name + "/" + map_name + "_edaphic_image.png")

    def prepare_water_calculation(self, map_name):
        """
        It loads the correct map and biom object, starts the hydrologic calculation, displays the result in the UI and
        saves it as an image.
        :param map_name: name of the current map
        """
        map = self.maps[map_name]
        biom = map.biom
        hydrology = Hydrology(self)
        self.image_water_map = hydrology.calculate_hydrology_map(map, self.image_edaphic_map, self.soil_ids_map,
                                                                 self.image_insolation_map, biom)
        self.main_window.frames['ProbabilityCloudWindow'].draw_hydrology_image(self.image_water_map)
        self.image_water_map.save_image("resources/results/" + map_name + "/" + map_name + "_water_image.png")

    def calculate_all(self, map_name, daylight_hours, sun_start_elevation, sun_start_azimuth, sun_max_elevation,
                      reflection_coefficient):
        """
        Starts all calculation. This function is used for very large maps. So the user does not have to check the
        state of the application all the time. You can start all calculations and come back a few hours later.
        :param map_name: String of the current map name.
        :param daylight_hours: Integer of the number of daylight hours.
        :param sun_start_elevation: Float of the start elevation of the sun.
        :param sun_start_azimuth: Float of the start azimuth of the sun.
        :param sun_max_elevation: Float of the maximal sun elevation (noon).
        :param reflection_coefficient: Float of the reflection coeficient. It states how much light of the neighbour
                pixel will be reflected.
        """
        self.prepare_insolation_calculation(map_name, daylight_hours, sun_start_elevation, sun_start_azimuth,
                                            sun_max_elevation, reflection_coefficient)
        self.prepare_orographic_calculation(map_name)
        self.prepare_edaphic_calculation(map_name)
        self.prepare_water_calculation(map_name)

    def prepare_probabilites_calculation(self, vegetation_name, map_name):
        """
        It creates a probability object, starts the calculation of the probabilites and displays the result in the UI.
        :param vegetation_name: String. Name of the vegetation for which the probabilites will be calculated.
        :param map_name: String. Name of the current map.
        """
        probability_calculator = Probabilities(self, vegetation_name, map_name)
        self.image_probabilities = probability_calculator.calculate_probabilities()
        self.main_window.frames['ProbabilityCloudWindow'].draw_probability_image(self.image_probabilities)

    def search_soil(self, soil_id):
        """
        Searches for the soil object in the soil list with the help of the soil ID.
        :param soil_id: ID of the soil for which the object from the list shall be found.
        :return: soil_value: Soil object from the soil list.
        """
        for soil_name, soil_value in self.soils.items():
            if soil_value.id == soil_id:
                return soil_value
        print('Soil id (' + str(soil_id) + ') could not be found!')
        return self.soils['NotFound']  # raise Exception('Soil id could not be found!')

    def load_bioms(self):
        """
        Loads all bioms from the bioms.yml into a list.
        """
        bioms = {}
        bioms_file = Path("resources/data/bioms.yml")
        if bioms_file.is_file():
            with open(bioms_file, 'r') as stream:
                try:
                    bioms_dict = yaml.safe_load(stream)
                    if bioms_dict is not None:
                        for biom_name, biom_values in bioms_dict.items():
                            bioms[biom_name] = Biom(biom_name,
                                                    float(biom_values['atmospheric_diffusion']),
                                                    float(biom_values['atmospheric_absorption']),
                                                    float(biom_values['cloud_reflection']),
                                                    float(biom_values['avg_rainfall_per_day']),
                                                    float(biom_values['groundwater']))
                except yaml.YAMLError as exc:
                    print(exc)
        self.bioms = bioms

    def load_soils(self):
        """
        Loads all soils from the soils.yml into a list.
        """
        soils = {}
        soils_file = Path("resources/data/soil_types.yml")
        if soils_file.is_file():
            with open(soils_file, 'r') as stream:
                try:
                    soils_dict = yaml.safe_load(stream)
                    if soils_dict is not None:
                        for soil_name, soil_values in soils_dict.items():
                            soils[soil_name] = Soil(int(soil_values['id']), soil_name, float(soil_values['albedo']),
                                                    float(soil_values['water_absorption']))
                except yaml.YAMLError as exc:
                    print(exc)
        self.soils = soils

    def load_vegetations(self):
        """
        Loads all vegetations from the vegetation_types.yml into a list.
        """
        vegetations = {}
        vegetations_file = Path("resources/data/vegetation_types.yml")
        if vegetations_file.is_file():
            with open(vegetations_file, 'r') as stream:
                try:
                    vegetations_dict = yaml.safe_load(stream)
                    if vegetations_dict is not None:
                        for vegetation_name, vegetation_values in vegetations_dict.items():
                            vegetations[vegetation_name] = Vegetation(vegetation_name,
                                                                      float(vegetation_values['energy_demand']),
                                                                      float(vegetation_values['water_demand']),
                                                                      self.soils[vegetation_values['soil_demand']],
                                                                      float(vegetation_values['soil_depth_demand']))
                except yaml.YAMLError as exc:
                    print(exc)
        self.vegetations = vegetations

    def load_maps(self):
        """
        Loads all maps from the maps.yml into a list.
        """
        maps = {}
        maps_file = Path("resources/data/maps.yml")
        if maps_file.is_file():
            with open(maps_file, 'r') as stream:
                try:
                    maps_dict = yaml.safe_load(stream)
                    if maps_dict is not None:
                        for map_name, map_values in maps_dict.items():
                            maps[map_name] = Map(map_name, self.bioms[map_values['biom']],
                                                 map_values['height_map_path'],
                                                 map_values['texture_map_path'],
                                                 map_values['height_conversion'],
                                                 map_values['max_soil_depth'],
                                                 map_values['pixel_size'])
                except yaml.YAMLError as exc:
                    print(exc)
        self.maps = maps

    @staticmethod
    def save_3d_list(list, path):
        """
        It saves the result of the orographic class (normal map) as a file.
        :param list: List of the normal vectors of the current map.
        :param path: String of the path where the file should be saved.
        """
        if not os.path.exists(os.path.dirname(path)):
            os.makedirs(os.path.dirname(path))
        output = open(path, 'wb')
        pickle.dump(list, output)
        output.close()

    @staticmethod
    def load_3d_list(path):
        """
        Loads the normal vectors from a file.
        :param path: String of the path of the file that will be loaded.
        :return: list: List of the loaded normal vectors (normal map).
        """
        pkl_file = open(path, 'rb')
        list = pickle.load(pkl_file)
        pkl_file.close()
        return list

    @staticmethod
    def save_image_as_csv(image):
        np.savetxt("resources/height_map.csv", image, delimiter=',', fmt='%s')


if __name__ == "__main__":
    Controller().main_window.mainloop()

*/
