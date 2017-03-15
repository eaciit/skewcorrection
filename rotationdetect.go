package skewcorrection

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"math"
)

var black = color.Gray{0}
var white = color.Gray{255}
var S = 30
var DEBUG = true

//input must GrayImage
func ThresholdImage(input draw.Image, threshold uint8) *image.Gray {
	bounds := input.Bounds()
	m := image.NewGray(bounds)
	width := bounds.Max.X
	height := bounds.Max.Y

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if input.At(j, i).(color.Gray).Y > threshold {
				m.Set(j, i, white)
			} else {
				m.Set(j, i, black)
			}
		}
	}
	return m
}

//return largest,earliest value and its index
func Max(list []int) (int, int) {
	biggest := list[0]
	bigIndex := 0
	for i, _ := range list {
		if list[i] > biggest {
			biggest = list[i]
			bigIndex = i
		}
	}
	return biggest, bigIndex
}

//return the product of all element in array
func ProdZero(list []int) int64 {
	for _, i := range list {
		if i == 0 {
			return 0
		}

	}
	return 1
}
func getRS(face *image.Gray, Y1 int, Y2 int, s int, d int) int {
	bounds := face.Bounds()
	Y := bounds.Max.Y
	sum := 0
	for y := S; y < Y-S; y++ {
		if face.At(Y1, y).(color.Gray).Y == 255 && face.At(Y2, y+s).(color.Gray).Y == 255 {

			sum++
		}
	}
	return sum
}
func Rad2Deg(rad float64) float64 {
	return rad * 180 / math.Pi
}
func Rad2DegF32(rad float64) float32 {
	return float32(rad * 180 / math.Pi)
}
func DetectRotation(image1 image.Image) (float64, draw.Image, int, int) {
	m := image.NewGray(image1.Bounds())
	mRect := m.Bounds()
	width := mRect.Max.X
	height := mRect.Max.Y
	draw.Draw(m, m.Bounds(), image1, image.ZP, draw.Src)
	m2 := ThresholdImage(m, 245)
	//fmt.Println()

	//toimg, _ := os.Create("new.jpg")
	//defer toimg.Close()
	//jpeg.Encode(toimg, m2, &jpeg.Options{jpeg.DefaultQuality})

	hist := []int{}
	for i := 0; i < width; i++ {
		hist = append(hist, 0)
	}
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if m2.At(i, j).(color.Gray).Y == black.Y {

				hist[i] += 1
			}
		}
	}
	histLeft := make([]int, len(hist)/2)
	histRight := make([]int, len(hist)/2)
	copy(histLeft, hist[:len(hist)/2])
	//histRight := hist[len(hist)/2:]
	copy(histRight, hist[len(hist)/2:])

	Y1 := 0
	Y2 := 0
	count := 0
	for true {
		_, mIndx := Max(histLeft)
		if mIndx < 5 {
			histLeft[mIndx] = 0
			count += 1
			continue
		}

		if ProdZero(hist[mIndx-5:mIndx+5]) > 0 {
			Y1 = mIndx
			break
		} else {
			histLeft[mIndx] = 0
			count += 1
			continue
		}
	}
	for true {
		_, mIndx := Max(histRight)
		if mIndx > len(histRight)-5 {
			histRight[mIndx] = 0
			continue
		}
		if ProdZero(hist[mIndx-5+len(histRight):mIndx+5+len(histRight)]) > 0 {
			Y2 = mIndx + len(histRight)
			break
		} else {
			histRight[mIndx] = 0
			continue
		}
	}
	d := Y2 - Y1
	fmt.Println(Y1, Y2)
	maxRS := 0
	optS := -100
	for s := -S; s <= S; s++ {
		curRS := getRS(m, Y1, Y2, s, d)
		//fmt.Println(curRS)
		if curRS > maxRS {
			maxRS = curRS
			optS = s
		}
	}
	counterRad := math.Atan(float64(optS) / float64(d))
	return counterRad, m, Y1, Y2
}
