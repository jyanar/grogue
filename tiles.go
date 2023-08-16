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

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/tiles"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
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

func (t *TileDrawer) GetImage(c gruid.Cell) image.Image {
	// Selenized theme
	fg := image.NewUniform(color.RGBA{0xad, 0xbc, 0xbc, 255})
	bg := image.NewUniform(color.RGBA{0x10, 0x3c, 0x48, 255})

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
}

func (t *TileDrawer) TileSize() gruid.Point {
	return t.drawer.Size()
}

func NewTileDrawer() (*TileDrawer, error) {
	t := &TileDrawer{}

	// Grab the monospace font TTF.
	font, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}

	// Retrieve the font face.
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 12,
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
