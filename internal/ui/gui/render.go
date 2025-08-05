package gui

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"trackLogicChess/internal/game"
	"trackLogicChess/internal/player"
)

// ─── 布局常量 ────────────────────────────────────────────────
const (
	boardSize = 4 * cellSize
)

var (
	lineColor = color.RGBA{0xff, 0xff, 0xff, 0xff} // 网格白
	arrowBlue = color.RGBA{0x00, 0x96, 0xff, 0xff} // 箭头蓝
)

// ────────────────────────────────────────────────────────────
// DrawBoard 绘制网格、两圈箭头、棋子。
// 方向由 GameState 决定：dirOuter/dirInner == game.Clockwise 表示顺时针。
// ────────────────────────────────────────────────────────────
func DrawBoard(screen *ebiten.Image, b *game.Board, imgA, imgB *ebiten.Image, dirOuter, dirInner game.Direction) {
	drawGrid(screen)
	drawRingArrows(screen, 0, dirOuter == game.Clockwise)
	drawRingArrows(screen, 1, dirInner == game.Clockwise)
	drawPieces(screen, b, imgA, imgB)
}

// ──────────────────────────────
// 1. 网格
// ──────────────────────────────
func drawGrid(screen *ebiten.Image) {
	for i := 0; i <= 4; i++ {
		x := float64(boardOriginX + i*cellSize)
		ebitenutil.DrawLine(screen, x, float64(boardOriginY), x, float64(boardOriginY+boardSize), lineColor)
		y := float64(boardOriginY + i*cellSize)
		ebitenutil.DrawLine(screen, float64(boardOriginX), y, float64(boardOriginX+boardSize), y, lineColor)
	}
}

// ──────────────────────────────
// 2. 两圈箭头
// ──────────────────────────────
func drawRingArrows(screen *ebiten.Image, ring int, clockwise bool) {
	coords := ringCoords(ring)
	n := len(coords)
	if n == 0 {
		return
	}
	// 每段的箭头：取连续两个坐标间的中点
	for i := 0; i < n; i++ {
		var from, to [2]int
		if clockwise {
			from = coords[i]
			to = coords[(i+1)%n]
		} else {
			from = coords[(i+1)%n]
			to = coords[i]
		}
		drawSegmentArrow(screen, from, to)
	}
}

// ringCoords 返回 ring=0 外圈 12 格，ring=1 内圈 4 格 的格子整数坐标
func ringCoords(ring int) [][2]int {
	if ring == 0 {
		return [][2]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {2, 0}, {1, 0}}
	}
	return [][2]int{{1, 1}, {1, 2}, {2, 2}, {2, 1}}
}

// 在两个格子中心连线的中点画蓝色箭头，长度≈cellSize/3
func drawSegmentArrow(screen *ebiten.Image, from, to [2]int) {
	//opts := &ebiten.DrawImageOptions{}
	fx := float64(boardOriginX + from[1]*cellSize + cellSize/2)
	fy := float64(boardOriginY + from[0]*cellSize + cellSize/2)
	tx := float64(boardOriginX + to[1]*cellSize + cellSize/2)
	ty := float64(boardOriginY + to[0]*cellSize + cellSize/2)

	mx := (fx + tx) / 2
	my := (fy + ty) / 2

	dx := tx - fx
	dy := ty - fy
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	ux, uy := dx/length, dy/length

	body := float64(cellSize) * 0.3
	x1 := mx - ux*body/2
	y1 := my - uy*body/2
	x2 := mx + ux*body/2
	y2 := my + uy*body/2
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1.0, arrowBlue, false)
	//ebitenutil.DrawLine(screen, x1, y1, x2, y2, arrowBlue)

	head := body * 0.4
	vx, vy := -uy, ux
	lx := x2 - ux*head + vx*head*0.6
	ly := y2 - uy*head + vy*head*0.6
	rx := x2 - ux*head - vx*head*0.6
	ry := y2 - uy*head - vy*head*0.6

	vector.StrokeLine(screen, float32(x2), float32(y2), float32(lx), float32(ly), 1.0, arrowBlue, false)
	vector.StrokeLine(screen, float32(x2), float32(y2), float32(rx), float32(ry), 1.0, arrowBlue, false)
	//ebitenutil.DrawLine(screen, x2, y2, lx, ly, arrowBlue)
	//ebitenutil.DrawLine(screen, x2, y2, rx, ry, arrowBlue)
}

// ──────────────────────────────
// 3. 棋子
// ──────────────────────────────
func drawPieces(screen *ebiten.Image, b *game.Board, imgA, imgB *ebiten.Image) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			var src *ebiten.Image
			switch b.Cell(r, c) {
			case player.Black:
				src = imgB
			case player.White:
				src = imgA
			default:
				continue
			}
			bw, bh := src.Bounds().Dx(), src.Bounds().Dy()
			offX := float64((cellSize - bw) / 2)
			offY := float64((cellSize - bh) / 2)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(boardOriginX+c*cellSize)+offX, float64(boardOriginY+r*cellSize)+offY)
			screen.DrawImage(src, op)
		}
	}
}

// Draw 绘制补间动画，结束后交给 DrawBoard 处理
func (a *animator) Draw(
	screen *ebiten.Image,
	current *game.Board,
	dirOuter, dirInner game.Direction,
	imgA, imgB *ebiten.Image,
) {
	if !a.active {
		// 动画结束，全量重绘
		DrawBoard(screen, current, imgA, imgB, dirOuter, dirInner)
		return
	}

	// 1) 画背景：网格 + 箭头
	drawGrid(screen)
	drawRingArrows(screen, 0, dirOuter == game.Clockwise)
	drawRingArrows(screen, 1, dirInner == game.Clockwise)

	// 2) 只绘制 those that move: initPos != targetPos
	t := float64(time.Since(a.startAt)) / float64(rotateDur)
	if t > 1 {
		t = 1
	}
	phase := t * t * (3 - 2*t) // smoothstep

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			tl := a.from[r][c]
			if tl.img == nil {
				continue
			}
			dst := a.to[r][c].targetPos
			if tl.initPos == dst {
				continue // 静止的跳过
			}
			x0, y0 := float64(tl.initPos.X), float64(tl.initPos.Y)
			x1, y1 := float64(dst.X), float64(dst.Y)
			x := x0 + (x1-x0)*phase
			y := y0 + (y1-y0)*phase

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(tl.img, op)
		}
	}

	// 3) 完成后跳回静态渲染
	if t >= 1 {
		a.active = false
	}
}

var firstFrame = true

// Update 每帧调用，检查是否结束动画
func (a *animator) Update() {
	if firstFrame {
		firstFrame = false
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum) // 已经进入事件循环，安全
	}
	if !a.active {
		return
	}
	if time.Since(a.startAt) >= rotateDur {
		a.active = false
	}
}
