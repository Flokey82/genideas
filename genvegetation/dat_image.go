package genvegetation

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/Flokey82/go_gens/vectors"
)

type DataType int

func saveVec3ToPNG(path string, imageData [][]vectors.Vec3) {
	width, height := len(imageData), len(imageData[0])

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			vec := imageData[y][x]
			vec.X = (vec.X * 0.5) + 0.5
			vec.Y = (vec.Y * 0.5) + 0.5
			vec.Z = (vec.Z * 0.5) + 0.5
			col := color.NRGBA{
				R: uint8(vec.X * 255),
				G: uint8(vec.Y * 255),
				B: uint8(vec.Z * 255),
				A: 255,
			}
			img.Set(x, y, col)
		}
	}

	// save the new image
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func minMax(imageData [][]float64) (float64, float64) {
	if len(imageData) == 0 {
		return 0, 0
	}
	min, max := imageData[0][0], imageData[0][0]
	for _, row := range imageData {
		for _, val := range row {
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
	}
	return min, max
}
func saveFloatPNG(path string, imageData [][]float64) {
	width, height := len(imageData), len(imageData[0])

	// get min and max values
	min, max := minMax(imageData)

	img := image.NewGray16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colVal := uint8((imageData[y][x] - min) / (max - min) * 255)
			img.Set(x, y, color.NRGBA{
				R: colVal,
				G: colVal,
				B: colVal,
				A: 255,
			})
		}
	}

	// save the new image
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func loadFloatPNG(path string) [][]float64 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var image [][]float64
	for y := 0; y < height; y++ {
		var row []float64
		for x := 0; x < width; x++ {
			row = append(row, float64(img.At(x, y).(color.Gray16).Y)/255.0)
		}
		image = append(image, row)
	}

	return image
}

func saveUint16AsIntPNG(path string, imageData [][]int) {
	width, height := len(imageData), len(imageData[0])

	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colVal := uint8(imageData[y][x])
			img.Set(x, y, color.Gray{Y: colVal})
		}
	}

	// save the new image
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func loadUint16AsIntPNG(path string) [][]int {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var image [][]int
	for y := 0; y < height; y++ {
		var row []int
		for x := 0; x < width; x++ {
			row = append(row, int(img.At(x, y).(color.Gray).Y))
		}
		image = append(image, row)
	}

	return image
}

const (
	DtUint16 DataType = iota
	DtFloat
)

type Image struct {
	size          int
	fill_value    float64
	dtype         DataType
	image         [][]float64
	is_image_load bool
}

func NewImage(size int, fill_value float64, dtype DataType) *Image {
	return &Image{
		size:          size,
		fill_value:    fill_value,
		dtype:         dtype,
		image:         new2DArr(size),
		is_image_load: false,
	}
}

func (i Image) LoadImage(path string) error {
	// TODO: Implement
	return nil
}

func (i Image) SaveImage(path string) error {
	// TODO: Implement
	return nil
}

func (i Image) FilterUniqueNumbersFrom2DArray() []float64 {
	var listOfUniqueNumbers []float64
	seenNumbers := make(map[float64]bool)
	for _, row := range i.image {
		for _, number := range row {
			if !seenNumbers[number] {
				listOfUniqueNumbers = append(listOfUniqueNumbers, number)
				seenNumbers[number] = true
			}
		}
	}
	return listOfUniqueNumbers
}

func FilterUniqueNumbersFrom2DArray(image [][]int) []int {
	var listOfUniqueNumbers []int
	seenNumbers := make(map[int]bool)
	for _, row := range image {
		for _, number := range row {
			if !seenNumbers[number] {
				listOfUniqueNumbers = append(listOfUniqueNumbers, number)
				seenNumbers[number] = true
			}
		}
	}
	return listOfUniqueNumbers
}

func (i Image) TransformImageToValidSoils(transformationList map[float64]float64) {
	for y := 0; y < i.size; y++ {
		for x := 0; x < i.size; x++ {
			i.image[y][x] = transformationList[i.image[y][x]]
		}
	}
}

func TransformImageToValidSoils(image [][]int, transformationList map[int]int) {
	for y := 0; y < len(image); y++ {
		for x := 0; x < len(image); x++ {
			image[y][x] = transformationList[image[y][x]]
		}
	}
}

/*

import imageio
import numpy as np
import os


class Image:
    """
    The image class is used as a container of the loaded and calculated maps. The image can be loaded and saved.
    """

    def __init__(self, size=None, fill_color=None, dtype=None):
        if size is None:
            size = 1
        self.size = size
        if fill_color is None:
            fill_color = 0
        self.fill_color = fill_color
        if dtype is None:
            dtype = np.uint16
        self.dtype = dtype
        self.image = np.full(shape=(size, size), fill_value=fill_color, dtype=self.dtype)
        self.is_image_loaded = False

    def __eq__(self, other):
        assert self.size == other.size, "The sizes of the images are not equal!"
        return (self.image == other.image).all()

    def load_image(self, path):
        self.image = imageio.imread(path)
        self.size = self.image.shape[0]
        # float values need special care. they were stored as integers and will be transformed back to float values.
        if self.dtype == np.float:
            new_image = np.full(shape=(self.size, self.size), fill_value=0.0, dtype=np.float)
            for y in range(self.size):
                for x in range(self.size):
                    new_image[x][y] = self.image[x][y] / 1000.0
            self.image = new_image
        self.is_image_loaded = True

    def save_image(self, path):
        if not os.path.exists(os.path.dirname(path)):
            os.makedirs(os.path.dirname(path))
        # float values need special care. they will be transformed to integers.
        if self.dtype == np.float:
            image = np.zeros(shape=(self.size, self.size), dtype=np.uint16)
            for y in range(self.size):
                for x in range(self.size):
                    image[x][y] = self.image[x][y] * 1000
            imageio.imwrite(path, image)
        else:
            imageio.imwrite(path, self.image)

    def filter_unique_numbers_from_2d_array(self):
        list_of_unique_numbers = []
        for row in self.image:
            for number in row:
                if number not in list_of_unique_numbers:
                    list_of_unique_numbers.append(number)
        print(list_of_unique_numbers)

    def transform_image_to_valid_soils(self, transformation_list):
        for y in range(self.size):
            for x in range(self.size):
                self.image[x][y] = transformation_list[self.image[x][y]]
*/
