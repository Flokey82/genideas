package genvegetation

type Map struct {
	name              string
	biom              *Biom
	height_map_path   string
	texture_map_path  string
	height_conversion float64 // meters per unit
	max_soil_depth    float64 // in cm
	pixel_size        float64 // in m
}

func NewMap(name string, biom *Biom, height_map_path, texture_map_path string, height_conversion, max_soil_depth, pixel_size float64) *Map {
	return &Map{
		name:              name,
		biom:              biom,
		height_map_path:   height_map_path,
		texture_map_path:  texture_map_path,
		height_conversion: height_conversion,
		max_soil_depth:    max_soil_depth,
		pixel_size:        pixel_size,
	}
}

/*
func (m Map) SaveMap() {
	data := map[string]map[string]any{
		m.name: {
			"biom":              m.biom,
			"height_map_path":   m.height_map_path,
			"texture_map_path":  m.texture_map_path,
			"height_conversion": m.height_conversion,
			"max_soil_depth":    m.max_soil_depth,
			"pixel_size":        m.pixel_size,
		},
	}

	outfile, err := os.OpenFile("resources/data/maps.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	yaml.NewEncoder(outfile).Encode(data)
}
*/
/*
import yaml

class Map:
    def __init__(self, name, biom, height_map_path, texture_map_path, height_conversion, max_soil_depth, pixel_size):
        self.name = name
        self.biom = biom
        self.height_map_path = height_map_path
        self.texture_map_path = texture_map_path
        self.height_conversion = height_conversion  # The factor to convert a height value of the height-map to the actual height
        self.max_soil_depth = max_soil_depth  # in cm, states the maximal depth the ground can have when it has no tilt
        self.pixel_size = pixel_size  # the size that a pixel covers of the real map in m

    def save_map(self):
        data = {self.name: {
            'biom': self.biom,
            'height_map_path': self.height_map_path,
            'texture_map_path': self.texture_map_path,
            'height_conversion': self.height_conversion,
            'max_soil_depth': self.max_soil_depth,
            'pixel_size': self.pixel_size,
            }
        }

        with open('resources/data/maps.yml', 'a') as outfile:
            yaml.dump(data, outfile, default_flow_style=False)
*/
