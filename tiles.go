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

	ColorGrass

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
	ThemeBrogue
)

// Current theme.
const theme = ThemeNoir

type TileDrawer struct {
	drawer *tiles.Drawer
}

func inverseColor(c *image.Uniform) *image.Uniform {
	r, g, b, _ := c.RGBA()
	return image.NewUniform(rgba(uint8(255-r), uint8(255-g), uint8(255-b)))
}

func rgba(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}

// themeDefaults holds the base fg/bg color for each theme.
var themeDefaults = [2]struct{ fg, bg color.RGBA }{
	ThemeSelenized: {fg: rgba(0xad, 0xbc, 0xbc), bg: rgba(0x10, 0x3c, 0x48)},
	ThemeNoir:      {fg: rgba(80, 80, 80), bg: rgba(0, 0, 0)},
}

// fgTable maps a logical color to its per-theme foreground RGBA override.
// A zero color.RGBA (alpha == 0) means "use the theme default".
var fgTable = map[gruid.Color][3]color.RGBA{
	ColorPlayer:           {ThemeSelenized: rgba(0x46, 0x95, 0xf7), ThemeNoir: rgba(0xdb, 0xb3, 0x2d)},
	ColorBlood:            {ThemeSelenized: rgba(178, 3, 3), ThemeNoir: rgba(138, 3, 3)},
	ColorMonster:          {ThemeSelenized: rgba(0xfa, 0x57, 0x50), ThemeNoir: rgba(230, 0, 0)},
	ColorCorpse:           {ThemeSelenized: rgba(0xff, 0xa0, 0x30)},
	ColorLogPlayerAttack:  {ThemeSelenized: rgba(0x75, 0xb9, 0x38), ThemeNoir: rgba(0x75, 0xb9, 0x38)},
	ColorStatusHealthy:    {ThemeSelenized: rgba(0x75, 0xb9, 0x38), ThemeNoir: rgba(0x75, 0xb9, 0x38)},
	ColorLogMonsterAttack: {ThemeSelenized: rgba(0xed, 0x86, 0x49), ThemeNoir: rgba(230, 0, 0)},
	ColorStatusWounded:    {ThemeSelenized: rgba(0xed, 0x86, 0x49), ThemeNoir: rgba(230, 0, 0)},
	ColorLogSpecial:       {ThemeSelenized: rgba(0xf2, 0x75, 0xbe), ThemeNoir: rgba(0xdb, 0xb3, 0x2d)},
	ColorWater1:           {ThemeSelenized: rgba(148, 148, 255), ThemeNoir: rgba(148, 148, 255)},
	ColorFOV:              {ThemeNoir: rgba(200, 200, 200)},
	ColorTroll:            {ThemeNoir: rgba(20, 200, 20)},
	ColorHealthPotion:     {ThemeNoir: rgba(0xdb, 0xb3, 0x2d)},
	ColorScroll:           {ThemeNoir: rgba(0xdb, 0xb3, 0x2d)},
	ColorWater2:           {ThemeNoir: rgba(107, 107, 255)},
	ColorGrass:            {ThemeSelenized: rgba(0x44, 0x99, 0x33), ThemeNoir: rgba(0x44, 0x99, 0x33)},
}

// bgTable maps a logical color to its per-theme background RGBA override.
// A zero color.RGBA (alpha == 0) means "use the theme default".
var bgTable = map[gruid.Color][2]color.RGBA{
	ColorFOV:    {ThemeSelenized: rgba(0x18, 0x49, 0x56)},
	ColorBlood:  {ThemeSelenized: rgba(138, 3, 3), ThemeNoir: rgba(138, 3, 3)},
	ColorTarget: {ThemeSelenized: rgba(0x75, 0x75, 0x00), ThemeNoir: rgba(100, 100, 100)},
	ColorWater1: {ThemeSelenized: rgba(107, 107, 255), ThemeNoir: rgba(107, 107, 255)},
	ColorWater2: {ThemeNoir: rgba(148, 148, 255)},
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
