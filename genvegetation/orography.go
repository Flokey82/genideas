package genvegetation

import (
	"github.com/Flokey82/go_gens/vectors"
)

type Orgraphy struct{}

func NewOrgraphy() *Orgraphy {
	return &Orgraphy{}
}

// CalculateNormal calculates the normal vector from three vertices. First two direction vectors will be calculated and
// via cross product the normal vector will be calculated.
// :param vector1, vector2, vector3: three vertices.
// :return: Cross product of the calculated direction vectors.
func (o *Orgraphy) CalculateNormal(vector1, vector2, vector3 vectors.Vec3) vectors.Vec3 {
	vector_a := vector2.Sub(vector1)
	vector_b := vector3.Sub(vector1)
	return vector_a.Cross(vector_b)
}

// Normalize normalizes a vector by squaring each component and adding the results. The square root of the sum will
// be calculted. The result is used to divide each vector component. This normalizes the vector.
// :param raw: Unnormalized vector.
// :return: Normalized vector.
func (o *Orgraphy) Normalize(raw vectors.Vec3) vectors.Vec3 {
	return raw.Normalize()
}

// AddVectors adds two vectors.
// :param vector1, vector2: Vectors to add.
// :return: Sum of the two vectors.
func (o *Orgraphy) AddVectors(vector1, vector2 vectors.Vec3) vectors.Vec3 {
	return vector1.Add(vector2)
}

