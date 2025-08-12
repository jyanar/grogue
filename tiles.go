// Implements the TileManager interface for gruid-sdl:
//
// type TileManager interface {
//   GetImage(gruid.Cell) image.Image
//   TileSize() gruid.Point
// }
//

package main

import (
	"image"
	"image/color"
	"os"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/tiles"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

// Available colors. These are set to appropriate RGB values by the theme.
const (
	ColorNone gruid.Color = iota + 1
	ColorFOV
	ColorTarget
	ColorPlayer
	ColorMonster
	ColorTroll

	ColorCorpse
	ColorHealthPotion
	ColorScroll
	ColorBlood
	ColorWater1
	ColorWater2

	ColorLog
	ColorLogPlayerAttack
	ColorLogMonsterAttack
	ColorLogSpecial
	ColorStatusHealthy
	ColorStatusWounded
)

// const (
// 	ColorFOV     = "ColorFOV"
// 	ColorTarget  = "ColorTarget"
// 	ColorPlayer  = "ColorPlayer"
// 	ColorMonster = "ColorMonster"
// 	ColorTroll   = "ColorTroll"

// 	ColorCorpse       = "ColorCorpse"
// 	ColorHealthPotion = "ColorHealthPotion"
// 	ColorScroll       = "ColorScroll"
// 	ColorBlood        = "ColorBlood"
// 	ColorWater1       = "ColorWater1"
// 	ColorWater2       = "ColorWater2"

// 	ColorLog              = "ColorLog"
// 	ColorLogPlayerAttack  = "ColorLogPlayerAttack"
// 	ColorLogMonsterAttack = "ColorLogMonsterAttack"
// 	ColorLogSpecial       = "ColorLogSpecial"
// 	ColorStatusHealthy    = "ColorStatusHealthy"
// 	ColorStatusWounded    = "ColorStatusWounded"
// )

// A list of available themes.
const (
	ThemeSelenized = iota
	ThemeNoir
)

// Current theme.
const theme = ThemeNoir

type TileDrawer struct {
	drawer *tiles.Drawer
}

func inverseColor(c *image.Uniform) *image.Uniform {
	r, g, b, _ := c.RGBA()
	return image.NewUniform(color.RGBA{uint8(255 - r), uint8(255 - g), uint8(255 - b), 255})
}

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	fg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 255})
	bg := image.NewUniform(color.RGBA{0x00, 0x00, 0x00, 255})
	switch theme {

	case ThemeSelenized:
		fg = image.NewUniform(color.RGBA{0xad, 0xbc, 0xbc, 255})
		bg = image.NewUniform(color.RGBA{0x10, 0x3c, 0x48, 255})

		switch c.Style.Fg {
		case ColorPlayer:
			fg = image.NewUniform(color.RGBA{0x46, 0x95, 0xf7, 255})
		case ColorBlood:
			fg = image.NewUniform(color.RGBA{178, 3, 3, 255})
		case ColorMonster:
			fg = image.NewUniform(color.RGBA{0xfa, 0x57, 0x50, 255})
		case ColorCorpse:
			fg = image.NewUniform(color.RGBA{0xff, 0xa0, 0x30, 255})
		case ColorLogPlayerAttack, ColorStatusHealthy:
			fg = image.NewUniform(color.RGBA{0x75, 0xb9, 0x38, 255})
		case ColorLogMonsterAttack, ColorStatusWounded:
			fg = image.NewUniform(color.RGBA{0xed, 0x86, 0x49, 255})
		case ColorLogSpecial:
			fg = image.NewUniform(color.RGBA{0xf2, 0x75, 0xbe, 255})
		case ColorWater1:
			fg = image.NewUniform(color.RGBA{148, 148, 255, 255})
		}

		switch c.Style.Bg {
		case ColorFOV:
			bg = image.NewUniform(color.RGBA{0x18, 0x49, 0x56, 255})
		case ColorBlood:
			bg = image.NewUniform(color.RGBA{138, 3, 3, 255})
		case ColorTarget:
			bg = image.NewUniform(color.RGBA{0x75, 0x75, 0x00, 255})
		case ColorWater1:
			bg = image.NewUniform(color.RGBA{107, 107, 255, 255})
		}

	case ThemeNoir:
		fg = image.NewUniform(color.RGBA{80, 80, 80, 255})
		bg = image.NewUniform(color.RGBA{0x00, 0x00, 0x00, 255})

		switch c.Style.Fg {
		case ColorPlayer:
			fg = image.NewUniform(color.RGBA{0xdb, 0xb3, 0x2d, 255})
		case ColorBlood:
			fg = image.NewUniform(color.RGBA{138, 3, 3, 255})
		case ColorFOV:
			fg = image.NewUniform(color.RGBA{200, 200, 200, 255})
		case ColorMonster:
			fg = image.NewUniform(color.RGBA{230, 0, 0, 255})
		case ColorTroll:
			fg = image.NewUniform(color.RGBA{20, 200, 20, 255})
		case ColorHealthPotion:
			fg = image.NewUniform(color.RGBA{0xdb, 0xb3, 0x2d, 255})
		case ColorScroll:
			fg = image.NewUniform(color.RGBA{0xdb, 0xb3, 0x2d, 255})
		case ColorLogPlayerAttack, ColorStatusHealthy:
			fg = image.NewUniform(color.RGBA{0x75, 0xb9, 0x38, 255})
		case ColorLogMonsterAttack, ColorStatusWounded:
			fg = image.NewUniform(color.RGBA{230, 0, 0, 255})
		case ColorLogSpecial:
			fg = image.NewUniform(color.RGBA{0xdb, 0xb3, 0x2d, 255})
		case ColorWater1:
			fg = image.NewUniform(color.RGBA{148, 148, 255, 255})
		case ColorWater2:
			fg = image.NewUniform(color.RGBA{107, 107, 255, 255})
		}

		switch c.Style.Bg {
		case ColorBlood:
			bg = image.NewUniform(color.RGBA{138, 3, 3, 255})
		case ColorTarget:
			bg = image.NewUniform(color.RGBA{100, 100, 100, 255})
		case ColorWater1:
			bg = image.NewUniform(color.RGBA{107, 107, 255, 255})
		case ColorWater2:
			bg = image.NewUniform(color.RGBA{148, 148, 255, 255})
		}

		switch c.Style.Attrs { // If you want to invert fg and bg
		case AttrReverse:
			fg, bg = bg, fg
		}

	}
	return t.drawer.Draw(c.Rune, fg, bg)
}

func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}

func readTTF(filepath string) (*sfnt.Font, error) {
	fontBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	font, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}

func NewTileDrawer() (*TileDrawer, error) {
	t := &TileDrawer{}

	// Grab the monospace font TTF.

	// // GoMono
	// font, err := opentype.Parse(gomono.TTF)
	// if err != nil {
	// 	return nil, err
	// }

	// // IBM EGA
	// font, err := readTTF("assets/MxPlus_IBM_EGA_8x14.ttf")
	// if err != nil {
	// 	return nil, err
	// }

	// IBM MDA
	font, err := opentype.Parse(ibm_mda)
	if err != nil {
		return nil, err
	}

	// Retrieve the font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72 * 2,
	})
	if err != nil {
		return nil, err
	}

	// Create new drawer for tiles using the face. Note that we could use
	// multiple faces (e.g. italic/bold/etc) -- in that case we would simply
	// define drawers for those as well and call the appropriate one in the
	// GetImage method.
	t.drawer, err = tiles.NewDrawer(face)
	if err != nil {
		return nil, err
	}
	return t, nil
}
