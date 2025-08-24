// Copyright (C) Egor
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
// GNU Lesser General Public Licence is available at
// http://www.gnu.org/copyleft/lesser.html
package gfb

import (
	"os"
	"strconv"
	"strings"

	"github.com/crazy3lf/colorconv"
	"golang.org/x/sys/unix"
)

var resX, resY int = GetResolution("fb0")
var fbSize int = resX * resY * 4

var fbMmap []byte

func InitFb() []byte {
	fbFile, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	fbMmap, err = unix.Mmap(
		int(fbFile.Fd()),
		0,
		fbSize,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED,
	)
	if err != nil {
		panic(err)
	}
	return make([]byte, fbSize)
}

func ClearScreen(fb []byte) {
	copy(fb, make([]byte, fbSize))
}

func GetResolution(fbName string) (resX, resY int) {
	fbrel, _ := os.ReadFile("/sys/class/graphics/" + fbName + "/virtual_size")
	fbstr := string(fbrel[:len(fbrel)-1])
	fblist := strings.Split(fbstr, ",")
	resX, _ = strconv.Atoi(fblist[0])
	resY, _ = strconv.Atoi(fblist[1])
	return resX, resY
}

func SetPoint(fb []uint8, x int, y int, r uint8, g uint8, b uint8) {
	if (resX > x) && (x > 0) && (resY > y) && (y > 0) {
		offset := (resX*y + x) * 4
		fb[offset] = b
		fb[offset+1] = g
		fb[offset+2] = r
		fb[offset+3] = 0
	}
}
func SetPointHue(fb []uint8, x int, y int, hue float64, saturation float64, value float64) {
	r, g, b, _ := colorconv.HSVToRGB(hue, saturation, value)
	SetPoint(fb, x, y, r, g, b)
}

func DrawRectangle(fb []uint8, xstart int, xend int, ystart int, yend int, r uint8, g uint8, b uint8) {
	lenght := yend - ystart
	for x := xstart; x <= xend; x++ {
		DrawVLine(fb, x, ystart, lenght, r, g, b)
	}

}

func DrawTestRainbow(fb []uint8, xstart int, xend int, ystart int, yend int) {
	var n float64 = 0
	var add float64 = 360.0 / float64(xend-xstart)
	lenght := yend - ystart
	for x := xstart; x < xend; x++ {
		r, g, b, _ := colorconv.HSVToRGB(n, 0.9, 0.9)
		DrawVLine(fb, x, ystart, lenght, r, g, b)
		n += add
	}
}

func GetPoint(fb []uint8, x int, y int) (r, g, b uint8) {
	r, g, b = 0, 0, 0
	if (resX > x) && (x > 0) && (resY > y) && (y > 0) {
		offset := (resX*y + x) * 4
		b = fb[offset]
		g = fb[offset+1]
		r = fb[offset+2]
	}
	return r, g, b
}

// func DrawCircle(fb []uint8, y_center int, x_center int, radius int, r uint8, g uint8, b uint8) {

// 	antiAliasRadius := 1.5

// 	for x := x_center - radius - int(antiAliasRadius); x <= x_center+radius+int(antiAliasRadius); x++ {
// 		for y := y_center - radius - int(antiAliasRadius); y <= y_center+radius+int(antiAliasRadius); y++ {

// 			distSquared := float64((x-x_center)*(x-x_center) + (y-y_center)*(y-y_center))
// 			dist := math.Sqrt(distSquared)
// 			radiusSq := float64(radius * radius)
// 			if distSquared <= radiusSq {
// 				SetPoint(fb, x, y, r, g, b)
// 			} else {
// 				coverage := (dist - float64(radius)) / antiAliasRadius
// 				if coverage < 0 {
// 					coverage = 0
// 				}
// 				if coverage > 1 {
// 					coverage = 1
// 				}
// 				bg_r, bg_g, bg_b := GetPoint(fb, x, y)
// 				// blended_r := uint8(float64(r)*(1.0-coverage) + float64(bg_r)*coverage)
// 				// blended_g := uint8(float64(g)*(1.0-coverage) + float64(bg_g)*coverage)
// 				// blended_b := uint8(float64(b)*(1.0-coverage) + float64(bg_b)*coverage)
// 				blended_r := bg_r/2 + r
// 				blended_g := bg_g/2 + g
// 				blended_b := bg_b/2 + b

// 				SetPoint(fb, x, y, blended_r, blended_g, blended_b)
// 			}
// 		}
// 	}
// }