// CreateVertexList creates a vertex list. This uses the height map. Each pixel will be transformed to a vertex. The x- and y-position
// will be determined by the pixel position multiplied with the pixel size to receive the real position. The z-position
// is determined by the height value multiplied with height conversion value. Furthermore the list will get a
// padding by repeating the edge vertices. This is necessary for calculating the normals, which need the neighbours.
// :param map: Object of the map class. Used for the pixel size and the height conversion value.
// :param image_height_map: Image of the height map.
// :return: vertex_list: List of the determined vertices.
func (o *Orgraphy) CreateVertexList(m *Map, image_height_map [][]float64) [][]vectors.Vec3 {
	/*
	   vertex_list = []
	   row = [[0, 0, image_height_map.image[0][0] * map.height_conversion]]
	   for x in range(1, image_height_map.size + 1):
	       row.append([map.pixel_size * x, 0, image_height_map.image[0][x - 1] * map.height_conversion])
	   row.append([(image_height_map.size + 1) * map.pixel_size, 0,
	               image_height_map.image[0][image_height_map.size - 1] * map.height_conversion])
	   vertex_list.append(row)
	   for y in range(1, image_height_map.size + 1):
	       row = [[0, y * map.pixel_size, image_height_map.image[y - 1][0] * map.height_conversion]]
	       for x in range(1, image_height_map.size + 1):
	           row.append(
	               [x * map.pixel_size, y * map.pixel_size,
	                image_height_map.image[y - 1][x - 1] * map.height_conversion])
	       row.append(
	           [(image_height_map.size + 1) * map.pixel_size,
	            y * map.pixel_size,
	            image_height_map.image[y - 1][image_height_map.size - 1] * map.height_conversion])
	       vertex_list.append(row)
	   row = [[0, (image_height_map.size + 1) * map.pixel_size,
	           image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion]]
	   for x in range(1, image_height_map.size + 1):
	       row.append([x * map.pixel_size,
	                   (image_height_map.size + 1) * map.pixel_size,
	                   image_height_map.image[image_height_map.size - 1][x - 1] * map.height_conversion])
	   row.append([(image_height_map.size + 1) * map.pixel_size,
	               (image_height_map.size + 1) * map.pixel_size,
	               image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion])
	   vertex_list.append(row)
	   return vertex_list
	*/
	/*

	   vertex_list = []
	   row = [[0, 0, image_height_map.image[0][0] * map.height_conversion]]
	   for x in range(1, image_height_map.size + 1):
	       row.append([map.pixel_size * x, 0, image_height_map.image[0][x - 1] * map.height_conversion])
	   row.append([(image_height_map.size + 1) * map.pixel_size, 0,
	               image_height_map.image[0][image_height_map.size - 1] * map.height_conversion])
	   vertex_list.append(row)
	*/
	vertex_list := make([][]vectors.Vec3, len(image_height_map)+2)
	row := make([]vectors.Vec3, len(image_height_map)+2)
	row[0] = vectors.Vec3{X: 0, Y: 0, Z: image_height_map[0][0] * m.height_conversion}
	for x := 1; x < len(image_height_map)+1; x++ {
		row[x] = vectors.Vec3{X: float64(x) * m.pixel_size, Y: 0, Z: image_height_map[0][x-1] * m.height_conversion}
	}
	row[len(image_height_map)+1] = vectors.Vec3{X: float64(len(image_height_map)+1) * m.pixel_size, Y: 0, Z: image_height_map[0][len(image_height_map)-1] * m.height_conversion}
	vertex_list[0] = row

	/*

	   for y in range(1, image_height_map.size + 1):
	       row = [[0, y * map.pixel_size, image_height_map.image[y - 1][0] * map.height_conversion]]
	       for x in range(1, image_height_map.size + 1):
	           row.append(
	               [x * map.pixel_size, y * map.pixel_size,
	                image_height_map.image[y - 1][x - 1] * map.height_conversion])
	       row.append(
	           [(image_height_map.size + 1) * map.pixel_size,
	            y * map.pixel_size,
	            image_height_map.image[y - 1][image_height_map.size - 1] * map.height_conversion])
	       vertex_list.append(row)
	*/

	for y := 1; y < len(image_height_map)+1; y++ {
		row := make([]vectors.Vec3, len(image_height_map)+2)
		row[0] = vectors.Vec3{X: 0, Y: float64(y) * m.pixel_size, Z: image_height_map[y-1][0] * m.height_conversion}
		for x := 1; x < len(image_height_map)+1; x++ {
			row[x] = vectors.Vec3{X: float64(x) * m.pixel_size, Y: float64(y) * m.pixel_size, Z: image_height_map[y-1][x-1] * m.height_conversion}
		}
		row[len(image_height_map)+1] = vectors.Vec3{X: float64(len(image_height_map)+1) * m.pixel_size, Y: float64(y) * m.pixel_size, Z: image_height_map[y-1][len(image_height_map)-1] * m.height_conversion}
		vertex_list[y] = row
	}

	/*
	 row = [[0, (image_height_map.size + 1) * map.pixel_size,
	           image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion]]
	   for x in range(1, image_height_map.size + 1):
	       row.append([x * map.pixel_size,
	                   (image_height_map.size + 1) * map.pixel_size,
	                   image_height_map.image[image_height_map.size - 1][x - 1] * map.height_conversion])
	   row.append([(image_height_map.size + 1) * map.pixel_size,
	               (image_height_map.size + 1) * map.pixel_size,
	               image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion])
	   vertex_list.append(row)
	   return vertex_list
	*/
	row = make([]vectors.Vec3, len(image_height_map)+2)
	row[0] = vectors.Vec3{X: 0, Y: float64(len(image_height_map)+1) * m.pixel_size, Z: image_height_map[len(image_height_map)-1][0] * m.height_conversion}
	for x := 1; x < len(image_height_map)+1; x++ {
		row[x] = vectors.Vec3{X: float64(x) * m.pixel_size, Y: float64(len(image_height_map)+1) * m.pixel_size, Z: image_height_map[len(image_height_map)-1][x-1] * m.height_conversion}
	}
	row[len(image_height_map)+1] = vectors.Vec3{X: float64(len(image_height_map)+1) * m.pixel_size, Y: float64(len(image_height_map)+1) * m.pixel_size, Z: image_height_map[len(image_height_map)-1][len(image_height_map)-1] * m.height_conversion}
	vertex_list[len(image_height_map)+1] = row
	return vertex_list
}

