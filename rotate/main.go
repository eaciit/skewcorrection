package rotate

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strings"
)

func Deg2Rad(p float64) float64 {
	return p * math.Pi / 180
}
func Rad2Deg(rad float64) float64 {
	return rad * 180 / math.Pi
}
func rotatePointFloat(start *image.Point, degree float64, x0, y0 int) (float64, float64) {

	rad := Deg2Rad(degree)
	p := math.Cos(rad)*float64(start.X-x0) - math.Sin(rad)*float64(start.Y-y0) + float64(x0)
	q := math.Sin(rad)*float64(start.X-x0) + math.Cos(rad)*float64(start.Y-y0) + float64(y0)
	return p, q
}
func rotatePoint(start *image.Point, degree float64, x0, y0 int) *image.Point {
	newPoint := &image.Point{}
	rad := Deg2Rad(degree)
	newPoint.X = int(math.Floor(math.Cos(rad)*float64(start.X-x0) - math.Sin(rad)*float64(start.Y-y0) + float64(x0)))
	newPoint.Y = int(math.Floor(math.Sin(rad)*float64(start.X-x0) + math.Cos(rad)*float64(start.Y-y0) + float64(y0)))
	return newPoint
}
func RGBAvg(colors []color.Color) color.Color {
	red := uint32(0)
	green := uint32(0)
	blue := uint32(0)
	alpha := uint8(255)
	for _, i := range colors {
		r, g, b, _ := i.RGBA()
		r8bit := uint8(r * 255 / 65535)
		g8bit := uint8(g * 255 / 65535)
		b8bit := uint8(b * 255 / 65535)
		//fmt.Println("RGB", r, g, b, r8bit, g8bit, b8bit)
		red = (red + uint32(r8bit))
		green = (green + uint32(g8bit))
		blue = (blue + uint32(b8bit))
	}
	//fmt.Println("preavg", red, green, blue)
	red = red / uint32(len(colors))

	green = green / uint32(len(colors))

	blue = blue / uint32(len(colors))
	//fmt.Println("prebuff", red, green, blue)
	if red > 255 {
		red = 255
	}
	if green > 255 {
		green = 255
	}
	if blue > 255 {
		blue = 255
	}
	//fmt.Println("postbuff", red, green, blue)
	return color.RGBA{uint8(red), uint8(green), uint8(blue), alpha}
}
func RotateImagePath(pathsrc, pathdest string, degree float64) {
	filename := pathsrc

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	image1, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//m := image.NewGray(image1.Bounds())
	mRect := image1.Bounds()
	kiriAtas := &image.Point{}
	kiriAtas.X = mRect.Min.X
	kiriAtas.Y = mRect.Min.Y

	kiriBawah := &image.Point{}
	kiriBawah.X = mRect.Min.X
	kiriBawah.Y = mRect.Max.Y

	kananAtas := &image.Point{}
	kananAtas.X = mRect.Max.X
	kananAtas.Y = mRect.Min.Y

	kananBawah := &image.Point{}
	kananBawah.X = mRect.Max.X
	kananBawah.Y = mRect.Max.Y

	x0 := mRect.Max.X / 2
	y0 := mRect.Max.Y / 2
	newPoint := [4]*image.Point{}
	newPoint[0] = rotatePoint(kiriAtas, degree, x0, y0)
	newPoint[1] = rotatePoint(kiriBawah, degree, x0, y0)
	newPoint[2] = rotatePoint(kananAtas, degree, x0, y0)
	newPoint[3] = rotatePoint(kananBawah, degree, x0, y0)
	minXNew := mRect.Max.X
	maxXNew := 0
	minYNew := mRect.Max.Y
	maxYNew := 0
	for _, i := range newPoint {
		if i.X < minXNew {
			minXNew = i.X
		}
		if i.X > maxXNew {
			maxXNew = i.X
		}
		if i.Y < minYNew {
			minYNew = i.Y
		}
		if i.Y > maxYNew {
			maxYNew = i.Y
		}
	}
	//fmt.Println(minXNew, maxXNew, minYNew, maxYNew)
	newRect := image.Rectangle{}
	newRect.Min.X = minXNew
	newRect.Min.Y = minYNew
	newRect.Max.X = maxXNew
	newRect.Max.Y = maxYNew
	m := image.NewRGBA(newRect)
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	checkmap := map[int]map[int][]color.Color{}
	for i := newRect.Min.Y - 1; i < newRect.Max.Y+1; i++ {
		checkmap[i] = map[int][]color.Color{}
		for j := newRect.Min.X - 2; j < newRect.Max.X+1; j++ {
			checkmap[i][j] = []color.Color{}
		}
	}
	for i := 0; i < mRect.Max.Y; i++ {
		for j := 0; j < mRect.Max.X; j++ {
			point := &image.Point{}
			point.X = j
			point.Y = i
			newX, newY := rotatePointFloat(point, degree, x0, y0)
			//fmt.Println(oldX, oldY)
			leftX := int(math.Floor(newX))
			rightX := int(math.Ceil(newX))
			upY := int(math.Floor(newY))
			downY := int(math.Ceil(newY))
			c1 := image1.At(j, i)
			//fmt.Println(newX, newY)
			checkmap[upY][leftX] = append(checkmap[upY][leftX], c1)
			checkmap[upY][rightX] = append(checkmap[upY][rightX], c1)
			checkmap[downY][leftX] = append(checkmap[downY][leftX], c1)
			checkmap[downY][rightX] = append(checkmap[downY][rightX], c1)
			//g := uint8((pix1G + pix2G + pix3G + pix4G) / 4)
			//b := uint8((pix1B + pix2B + pix3B + pix4B) / 4)
			//a := uint8((pix1A + pix2A + pix3A + pix4A) / 4)
			//fmt.Println(newPoint.Y)
			//checkmap[newPoint.Y][newPoint.X] = true
			//m.Set(j, i, color.RGBA{r, r, r, 255})
		}
	}
	for k, _ := range checkmap {
		for j, _ := range checkmap[k] {
			if len(checkmap[k][j]) > 0 {
				m.Set(j, k, RGBAvg(checkmap[k][j]))
			}

		}
	}
	/*
		for i := 0; i < newRect.Max.Y; i++ {
			for j := 0; j < newRect.Max.X; j++ {

			}
		}*/
	toimg, _ := os.Create(pathdest)
	defer toimg.Close()
	if strings.HasSuffix(pathdest, ".png") {
		png.Encode(toimg, m)
	} else if strings.HasSuffix(pathdest, ".jpg") || strings.HasSuffix(pathdest, ".jpeg") {
		jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})
	}

}
