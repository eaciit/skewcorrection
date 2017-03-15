package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/gift"
	"github.com/eaciit/skewcorrection"
)

var DEBUG = false //true
var black = color.Gray{0}
var white = color.Gray{255}

func main() {
	filename := os.Args[1]
	outname := os.Args[2]

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
	counterRad, m, Y1, Y2 := skewcorrection.DetectRotation(image1)
	mRect := m.Bounds()
	//width := mRect.Max.X
	height := mRect.Max.Y

	//use this to use ImageMagick convert to rotate image
	//fmt.Println(optS, d, counterRad, Rad2Deg(counterRad))
	//rotateStr := strconv.FormatFloat(Rad2Deg(counterRad), 'f', -1, 64)
	//cmd := exec.Command("convert", filename, "-rotate", "-"+rotateStr, outname)
	//cmd.Run()

	//use this to rotate with float64 preccission
	//rotate.RotateImagePath(filename, outname, -rotationdetection.Rad2Deg(counterRad))

	//this example use gift filter
	//fmt.Println(rotationdetection.Rad2Deg(counterRad))
	//fmt.Println(rotationdetection.Rad2DegF32(counterRad))
	g := gift.New(
		gift.Rotate(skewcorrection.Rad2DegF32(counterRad), white, gift.CubicInterpolation),
	)
	toimg, _ := os.Create(outname)
	defer toimg.Close()
	dst := image.NewRGBA(g.Bounds(image1.Bounds()))
	g.Draw(dst, image1)
	if strings.HasSuffix(outname, ".png") {
		png.Encode(toimg, dst)
	} else if strings.HasSuffix(outname, ".jpg") || strings.HasSuffix(outname, ".jpeg") {
		jpeg.Encode(toimg, dst, &jpeg.Options{jpeg.DefaultQuality})
	}

	if DEBUG {
		for i := 0; i < height; i++ {
			//fmt.Println(m.At(Y1, i).(color.Gray).Y, m.At(Y2, i).(color.Gray).Y)
			m.Set(Y1, i, black)
			m.Set(Y2, i, black)
		}
		toimg, _ := os.Create("debug.jpg")
		defer toimg.Close()
		jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})
	}

	//fmt.Println(reflect.TypeOf(image1))
}
