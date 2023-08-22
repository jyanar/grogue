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

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/tiles"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const (
	ColorFOV gruid.Color = iota + 1
	ColorPlayer
	ColorMonster
	ColorCorpse
)

type TileDrawer struct {
	drawer *tiles.Drawer
}

const (
	ThemeSelenized = iota
	ThemeNoir
)

const theme = ThemeSelenized

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
		case ColorMonster:
			fg = image.NewUniform(color.RGBA{0xfa, 0x57, 0x50, 255})
		case ColorCorpse:
			fg = image.NewUniform(color.RGBA{0xff, 0xa0, 0x30, 255})
		}

		switch c.Style.Bg {
		case ColorFOV:
			bg = image.NewUniform(color.RGBA{0x18, 0x49, 0x56, 255})
		}

		return t.drawer.Draw(c.Rune, fg, bg)

	case ThemeNoir:
		fg = image.NewUniform(color.RGBA{100, 100, 100, 255})
		bg = image.NewUniform(color.RGBA{0x00, 0x00, 0x00, 255})

		switch c.Style.Fg {
		case ColorPlayer, ColorFOV:
			fg = image.NewUniform(color.RGBA{200, 200, 200, 255})
		case ColorMonster:
			fg = image.NewUniform(color.RGBA{230, 0, 0, 255})
		}

		return t.drawer.Draw(c.Rune, fg, bg)
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

	// IBM EGA
	// font, err := readTTF("assets/MxPlus_IBM_EGA_8x14.ttf")
	// if err != nil {
	// 	return nil, err
	// }

	// IBM MDA
	font, err := readTTF("assets/Mx437_IBM_MDA.ttf")
	if err != nil {
		return nil, err
	}

	// Retrieve the font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72 * 2,
		// Size: 12,
		// DPI:  72 * 2,
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
