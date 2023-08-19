package genvegetation

type Vegetation struct {
	name              string
	energy_demand     float64 // in kcal/day
	water_demand      float64 // in l/day
	soil_demand       *Soil
	soil_depth_demand float64 // in cm
}

func NewVegetation(name string, energy_demand, water_demand, soil_depth_demand float64, soil_demand *Soil) *Vegetation {
	return &Vegetation{
		name:              name,
		energy_demand:     energy_demand,
		water_demand:      water_demand,
		soil_demand:       soil_demand,
		soil_depth_demand: soil_depth_demand,
	}
}

/*
func (v Vegetation) SaveVegetation() {
	data := map[string]map[string]any{
		v.name: {
			"energy_demand":     v.energy_demand,
			"water_demand":      v.water_demand,
			"soil_demand":       v.soil_demand,
			"soil_depth_demand": v.soil_depth_demand,
		},
	}

	outfile, err := os.OpenFile("resources/data/vegetation_types.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	yaml.NewEncoder(outfile).Encode(data)
}
*/
/*
import yaml


class Vegetation:
    def __init__(self, name, energy_demand, water_demand, soil_demand, soil_depth_demand):
        self.name = name
        self.energy_demand = energy_demand
        self.water_demand = water_demand
        self.soil_demand = soil_demand
        self.soil_depth_demand = soil_depth_demand

    def save_vegetation(self):
        data = {self.name: {
            'energy_demand': self.energy_demand,
            'water_demand': self.water_demand,
            'soil_demand': self.soil_demand,
            'soil_depth_demand': self.soil_depth_demand }
        }

        with open('resources/data/vegetation_types.yml', 'a') as outfile:
            yaml.dump(data, outfile, default_flow_style=False)
*/
