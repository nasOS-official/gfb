// SPDX-License-Identifier: MPL-2.0
// Package gfb provides a simple API for drawing graphics directly to a Linux framebuffer.
// It allows low-level manipulation of pixels, shapes, and colors using raw memory-mapped buffers.
//
// Typical usage:
//
//	var fb gfb.Gfb
//	fb.InitFb(0)
//	fb.ClearScreen()
//	fb.DrawRectangle(40, 200, 50, 100, 255, 0, 0)
//	fb.DrawCircle(300, 300, 100, 0, 255, 0)
//	fb.DrawLine(100, 400, 200, 500, 0, 0, 255)
//	fb.UpdateScreen()

package gfb

import (
	"os"
	"strconv"
	"strings"

	"github.com/crazy3lf/colorconv"
	"golang.org/x/sys/unix"
)

var fbMmap []byte

// Gfb represents a framebuffer device with resolution and pixel buffer.
type Gfb struct {
	ResX   int    // Horizontal resolution (pixels)
	ResY   int    // Vertical resolution (pixels)
	Fb     []byte // Framebuffer copy in memory
	fbSize int    // Size of framebuffer (bytes)
}

// Initializes the framebuffer device /dev/fb<num>.
// Reads resolution from /sys/class/graphics/fb<num>/virtual_size.
// Memory maps the framebuffer into process memory.
// Allocates an internal buffer (Fb) for drawing.
func (gfb *Gfb) InitFb(num int) {
	gfb.ResX, gfb.ResY = GetResolution("fb" + strconv.Itoa(num))
	gfb.fbSize = gfb.ResX * gfb.ResY * 4

	fbFile, err := os.OpenFile("/dev/fb"+strconv.Itoa(num), os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	fbMmap, err = unix.Mmap(
		int(fbFile.Fd()),
		0,
		gfb.fbSize,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED,
	)
	if err != nil {
		panic(err)
	}
	gfb.Fb = make([]byte, gfb.fbSize)
}

// Clears the internal framebuffer buffer (sets all pixels to black).
func (gfb *Gfb) ClearScreen() {
	copy(gfb.Fb, make([]byte, gfb.fbSize))
}

// GetResolution reads the resolution of a framebuffer device from sysfs.
// fbName should be like "fb0".
// Returns horizontal and vertical resolution in pixels.
func GetResolution(fbName string) (ResX, ResY int) {
	fbrel, _ := os.ReadFile("/sys/class/graphics/" + fbName + "/virtual_size")
	fbstr := string(fbrel[:len(fbrel)-1])
	fblist := strings.Split(fbstr, ",")
	ResX, _ = strconv.Atoi(fblist[0])
	ResY, _ = strconv.Atoi(fblist[1])
	return ResX, ResY
}

// SetPoint sets a pixel at (x, y) to the given RGB color.
// Pixels outside the screen bounds are ignored.
func (gfb *Gfb) SetPoint(x int, y int, r uint8, g uint8, b uint8) {
	if (gfb.ResX > x) && (x > 0) && (gfb.ResY > y) && (y > 0) {
		offset := (gfb.ResX*y + x) * 4
		gfb.Fb[offset] = b
		gfb.Fb[offset+1] = g
		gfb.Fb[offset+2] = r
		gfb.Fb[offset+3] = 0
	}
}

// SetPointHue sets a pixel at (x, y) using HSV color space.
// Hue, saturation, and value are converted to RGB before writing.
func (gfb *Gfb) SetPointHue(x int, y int, hue float64, saturation float64, value float64) {
	r, g, b, _ := colorconv.HSVToRGB(hue, saturation, value)
	gfb.SetPoint(x, y, r, g, b)
}

// DrawRectangle draws a filled rectangle between (xstart, ystart) and (xend, yend)
// with the given RGB color.
func (gfb *Gfb) DrawRectangle(xstart int, xend int, ystart int, yend int, r uint8, g uint8, b uint8) {
	lenght := yend - ystart
	for x := xstart; x <= xend; x++ {
		gfb.DrawVLine(x, ystart, lenght, r, g, b)
	}

}

// DrawTestRainbow draws a horizontal rainbow gradient within the given rectangle.
func (gfb *Gfb) DrawTestRainbow(xstart int, xend int, ystart int, yend int) {
	var n float64 = 0
	var add float64 = 360.0 / float64(xend-xstart)
	lenght := yend - ystart
	for x := xstart; x < xend; x++ {
		r, g, b, _ := colorconv.HSVToRGB(n, 0.9, 0.9)
		gfb.DrawVLine(x, ystart, lenght, r, g, b)
		n += add
	}
}

// GetPoint returns the RGB color of the pixel at (x, y).
// If the coordinates are out of bounds, returns black (0,0,0).
func (gfb *Gfb) GetPoint(x int, y int) (r, g, b uint8) {
	r, g, b = 0, 0, 0
	if (gfb.ResX > x) && (x > 0) && (gfb.ResY > y) && (y > 0) {
		offset := (gfb.ResX*y + x) * 4
		b = gfb.Fb[offset]
		g = gfb.Fb[offset+1]
		r = gfb.Fb[offset+2]
	}
	return r, g, b
}

// func DrawCircle(Fb []uint8, y_center int, x_center int, radius int, r uint8, g uint8, b uint8) {

// 	antiAliasRadius := 1.5

// 	for x := x_center - radius - int(antiAliasRadius); x <= x_center+radius+int(antiAliasRadius); x++ {
// 		for y := y_center - radius - int(antiAliasRadius); y <= y_center+radius+int(antiAliasRadius); y++ {

// 			distSquared := float64((x-x_center)*(x-x_center) + (y-y_center)*(y-y_center))
// 			dist := math.Sqrt(distSquared)
// 			radiusSq := float64(radius * radius)
// 			if distSquared <= radiusSq {
// 				SetPoint(Fb, x, y, r, g, b)
// 			} else {
// 				coverage := (dist - float64(radius)) / antiAliasRadius
// 				if coverage < 0 {
// 					coverage = 0
// 				}
// 				if coverage > 1 {
// 					coverage = 1
// 				}
// 				bg_r, bg_g, bg_b := GetPoint(Fb, x, y)
// 				// blended_r := uint8(float64(r)*(1.0-coverage) + float64(bg_r)*coverage)
// 				// blended_g := uint8(float64(g)*(1.0-coverage) + float64(bg_g)*coverage)
// 				// blended_b := uint8(float64(b)*(1.0-coverage) + float64(bg_b)*coverage)
// 				blended_r := bg_r/2 + r
// 				blended_g := bg_g/2 + g
// 				blended_b := bg_b/2 + b

// 				SetPoint(Fb, x, y, blended_r, blended_g, blended_b)
// 			}
// 		}
// 	}
// }

// blendPoint draws a pixel at (x, y) with RGB color and alpha blending
// against the existing buffer
func (gfb *Gfb) blendPoint(x, y int, r, g, b uint8, alpha uint8) {
	if (gfb.ResX > x) && (x > 0) && (gfb.ResY > y) && (y > 0) {
		offset := (gfb.ResX*y + x) * 4
		br := gfb.Fb[offset+2]
		bg := gfb.Fb[offset+1]
		bb := gfb.Fb[offset]

		inv := 255 - alpha

		gfb.Fb[offset] = uint8((int(bb)*int(inv) + int(b)*int(alpha)) / 255)
		gfb.Fb[offset+1] = uint8((int(bg)*int(inv) + int(g)*int(alpha)) / 255)
		gfb.Fb[offset+2] = uint8((int(br)*int(inv) + int(r)*int(alpha)) / 255)
		gfb.Fb[offset+3] = 0
	}
}

// DrawCircle draws a circle outline with the given center, radius, and RGB color.
func (gfb *Gfb) DrawCircle(y_center, x_center, radius int, r, g, b uint8) {
	x, y := radius, 0
	p := 1 - radius

	for x >= y {
		gfb.DrawHLine(x_center-x, y_center+y, x<<1, r, g, b)
		if y != 0 {
			gfb.DrawHLine(x_center-x, y_center-y, x<<1, r, g, b)
		}

		if x != y && y != 0 {
			gfb.DrawHLine(x_center-y, y_center+x, y<<1, r, g, b)
			gfb.DrawHLine(x_center-y, y_center-x, y<<1, r, g, b)
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

// DrawLine draws an anti-aliased line between (x0, y0) and (x1, y1)
// with the given RGB color.
func (gfb *Gfb) DrawLine(x0 int, x1 int, y0 int, y1 int, r uint8, g uint8, b uint8) {
	const M = 15
	const Ms = 1 << M
	const I = 0xff

	if x1 == x0 {
		gfb.DrawVLine(x0, y0, y1-y0, r, g, b)
		return
	} else if y1 == y0 {
		gfb.DrawHLine(x0, y0, x1-x0, r, g, b)
		return
	}

	dx := x1 - x0
	dy := y1 - y0
	d := (dy << M) / dx

	gfb.SetPoint(x0, y0, r, g, b)
	gfb.SetPoint(x1, y1, r, g, b)

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

		gfb.blendPoint(x, y0, r, g, b, c1)
		gfb.blendPoint(x, y0+1, r, g, b, c2)
	}
}

// DrawHLine draws a horizontal line at y, starting at xstart with given length and RGB color.
func (gfb *Gfb) DrawHLine(xstart, y, length int, r, g, b uint8) {
	if y < 0 || y >= gfb.ResY {
		return
	}
	if xstart < 0 {
		length += xstart
		xstart = 0
	}
	xend := xstart + length
	if xend >= gfb.ResX {
		xend = gfb.ResX - 1
	}

	offset := (y*gfb.ResX + xstart) * 4
	for x := xstart; x <= xend; x++ {
		gfb.Fb[offset] = b
		gfb.Fb[offset+1] = g
		gfb.Fb[offset+2] = r
		gfb.Fb[offset+3] = 0
		offset += 4
	}
}

// DrawVLine draws a vertical line at x, starting at ystart with given length and RGB color.
func (gfb *Gfb) DrawVLine(x, ystart, length int, r, g, b uint8) {
	if x < 0 || x >= gfb.ResX {
		return
	}
	if ystart < 0 {
		length += ystart
		ystart = 0
	}
	yend := ystart + length
	if yend >= gfb.ResY {
		yend = gfb.ResY - 1
	}

	offset := (ystart*gfb.ResX + x) * 4
	for y := ystart; y <= yend; y++ {
		gfb.Fb[offset] = b
		gfb.Fb[offset+1] = g
		gfb.Fb[offset+2] = r
		gfb.Fb[offset+3] = 0
		offset += gfb.ResX * 4
	}
}

// UpdateScreen copies the internal framebuffer buffer to the actual framebuffer memory.
// This must be called after drawing operations to update the display.
func (gfb *Gfb) UpdateScreen() {
	copy(fbMmap, gfb.Fb)
}
