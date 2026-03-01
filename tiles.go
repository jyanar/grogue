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

// themeDefaults holds the base fg/bg color for each theme.
var themeDefaults = [2]struct{ fg, bg color.RGBA }{
	ThemeSelenized: {fg: color.RGBA{0xad, 0xbc, 0xbc, 255}, bg: color.RGBA{0x10, 0x3c, 0x48, 255}},
	ThemeNoir:      {fg: color.RGBA{80, 80, 80, 255}, bg: color.RGBA{0x00, 0x00, 0x00, 255}},
}

// fgTable maps a logical color to its per-theme foreground RGBA override.
// A zero color.RGBA (alpha == 0) means "use the theme default".
var fgTable = map[gruid.Color][2]color.RGBA{
	ColorPlayer:           {ThemeSelenized: {0x46, 0x95, 0xf7, 255}, ThemeNoir: {0xdb, 0xb3, 0x2d, 255}},
	ColorBlood:            {ThemeSelenized: {178, 3, 3, 255}, ThemeNoir: {138, 3, 3, 255}},
	ColorMonster:          {ThemeSelenized: {0xfa, 0x57, 0x50, 255}, ThemeNoir: {230, 0, 0, 255}},
	ColorCorpse:           {ThemeSelenized: {0xff, 0xa0, 0x30, 255}},
	ColorLogPlayerAttack:  {ThemeSelenized: {0x75, 0xb9, 0x38, 255}, ThemeNoir: {0x75, 0xb9, 0x38, 255}},
	ColorStatusHealthy:    {ThemeSelenized: {0x75, 0xb9, 0x38, 255}, ThemeNoir: {0x75, 0xb9, 0x38, 255}},
	ColorLogMonsterAttack: {ThemeSelenized: {0xed, 0x86, 0x49, 255}, ThemeNoir: {230, 0, 0, 255}},
	ColorStatusWounded:    {ThemeSelenized: {0xed, 0x86, 0x49, 255}, ThemeNoir: {230, 0, 0, 255}},
	ColorLogSpecial:       {ThemeSelenized: {0xf2, 0x75, 0xbe, 255}, ThemeNoir: {0xdb, 0xb3, 0x2d, 255}},
	ColorWater1:           {ThemeSelenized: {148, 148, 255, 255}, ThemeNoir: {148, 148, 255, 255}},
	ColorFOV:              {ThemeNoir: {200, 200, 200, 255}},
	ColorTroll:            {ThemeNoir: {20, 200, 20, 255}},
	ColorHealthPotion:     {ThemeNoir: {0xdb, 0xb3, 0x2d, 255}},
	ColorScroll:           {ThemeNoir: {0xdb, 0xb3, 0x2d, 255}},
	ColorWater2:           {ThemeNoir: {107, 107, 255, 255}},
}

// bgTable maps a logical color to its per-theme background RGBA override.
// A zero color.RGBA (alpha == 0) means "use the theme default".
var bgTable = map[gruid.Color][2]color.RGBA{
	ColorFOV:    {ThemeSelenized: {0x18, 0x49, 0x56, 255}},
	ColorBlood:  {ThemeSelenized: {138, 3, 3, 255}, ThemeNoir: {138, 3, 3, 255}},
	ColorTarget: {ThemeSelenized: {0x75, 0x75, 0x00, 255}, ThemeNoir: {100, 100, 100, 255}},
	ColorWater1: {ThemeSelenized: {107, 107, 255, 255}, ThemeNoir: {107, 107, 255, 255}},
	ColorWater2: {ThemeNoir: {148, 148, 255, 255}},
}

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	d := themeDefaults[theme]
	fg := image.NewUniform(d.fg)
	bg := image.NewUniform(d.bg)
	if rgba, ok := fgTable[c.Style.Fg]; ok {
		if v := rgba[theme]; v.A != 0 {
			fg = image.NewUniform(v)
		}
	}
	if rgba, ok := bgTable[c.Style.Bg]; ok {
		if v := rgba[theme]; v.A != 0 {
			bg = image.NewUniform(v)
		}
	}
	if c.Style.Attrs == AttrReverse {
		fg, bg = bg, fg
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
