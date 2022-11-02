package generator

import (
	"math"

	"github.com/aquilax/go-perlin"
	goimage "github.com/egor-erm/goimager/imager"
)

const (
	alpha  = 2.
	beta   = 2.
	n      = 5
	radius = 125
	step   = (math.Pi / 2) / radius
)

var blocks []string = []string{
	"#F5FC6D",
	"#A0FF50",
	"#52D52A",
	"#5E5E63",
	"#FCF7F2",
}

type Pos struct {
	X int
	Y int
}

type Island struct {
	Seed      int64
	Structure map[Pos]float64
}

func NewIsland(seed int64) *Island {
	return &Island{
		Seed:      seed,
		Structure: make(map[Pos]float64),
	}
}

func (island *Island) GenIsland() {
	island.genLandscape()
	island.genIslandShape()
}

func (island *Island) genLandscape() {
	p := perlin.NewPerlin(alpha, beta, n, island.Seed)
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			var max float64
			if int(math.Abs(float64(x))) >= int(math.Abs(float64(y))) {
				max = math.Abs(float64(x))
			} else {
				max = math.Abs(float64(y))
			}

			pos := Pos{
				X: x,
				Y: y,
			}

			h := math.Abs(p.Noise2D(float64(x)/65, float64(y)/65))
			h = math.Ceil(h*1000) / 1000

			island.Structure[pos] = h*math.Cos(max*step) + 0.036
		}
	}
}

func (island *Island) genIslandShape() {
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			var max float64
			if int(math.Abs(float64(x))) >= int(math.Abs(float64(y))) {
				max = math.Abs(float64(x))
			} else {
				max = math.Abs(float64(y))
			}

			k := 0.18 * math.Tan(max*step)
			pos := Pos{
				X: x,
				Y: y,
			}
			if h, ok := island.Structure[pos]; ok {
				if h-k < 0 {
					island.Structure[pos] = 0
				} else {
					island.Structure[pos] = h - k
				}

			} else {
				island.Structure[pos] = 0
			}
		}
	}
}

func (island *Island) SaveImage(name string) {
	img := goimage.NewWithCorners(name+".png", -radius, -radius, radius, radius)
	img.FillAllHex("#41B2D8")

	for k, value := range island.Structure {
		if value == 0 {
			continue
		}
		value += 0.16
		step := float64(1) / float64(len(blocks))
		otv := value / step
		var bl string

		if int(otv) < len(blocks) {
			bl = blocks[int(otv)]
		} else {
			bl = blocks[len(blocks)-1]
		}

		if value > 0 {
			img.SetHexPixel(k.X, k.Y, bl)
		}
	}

	img.Save()
}