// CalculateNormalMap calculates all normal vector of a map. It needs the previously calculated vertex list. For calculating the
// normal of a vertex all neighbour normal will be calculated, added up and normalized. Each vertex has six
// neighbour faces, which have a normal vector. These normals have to be calculated. The calculation is done by
// determining the direction vectors of three vertices (the oberserved vertex + 2 neighbour vertices). The
// cross product of the two directions vectors will be calculated resulting in the normal of that surface.
// :param map: Object of the map class. Used for creating the vertex list.
// :param image_height_map: Image of the height map. Used for creating the vertex list.
// :return: normals: List of all calculated normals of each pixel.
func (o *Orgraphy) CalculateNormalMap(m *Map, image_height_map [][]float64) [][]vectors.Vec3 {
	vertex_list := o.CreateVertexList(m, image_height_map)
	array_size := len(vertex_list[0]) - 2
	padded_array := vertex_list
	normals := make([][]vectors.Vec3, array_size)
	for y := range normals {
		normals[y] = make([]vectors.Vec3, array_size)
	}
	alternative := false
	if alternative {
		for y := 1; y < array_size+1; y++ {
			for x := 1; x < array_size+1; x++ {
				// first neighbour triangle
				triangle_normal := o.CalculateNormal(padded_array[y][x], padded_array[y][x-1], padded_array[y-1][x-1])
				sum := triangle_normal
				// second neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y-1][x-1], padded_array[y-1][x])
				sum = o.AddVectors(sum, triangle_normal)
				// third neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y-1][x], padded_array[y][x+1])
				sum = o.AddVectors(sum, triangle_normal)
				// fourth neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y][x+1], padded_array[y+1][x+1])
				sum = o.AddVectors(sum, triangle_normal)
				// fifth neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y+1][x+1], padded_array[y+1][x])
				sum = o.AddVectors(sum, triangle_normal)
				// sixth neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y+1][x], padded_array[y][x-1])
				sum = o.AddVectors(sum, triangle_normal)
				sum = o.Normalize(sum)
				normals[y-1][x-1] = sum
			}
		}
	} else {
		// Calculate a more simple normal map
		for y := 1; y < array_size+1; y++ {
			for x := 1; x < array_size+1; x++ {
				// first neighbour triangle
				triangle_normal := o.CalculateNormal(padded_array[y][x], padded_array[y][x-1], padded_array[y-1][x-1])
				sum := triangle_normal
				// second neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y-1][x-1], padded_array[y-1][x])
				sum = o.AddVectors(sum, triangle_normal)
				// third neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y-1][x], padded_array[y][x+1])
				sum = o.AddVectors(sum, triangle_normal)
				// fourth neighbour triangle
				triangle_normal = o.CalculateNormal(padded_array[y][x], padded_array[y][x+1], padded_array[y+1][x+1])
				sum = o.AddVectors(sum, triangle_normal)
				sum = o.Normalize(sum)
				normals[y-1][x-1] = sum
			}
		}
		/*
			for y := 0; y < array_size; y++ {
				for x := 0; x < array_size; x++ {
					// North neighbour
					n := y - 1
					if n < 0 {
						n = 0
					}

					// South neighbour
					s := y + 1
					if s > array_size-1 {
						s = array_size - 1
					}

					// West neighbour
					w := x - 1
					if w < 0 {
						w = 0
					}

					// East neighbour
					e := x + 1
					if e > array_size-1 {
						e = array_size - 1
					}

					// Calculate the normal of the current pixel
					hNw := image_height_map[n][w]
					hN := image_height_map[n][x]
					hNe := image_height_map[n][e]
					hW := image_height_map[y][w]
					hE := image_height_map[y][e]
					hSw := image_height_map[s][w]
					hS := image_height_map[s][x]
					hSe := image_height_map[s][e]

					// Calculate the normal of the current pixel
					xNormal := -(hSe - hSw + 2*(hE-hW) + hNe - hNw)
					yNormal := -(hNw - hSw + 2*(hN-hS) + hNe - hSe)

					// Calculate the normal of the current pixel
					normal := vectors.Vec3{X: xNormal, Y: yNormal, Z: 1 / 6}.Normalize()
					normals[y][x] = normal
				}
			}
		*/
	}

	return normals
}

