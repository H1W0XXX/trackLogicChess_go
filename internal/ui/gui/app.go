package gui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"trackLogicChess/internal/game"
	"trackLogicChess/internal/player"
)

const (
	// AI 落子延迟
	aiDelay = 200 * time.Millisecond
)

// App 实现 ebiten.Game，管理输入、AI、动画与渲染
type App struct {
	state      *game.GameState
	useAI      bool
	anim       animator
	imgA, imgB *ebiten.Image

	// AI 延迟缓存
	pendingPrev *game.Board
	pendingRC   [2]int
	pendingTime time.Time
	lastAI      time.Time
}

// Update 处理输入、AI 触发和动画逻辑
func (a *App) Update() error {
	if a.state.IsGameOver() {
		return nil
	}
	now := time.Now()
	// 1) 动画进行中
	if a.anim.active {
		a.anim.Update()
		return nil
	}
	// 2) AI 回合（带延迟）
	if a.useAI && a.state.CurrentPlayer == player.White {
		// 第一次触发时记录 prev 和时间
		if a.pendingPrev == nil {
			a.pendingPrev = a.state.Board.Clone()
			mv := game.FindBestMoveDeep(a.state, 6)
			a.pendingRC = [2]int{mv.Row, mv.Col}
			a.pendingTime = now
		}
		// 延迟后执行落子+动画
		if now.Sub(a.pendingTime) >= aiDelay {
			_ = a.state.ApplyMove(a.pendingRC[0], a.pendingRC[1])

			a.anim.Start(
				a.pendingPrev,
				a.state.Board,
				a.state.DirOuter,
				a.state.DirInner,
				a.imgA,
				a.imgB,
			)
			a.pendingPrev = nil
		}
		return nil
	}
	// 3) 人类回合：点击立刻落子并动画
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		r := (y - boardOriginY) / cellSize
		c := (x - boardOriginX) / cellSize
		if r >= 0 && r < 4 && c >= 0 && c < 4 && a.state.Board.IsEmpty(r, c) {
			prev := a.state.Board.Clone()
			if err := a.state.ApplyMove(r, c); err == nil {
				
				a.anim.Start(
					prev,
					a.state.Board,
					a.state.DirOuter,
					a.state.DirInner,
					a.imgA,
					a.imgB,
				)
			}
		}
	}
	return nil
}

// Draw 渲染：动画中、AI延迟预览、默认渲染
func (a *App) Draw(screen *ebiten.Image) {
	// 1) 动画进行中
	if a.anim.active {
		a.anim.Draw(screen,
			a.state.Board,
			a.state.DirOuter,
			a.state.DirInner,
			a.imgA,
			a.imgB,
		)
		return
	}
	// 2) AI 延迟预览阶段，仅画原始棋盘
	if a.pendingPrev != nil {
		DrawBoard(screen,
			a.pendingPrev,
			a.imgA, a.imgB,
			a.state.DirOuter, a.state.DirInner,
		)
		return
	}
	// 3) 默认完整渲染
	DrawBoard(screen,
		a.state.Board,
		a.imgA, a.imgB,
		a.state.DirOuter, a.state.DirInner,
	)
}

// Layout 定义窗口尺寸
func (a *App) Layout(outW, outH int) (int, int) {
	size := 4 * cellSize
	return boardOriginX*2 + size, boardOriginY*2 + size
}
