// File game/board.go
package game

import (
	"fmt"
	"strings"

	"trackLogicChess/player"
)

// Board 表示一个 4×4 棋盘，cells[r][c] 存储第 r 行、第 c 列的棋子颜色。
// 值为 player.Empty、player.Black 或 player.White。
type Board struct {
	cells [4][4]player.Color
}

// NewBoard 返回一个全空（所有格子值为 player.Empty）的 4×4 棋盘。
func NewBoard() *Board {
	b := &Board{}
	// Go 的数组会默认初始化为 0，本项目中 player.Empty == 0，
	// 因此这里不需要显式赋值。若想确保，可如下写法：
	// for r := 0; r < 4; r++ {
	//     for c := 0; c < 4; c++ {
	//         b.cells[r][c] = player.Empty
	//     }
	// }
	return b
}

// IsEmpty 返回 (r, c) 位置是否为空（player.Empty）。
// 若 r 或 c 越界，则返回 false。
func (b *Board) IsEmpty(r, c int) bool {
	if r < 0 || r >= 4 || c < 0 || c >= 4 {
		return false
	}
	return b.cells[r][c] == player.Empty
}

// Set 在 (r, c) 位置放置一个颜色为 col 的棋子。
// 不会做越界检查或重复落子检查，调用方需自行保证合法性。
func (b *Board) Set(r, c int, col player.Color) {
	b.cells[r][c] = col
}

// Get 返回 (r, c) 位置的棋子颜色。
// 若越界，返回 player.Empty（通常应由调用处保证索引合法）。
func (b *Board) Get(r, c int) player.Color {
	if r < 0 || r >= 4 || c < 0 || c >= 4 {
		return player.Empty
	}
	return b.cells[r][c]
}

// String 将棋盘渲染为多行字符串，可用于终端打印调试。
// 使用 “.” 表示空格，
//
//	“●” 表示黑子 (player.Black)，
//	“○” 表示白子 (player.White)。
func (b *Board) String() string {
	var sb strings.Builder
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			switch b.cells[r][c] {
			case player.Black:
				sb.WriteString("○ ")
			case player.White:
				sb.WriteString("● ")
			default:
				sb.WriteString(". ")
			}
		}
		// 去掉行末最后一个空格，然后换行
		line := strings.TrimRight(sb.String(), " ")
		sb.Reset()
		sb.WriteString(line)
		if r < 3 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// DebugPrint 简单地将棋盘输出到标准输出，方便调试。
func (b *Board) DebugPrint() {
	fmt.Println(b.String())
}
