package genvegetation

type Soil struct {
	id               int
	name             string
	albedo           float64
	water_absorption float64 // l/m³
}

func NewSoil(id int, name string, albedo, water_absorption float64) *Soil {
	return &Soil{
		id:               id,
		name:             name,
		albedo:           albedo,
		water_absorption: water_absorption,
	}
}

/*
func (s Soil) SaveSoil() {
	data := map[string]map[string]any{
		s.name: {
			"id":               s.id,
			"albedo":           s.albedo,
			"water_absorption": s.water_absorption,
		},
	}

	outfile, err := os.OpenFile("resources/data/soil_types.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	yaml.NewEncoder(outfile).Encode(data)
}
*/
/*
import yaml


class Soil:
    def __init__(self, id, name, albedo, water_absorption):
        self.id = id
        self.name = name
        self.albedo = albedo
        self.water_absorption = water_absorption

    def save_soil(self):
        data = {self.name: {
            'id': self.id,
            'albedo': self.albedo,
            'water_absorption': self.water_absorption }
        }

        with open('resources/data/soil_types.yml', 'a') as outfile:
            yaml.dump(data, outfile, default_flow_style=False)
*/
