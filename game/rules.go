// File game/rules.go
package game

import "trackLogicChess/player"

// CheckWin 检查指定颜色 col 是否在棋盘 b 上已连成 4 子。
// 返回 true 表示该颜色已在某一行、某一列或两条对角线上有 4 个连续的棋子。
func CheckWin(b *Board, col player.Color) bool {
	if col == player.Empty {
		return false
	}

	// 横向 & 纵向
	for i := 0; i < 4; i++ {
		// 横
		if b.cells[i][0] == col && b.cells[i][1] == col &&
			b.cells[i][2] == col && b.cells[i][3] == col {
			return true
		}
		// 纵
		if b.cells[0][i] == col && b.cells[1][i] == col &&
			b.cells[2][i] == col && b.cells[3][i] == col {
			return true
		}
	}

	// 主对角线
	if b.cells[0][0] == col && b.cells[1][1] == col &&
		b.cells[2][2] == col && b.cells[3][3] == col {
		return true
	}
	// 副对角线
	if b.cells[0][3] == col && b.cells[1][2] == col &&
		b.cells[2][1] == col && b.cells[3][0] == col {
		return true
	}
	return false
}
