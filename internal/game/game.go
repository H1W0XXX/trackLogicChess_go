// File game/game.go
package game

import (
	"errors"
	"fmt"
	"trackLogicChess/internal/player"
)

// Direction 表示旋转方向。
// Clockwise 表示顺时针，CounterClockwise 表示逆时针。
type Direction int

const (
	Clockwise Direction = iota
	CounterClockwise
)

// GameState 保存当前游戏的状态，包括棋盘、当前玩家、固定的旋转方向、胜者和是否结束。
type GameState struct {
	Board         *Board       // 4×4 棋盘
	CurrentPlayer player.Color // 当前玩家 (Black 或 White)
	DirOuter      Direction    // 启动时固定的“外圈”旋转方向
	DirInner      Direction    // 启动时固定的“内圈”旋转方向
	Winner        player.Color // 胜者 (Black、White，或 Empty 表示平局/无胜者)
	GameOver      bool         // 游戏是否结束
}

// NewGame 新建一个 GameState，需要传入固定的外圈和内圈方向。
// 例如：
//
//	g := NewGame(Clockwise, CounterClockwise)
func NewGame(dirOuter, dirInner Direction) *GameState {
	return &GameState{
		Board:         NewBoard(),
		CurrentPlayer: player.Black,
		DirOuter:      dirOuter,
		DirInner:      dirInner,
		Winner:        player.Empty,
		GameOver:      false,
	}
}

// ApplyMove 在 (r,c) 位置落子，然后对外圈和内圈执行“固定方向”旋转。
// 旋转方向由 g.DirOuter 和 g.DirInner 决定，后续不允许修改。
// 落子完成并旋转后，再判断当前玩家是否连成 4 子。
// 参数 r,c 均在 0–3 范围内；如果出错（格子已占用或游戏已结束），返回非 nil 错误。
func (g *GameState) ApplyMove(r, c int) error {
	// 1. 检查游戏状态与目标格合法性
	if g.GameOver {
		return errors.New("game already over")
	}
	if !g.Board.IsEmpty(r, c) {
		return errors.New("cell not empty")
	}

	// 2. 在 (r,c) 放置当前玩家的棋子
	g.Board.Set(r, c, g.CurrentPlayer)

	// 打印落子后但未旋转前的棋盘（可选调试）
	//fmt.Println("\n--- 旋转前棋盘 ---")
	//g.Board.DebugPrint()

	// 3. 使用固定方向做旋转
	RotateOuter(g.Board, g.DirOuter)
	RotateInner(g.Board, g.DirInner)

	//打印旋转后的棋盘（可选调试）
	fmt.Println("\n--- 旋转后棋盘 ---")
	g.Board.DebugPrint()

	// 4. 旋转完成后，先判断当前玩家和对手是否同时连成 4
	selfWin := CheckWin(g.Board, g.CurrentPlayer)
	opp := player.Empty
	if g.CurrentPlayer == player.Black {
		opp = player.White
	} else {
		opp = player.Black
	}
	oppWin := CheckWin(g.Board, opp)

	if selfWin && oppWin {
		// 双方同时连 4，判平局
		g.Winner = player.Empty
		g.GameOver = true
		return nil
	}
	if selfWin {
		// 只有自己连接 4
		g.Winner = g.CurrentPlayer
		g.GameOver = true
		return nil
	}
	if oppWin {
		// 只有对手连接 4（因为旋转导致自杀）
		g.Winner = opp
		g.GameOver = true
		return nil
	}

	// 5. 如果棋盘已满且无人连成 4，则平局
	if g.isBoardFull() {
		g.GameOver = true
		return nil
	}

	// 6. 切换到下一玩家
	if g.CurrentPlayer == player.Black {
		g.CurrentPlayer = player.White
	} else {
		g.CurrentPlayer = player.Black
	}
	return nil
}

// isBoardFull 判断棋盘是否已满
func (g *GameState) isBoardFull() bool {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if g.Board.IsEmpty(r, c) {
				return false
			}
		}
	}
	return true
}

// IsGameOver 返回游戏是否结束。
func (g *GameState) IsGameOver() bool {
	return g.GameOver
}

// WinnerColor 返回本局胜者颜色；如果平局或未结束，返回 player.Empty。
func (g *GameState) WinnerColor() player.Color {
	return g.Winner
}
