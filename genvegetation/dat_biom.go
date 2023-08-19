package genvegetation

type Biom struct {
	name                   string
	atmospheric_diffusion  float64 // in %
	atmospheric_absorption float64 // in %
	cloud_reflection       float64 // in %
	avg_rainfall_per_day   float64 // in l/cm²
	groundwater            float64 // in l/m³
}

func NewBiom(name string, atmospheric_diffusion, atmospheric_absorption, cloud_reflection,
	avg_rainfall_per_day, groundwater float64) *Biom {
	return &Biom{
		name:                   name,
		atmospheric_diffusion:  atmospheric_diffusion,
		atmospheric_absorption: atmospheric_absorption,
		cloud_reflection:       cloud_reflection,
		avg_rainfall_per_day:   avg_rainfall_per_day,
		groundwater:            groundwater,
	}
}

/*
func (b Biom) SaveBiom() {
	data := map[string]map[string]float64{
		b.name: {
			"atmospheric_diffusion":  b.atmospheric_diffusion,
			"atmospheric_absorption": b.atmospheric_absorption,
			"cloud_reflection":       b.cloud_reflection,
			"avg_rainfall_per_day":   b.avg_rainfall_per_day,
			"groundwater":            b.groundwater,
		},
	}

	outfile, err := os.OpenFile("resources/data/bioms.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	yaml.NewEncoder(outfile).Encode(data)
}
*/
/*
import yaml
class Biom:
    def __init__(self, name, atmospheric_diffusion, atmospheric_absorption, cloud_reflection,
                 avg_rainfall_per_day, groundwater):
        self.name = name
        # this value corresponds to the diffuse sun beam scattering by the atmosphere
        self.atmospheric_diffusion = atmospheric_diffusion  # in percent
        self.atmospheric_absorption = atmospheric_absorption  # in percent
        self.cloud_reflection = cloud_reflection  # in percent
        self.avg_rainfall_per_day = avg_rainfall_per_day  # in l/cm²
        self.groundwater = groundwater  # in l/cm²

    def save_biom(self):
        data = {self.name: {
            'atmospheric_diffusion': self.atmospheric_diffusion,
            'atmospheric_absorption': self.atmospheric_absorption,
            'cloud_reflection': self.cloud_reflection,
            'avg_rainfall_per_day': self.avg_rainfall_per_day,
            'groundwater': self.groundwater}
        }

        with open('resources/data/bioms.yml', 'a') as outfile:
            yaml.dump(data, outfile, default_flow_style=False)

*/
