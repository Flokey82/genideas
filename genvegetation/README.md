# genvegetation

WIP port of Martin Lesser's https://github.com/MartinLesser/Procedural-distribution-of-vegetation to Go.

The application calculates the growth probability of a given vegetation type for every pixel on a given heightmap.

Input (heightmap and soil map):

![alt text](/genvegetation/images/heightmap.png "heightmap")

Note that the soil map is mocked up. Usually you'd generate it based on erosion, deposition, and all that.

![alt text](/genvegetation/images/soilmap.png "soilmap")

Normal map:

![alt text](/genvegetation/images/orographic.png "orographic")

Steepness and soil depth:

![alt text](/genvegetation/images/edaphic.png "edaphic")

![alt text](/genvegetation/images/edaphic_angles.png "edaphic_angles")

Water retention and evaporation based on soil type:

![alt text](/genvegetation/images/water.png "water")

Average insolation throughout the day:

![alt text](/genvegetation/images/insolation.png "insolation")

Probability of vegetation growth (of a specific type):

![alt text](/genvegetation/images/probability.png "probability")

## TODO

- [ ] Use soil types based on erosion and deposition
- [ ] Figure out insolation based on latitude and time of year
- [ ] Refactor to allow for spherical coordinates