/*

import math
import numpy as np


class Orography:
    """
    The orography class calculates the normal map based on the height map.
    """
    @staticmethod
    def calculate_normal(vector1, vector2, vector3):
        """
        Calculates the normal vector from three vertices. First two direction vectors will be calculated and
        via cross product the normal vector will be calculated.
        :param vector1, vector2, vector3: three vertices.
        :return: Cross product of the calculated direction vectors.
        """
        vector_a = [a_i - b_i for a_i, b_i in zip(vector2, vector1)]
        vector_b = [a_i - b_i for a_i, b_i in zip(vector3, vector1)]
        return np.cross(vector_a, vector_b).tolist()

    @staticmethod
    def normalize(raw):
        """
        Normalizes a vector by squaring each component and adding the results. The square root of the sum will
        be calculted. The result is used to divide each vector component. This normalizes the vector.
        :param raw: Unnormalized vector.
        :return: Normalized vector.
        """
        sum = 0
        for v in raw:
            sum += v ** 2
        length = math.sqrt(sum)
        return [float(i) / length for i in raw]

    @staticmethod
    def add_vectors(vector1, vector2):
        return [a_i + b_i for a_i, b_i in zip(vector1, vector2)]

    @staticmethod
    def create_vertex_list(map, image_height_map):
        """
        Creates a vertex list. This uses the height map. Each pixel will be transformed to a vertex. The x- and y-position
        will be determined by the pixel position multiplied with the pixel size to receive the real position. The z-position
        is determined by the height value multiplied with height conversion value. Furthermore the list will get a
        padding by repeating the edge vertices. This is necessary for calculating the normals, which need the neighbours.
        :param map: Object of the map class. Used for the pixel size and the height conversion value.
        :param image_height_map: Image of the height map.
        :return: vertex_list: List of the determined vertices.
        """
        vertex_list = []
        row = [[0, 0, image_height_map.image[0][0] * map.height_conversion]]
        for x in range(1, image_height_map.size + 1):
            row.append([map.pixel_size * x, 0, image_height_map.image[0][x - 1] * map.height_conversion])
        row.append([(image_height_map.size + 1) * map.pixel_size, 0,
                    image_height_map.image[0][image_height_map.size - 1] * map.height_conversion])
        vertex_list.append(row)
        for y in range(1, image_height_map.size + 1):
            row = [[0, y * map.pixel_size, image_height_map.image[y - 1][0] * map.height_conversion]]
            for x in range(1, image_height_map.size + 1):
                row.append(
                    [x * map.pixel_size, y * map.pixel_size,
                     image_height_map.image[y - 1][x - 1] * map.height_conversion])
            row.append(
                [(image_height_map.size + 1) * map.pixel_size,
                 y * map.pixel_size,
                 image_height_map.image[y - 1][image_height_map.size - 1] * map.height_conversion])
            vertex_list.append(row)
        row = [[0, (image_height_map.size + 1) * map.pixel_size,
                image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion]]
        for x in range(1, image_height_map.size + 1):
            row.append([x * map.pixel_size,
                        (image_height_map.size + 1) * map.pixel_size,
                        image_height_map.image[image_height_map.size - 1][x - 1] * map.height_conversion])
        row.append([(image_height_map.size + 1) * map.pixel_size,
                    (image_height_map.size + 1) * map.pixel_size,
                    image_height_map.image[image_height_map.size - 1][image_height_map.size - 1] * map.height_conversion])
        vertex_list.append(row)
        return vertex_list

    @staticmethod
    def calculate_normal_map(map, image_height_map):
        """
        Calculates all normal vector of a map. It needs the previously calculated vertex list. For calculating the
        normal of a vertex all neighbour normal will be calculated, added up and normalized. Each vertex has six
        neighbour faces, which have a normal vector. These normals have to be calculated. The calculation is done by
        determining the direction vectors of three vertices (the oberserved vertex + 2 neighbour vertices). The
        cross product of the two directions vectors will be calculated resulting in the normal of that surface.
        :param map: Object of the map class. Used for creating the vertex list.
        :param image_height_map: Image of the height map. Used for creating the vertex list.
        :return: normals: List of all calculated normals of each pixel.
        """
        vertex_list = Orography.create_vertex_list(map, image_height_map)
        array_size = len(vertex_list[0]) - 2
        padded_array = vertex_list
        normals = []
        for y in range(1, array_size + 1):
            print("Calculating normal map. Row: " + str(y))
            row = []
            for x in range(1, array_size + 1):
                # print("Calculating normal map. Column: " + str(x))
                sum = [0.0, 0.0, 0.0]
                # first neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y][x - 1],
                                                             padded_array[y - 1][x - 1],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                # second neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y - 1][x - 1],
                                                             padded_array[y - 1][x],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                # third neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y - 1][x],
                                                             padded_array[y][x + 1],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                # fourth neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y][x + 1],
                                                             padded_array[y + 1][x + 1],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                # fifth neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y + 1][x + 1],
                                                             padded_array[y + 1][x],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                # sixth neighbour triangle
                triangle_normal = Orography.calculate_normal(padded_array[y][x],
                                                             padded_array[y + 1][x],
                                                             padded_array[y][x - 1],
                                                             )
                sum = Orography.add_vectors(sum, triangle_normal)
                sum = Orography.normalize(sum)
                row.append(sum)
            normals.append(row)
        return normals
*/
