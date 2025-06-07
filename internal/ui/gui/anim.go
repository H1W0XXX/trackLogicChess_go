package gui

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"trackLogicChess/internal/game"
	"trackLogicChess/internal/player"
)

const (
	rotateDur    = 2000 * time.Millisecond
	boardOriginX = 48
	boardOriginY = 48
	cellSize     = 64
)

// animator 控制一次旋转的关键帧
// 它会保存开始前后的棋子位置与图片
type animator struct {
	active     bool
	startAt    time.Time
	from, to   [4][4]tile
	imgA, imgB *ebiten.Image
}

type tile struct {
	img       *ebiten.Image // 棋子贴图指针
	initPos   image.Point   // 起点像素坐标
	targetPos image.Point   // 终点像素坐标
}

// Start 由 GUI 在完成逻辑旋转后调用，参数 imgA/imgB 对应两种棋子贴图
func (a *animator) Start(
	prev, next *game.Board,
	dirOuter, dirInner game.Direction,
	imgA, imgB *ebiten.Image,
) {
	// 清空
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			a.from[r][c] = tile{}
			a.to[r][c] = tile{}
		}
	}
	a.imgA, a.imgB = imgA, imgB

	// 外圈 + 内圈
	rings := []struct {
		coords [][2]int
		dir    game.Direction
	}{
		{ringCoords(0), dirOuter},
		{ringCoords(1), dirInner},
	}

	for _, ring := range rings {
		coords, dir := ring.coords, ring.dir
		n := len(coords)

		for dstIdx, rc := range coords {
			dstR, dstC := rc[0], rc[1]
			clr := next.Cell(dstR, dstC)
			if clr == player.Empty {
				continue // 目标格无子，跳过
			}
			img := chooseImage(clr, imgA, imgB)

			// 反推源 idx
			var srcIdx int
			if dir == game.Clockwise {
				srcIdx = (dstIdx - 1 + n) % n
			} else {
				srcIdx = (dstIdx + 1) % n
			}
			srcR, srcC := coords[srcIdx][0], coords[srcIdx][1]

			// 像素坐标
			x0, y0 := cellCenter(srcR, srcC, img)
			x1, y1 := cellCenter(dstR, dstC, img)

			// 写入条目（索引用“源格”）
			a.from[srcR][srcC] = tile{
				img:       img,
				initPos:   image.Pt(x0, y0),
				targetPos: image.Pt(x0, y0),
			}
			a.to[srcR][srcC] = tile{
				img:       img,
				initPos:   image.Pt(x0, y0),
				targetPos: image.Pt(x1, y1),
			}
		}
	}

	a.startAt = time.Now()
	a.active = true
}

// 计算格子 (r,c) 对应棋子左上角坐标
func cellCenter(r, c int, img *ebiten.Image) (x, y int) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	offX := (cellSize - w) / 2
	offY := (cellSize - h) / 2
	x = boardOriginX + c*cellSize + offX
	y = boardOriginY + r*cellSize + offY
	return
}

// chooseImage 根据 cellState 选择对应贴图
func chooseImage(clr player.Color, imgA, imgB *ebiten.Image) *ebiten.Image {
	if clr == player.Black {
		return imgB
	}
	return imgA
}
