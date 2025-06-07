package assets

import (
	"image"
	_ "image/png"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// LoadPNG 创建 *ebiten.Image
func LoadPNG(path string) *ebiten.Image {
	f, err := Images.Open(path) // 例如 "images/marbleA.png"
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f.(io.Reader))
	if err != nil {
		log.Fatalf("decode %s: %v", path, err)
	}
	return ebiten.NewImageFromImage(img)
}
