package generator

import (
	"math"
	"strconv"

	"github.com/aquilax/go-perlin"
	goimage "github.com/egor-erm/goimager/imager"
)

var blocks []string = []string{
	"#F5FC6D",
	"#A0FF50",
	"#52D52A",
	"#5E5E63",
	"#FCF7F2",
}

const (
	alpha  = 2.
	beta   = 2.
	n      = 8
	radius = 120
	step   = (math.Pi / 2) / radius // radians
)

type Island struct {
	Seed   int64
	Blocks map[Location]string
}

type Location [2]int32

func NewIsland(seed int64) {
	island := &Island{
		Seed:   seed,
		Blocks: make(map[Location]string),
	}

	perlin := perlin.NewPerlin(alpha, beta, n, seed)

	var height_map map[Location]float64 = make(map[Location]float64)
	var max_height, min_height float64

	for x := int32(-radius); x <= radius; x++ {
		for y := int32(-radius); y <= radius; y++ {
			pos := Location{x, y}

			max := math.Hypot(float64(x), float64(y))

			height := math.Abs(perlin.Noise2D(float64(x)/60, float64(y)/60))

			height = 3 * height * math.Cos(max*step)

			m := 1.3
			height -= -1*m*math.Cos(max*step) + 1

			height = math.Round(height*1000) / 1000

			if height <= 0 {
				height = 0
			} else {
				max_height = math.Max(max_height, height)
				min_height = math.Min(min_height, height)
			}

			height_map[pos] = height
		}
	}

	img := goimage.NewWithCorners(strconv.Itoa(int(seed))+".png", -radius, -radius, radius, radius)
	img.FillAllHex("#329FC1") //вода

	height_diff := max_height - min_height
	for k, value := range height_map {
		if value == 0 {
			continue
		}

		block := int((value - min_height) / height_diff * float64(len(blocks)))
		if block >= len(blocks) {
			block = len(blocks) - 1
		}

		color := blocks[block]

		island.Blocks[Location{k[0], k[1]}] = color
	}

	for pos, bl := range island.Blocks {
		img.SetHexPixel(int(pos.X()), int(pos.Y()), bl)
	}

	img.Save()
}

func (loc *Location) X() int32 {
	return loc[0]
}

func (loc *Location) Y() int32 {
	return loc[1]
}