func blendPoint(fb []uint8, x, y int, r, g, b uint8, alpha uint8) {
	if (resX > x) && (x > 0) && (resY > y) && (y > 0) {
		offset := (resX*y + x) * 4
		br := fb[offset+2]
		bg := fb[offset+1]
		bb := fb[offset]

		inv := 255 - alpha

		fb[offset] = uint8((int(bb)*int(inv) + int(b)*int(alpha)) / 255)
		fb[offset+1] = uint8((int(bg)*int(inv) + int(g)*int(alpha)) / 255)
		fb[offset+2] = uint8((int(br)*int(inv) + int(r)*int(alpha)) / 255)
		fb[offset+3] = 0
	}
}

func DrawCircle(fb []uint8, y_center, x_center, radius int, r, g, b uint8) {
	x, y := radius, 0
	p := 1 - radius

	for x >= y {
		DrawHLine(fb, x_center-x, y_center+y, x<<1, r, g, b)
		if y != 0 {
			DrawHLine(fb, x_center-x, y_center-y, x<<1, r, g, b)
		}

		if x != y && y != 0 {
			DrawHLine(fb, x_center-y, y_center+x, y<<1, r, g, b)
			DrawHLine(fb, x_center-y, y_center-x, y<<1, r, g, b)
		}

		y++
		if p <= 0 {
			p += y<<1 + 1
		} else {
			x--
			p += y<<1 - x<<1 + 1
		}
	}
}

func DrawLine(fb []uint8, x0 int, x1 int, y0 int, y1 int, r uint8, g uint8, b uint8) {
	const M = 15
	const Ms = 1 << M
	const I = 0xff

	if x1 == x0 {
		DrawVLine(fb, x0, y0, y1-y0, r, g, b)
		return
	} else if y1 == y0 {
		DrawHLine(fb, x0, y0, x1-x0, r, g, b)
		return
	}

	dx := x1 - x0
	dy := y1 - y0
	d := (dy << M) / dx

	SetPoint(fb, x0, y0, r, g, b)
	SetPoint(fb, x1, y1, r, g, b)

	D := 0
	for x := x0; x <= x1; x++ {
		D += d
		if D >= Ms {
			D -= Ms
			y0++
		}

		v := (D >> 7) & I
		c1 := uint8(I - v)
		c2 := uint8(v)

		blendPoint(fb, x, y0, r, g, b, c1)
		blendPoint(fb, x, y0+1, r, g, b, c2)
	}
}

func DrawHLine(fb []uint8, xstart, y, length int, r, g, b uint8) {
	if y < 0 || y >= resY {
		return
	}
	if xstart < 0 {
		length += xstart
		xstart = 0
	}
	xend := xstart + length
	if xend >= resX {
		xend = resX - 1
	}

	offset := (y*resX + xstart) * 4
	for x := xstart; x <= xend; x++ {
		fb[offset] = b
		fb[offset+1] = g
		fb[offset+2] = r
		fb[offset+3] = 0
		offset += 4
	}
}

func DrawVLine(fb []uint8, x, ystart, length int, r, g, b uint8) {
	if x < 0 || x >= resX {
		return
	}
	if ystart < 0 {
		length += ystart
		ystart = 0
	}
	yend := ystart + length
	if yend >= resY {
		yend = resY - 1
	}

	offset := (ystart*resX + x) * 4
	for y := ystart; y <= yend; y++ {
		fb[offset] = b
		fb[offset+1] = g
		fb[offset+2] = r
		fb[offset+3] = 0
		offset += resX * 4
	}
}

func UpdateScreen(fb []uint8) {
	copy(fbMmap, fb)
}

// //
// func main() {

// 	fb := InitFb()
// 	// drawTestRainbow(fb, (resX-resY)/2, resY+((resX-resY)/2), 0, resY)
// 	DrawRectangle(fb, 40, 500, 50, 100, 0, 255, 26)
// 	DrawLine(fb, 80, 800, 50, 100, 0, 255, 26)
// 	DrawTestRainbow(fb, 50, 320, 50, 320)
// 	DrawCircle(fb, 600, 600, 300, 255, 255, 0)
// 	DrawCircle(fb, 70, 70, 50, 255, 0, 0)
// 	DrawCircle(fb, 70, 120, 50, 0, 255, 0)
// 	DrawCircle(fb, 70, 170, 50, 0, 0, 255)
// 	UpdateScreen(fb)

//		os.Exit(0)
//	}
//
