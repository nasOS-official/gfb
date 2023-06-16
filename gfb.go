package gfb

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/chai2010/webp"

	"github.com/crazy3lf/colorconv"
)

var resX, resY int = getResolution("fb0")

func initFb() []uint8 {
	return make([]uint8, (resX * resY * 4))
}

func getResolution(fbName string) (resX, resY int) {
	fbrel, _ := ioutil.ReadFile("/sys/class/graphics/" + fbName + "/virtual_size")
	fmt.Println(len(fbrel))
	fbstr := string(fbrel[:len(fbrel)-1])
	fblist := strings.Split(fbstr, ",")
	resX, _ = strconv.Atoi(fblist[0])
	resY, _ = strconv.Atoi(fblist[1])
	return resX, resY
}

func setpoint(fb []uint8, x int, y int, r uint8, g uint8, b uint8) []uint8 {
	fb[(resX*x+y)*4] = b
	fb[(resX*x+y)*4+1] = g
	fb[(resX*x+y)*4+2] = r
	fb[(resX*x+y)*4+3] = 0
	return fb
}
func setpointHue(fb []uint8, x int, y int, hue float64, saturation float64, value float64) []uint8 {
	r, g, b, _ := colorconv.HSVToRGB(hue, saturation, value)
	return setpoint(fb, x, y, r, g, b)
}
func drawRectangle(fb []uint8, xstart int, xend int, ystart int, yend int, r uint8, g uint8, b uint8) {
	for y := ystart; y <= yend; y++ {
		for x := xstart; x <= xend; x++ {
			fb = setpoint(fb, y, x, r, g, b)
		}
	}

}

func drawTestRainbow(fb []uint8, xstart int, xend int, ystart int, yend int) {
	var n float64 = 1

	for y := ystart; y < yend; y++ {
		for x := xstart; x < xend; x++ {
			fb = setpointHue(fb, x, y, n/(float64(yend-ystart)*3), 0.9, 0.9)
			n++
		}
	}
}

func drawCircle(fb []uint8, y_center int, x_center int, radius int, r uint8, g uint8, b uint8) {
	for y := y_center - radius; y <= y_center+radius; y++ {
		for x := x_center - radius; x <= x_center+radius; x++ {
			if (x-x_center)*(x-x_center)+(y-y_center)*(y-y_center) <= radius*radius {
				fb = setpoint(fb, y, x, r, g, b)
			}
		}
	}

}

func drawLine(fb []uint8, xstart int, xend int, ystart int, yend int, r uint8, g uint8, b uint8) {
	// Calculate the distance and direction of the line
	dx := xend - xstart
	dy := yend - ystart
	dist := math.Sqrt(float64(dx*dx + dy*dy))

	// Draw the line by setting the color of each pixel along its path
	for t := 0.0; t <= 1.0; t += 1.0 / dist {
		x := int(float64(xstart) + t*float64(dx))
		y := int(float64(ystart) + t*float64(dy))
		fb = setpoint(fb, x, y, r, g, b)
	}

}

//lint:ignore U1000 Ignore unused function temporarily for debugging

func writeWebp(data []uint8, width, height int, filepath string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 4
			b := data[index]
			g := data[index+1]
			r := data[index+2]
			a := uint8(255)
			img.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, nil); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath, buf.Bytes(), 0666); err != nil {
		return err
	}
	return nil
}
func update_screen(fb []uint8) {
	_ = os.WriteFile("/dev/fb0", fb, 0644)
}

// func main() {
// 	fmt.Println("gfbos -  Go FrameBuffer")
// 	fmt.Println("Current Screen resolution is " + strconv.Itoa(resX) + "x" + strconv.Itoa(resY) + "px")
// 	fb := initFb()
// 	// drawTestRainbow(fb, (resX-resY)/2, resY+((resX-resY)/2), 0, resY)
// 	drawRectangle(fb, 40, 500, 50, 100, 0, 255, 26)
// 	drawLine(fb, 80, 800, 50, 100, 0, 255, 26)
// 	drawTestRainbow(fb, 50, 320, 50, 320)
// 	drawCircle(fb, 600, 600, 300, 255, 255, 0)
// 	drawCircle(fb, 70, 70, 50, 255, 0, 0)
// 	drawCircle(fb, 70, 120, 50, 0, 255, 0)
// 	drawCircle(fb, 70, 170, 50, 0, 0, 255)
// 	update_screen(fb)
// 	// _ = writeWebp(fb, resX, resY, "./test.webp")

// 	os.Exit(0)
// }
