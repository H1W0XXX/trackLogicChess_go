// File game/ring.go
package game

import (
	"trackLogicChess/internal/player"
)

// RotateOuter 对 4×4 棋盘的外圈 12 个格子执行“环移”一格操作。
// 外圈坐标（顺时针顺序）为：
//
//	(0,0) → (0,1) → (0,2) → (0,3) → (1,3) → (2,3) → (3,3) → (3,2) → (3,1) → (3,0) → (2,0) → (1,0) → (0,0)
//
// 如果 dir == Clockwise，则每个格子向下一个位置（顺时针方向）移动；
// 如果 dir == CounterClockwise，则向上一个位置（逆时针方向）移动。
func RotateOuter(b *Board, dir Direction) {
	// 按顺时针顺序列出外圈坐标
	coords := [][2]int{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 3}, {2, 3},
		{3, 3}, {3, 2}, {3, 1}, {3, 0},
		{2, 0}, {1, 0},
	}

	// 将当前外圈所有格子的值依次存入 vals
	var vals []player.Color
	for _, rc := range coords {
		r, c := rc[0], rc[1]
		vals = append(vals, b.Get(r, c))
	}

	n := len(vals) // n == 12
	if n == 0 {
		return
	}

	// 根据 dir 计算新位置时的偏移量
	// 顺时针：当前位置 i 的新值来源于旧位置 (i-1+n)%n
	// 逆时针：当前位置 i 的新值来源于旧位置 (i+1)%n
	for i, rc := range coords {
		var srcIdx int
		if dir == Clockwise {
			srcIdx = (i - 1 + n) % n
		} else {
			// CounterClockwise
			srcIdx = (i + 1) % n
		}
		r, c := rc[0], rc[1]
		b.Set(r, c, vals[srcIdx])
	}
}

// RotateInner 对 4×4 棋盘的内圈 4 个格子执行“旋转”一格操作。
// 内圈坐标（顺时针顺序）为：
//
//	(1,1) → (1,2) → (2,2) → (2,1) → (1,1)
//
// 如果 dir == Clockwise，则每个格子向下一个位置（顺时针方向）移动；
// 如果 dir == CounterClockwise，则向上一个位置（逆时针方向）移动。
func RotateInner(b *Board, dir Direction) {
	// 按顺时针顺序列出内圈坐标
	coords := [][2]int{
		{1, 1}, {1, 2},
		{2, 2}, {2, 1},
	}

	// 将当前内圈所有格子的值依次存入 vals
	var vals []player.Color
	for _, rc := range coords {
		r, c := rc[0], rc[1]
		vals = append(vals, b.Get(r, c))
	}

	n := len(vals) // n == 4
	if n == 0 {
		return
	}

	// 根据 dir 计算新位置时的偏移量
	for i, rc := range coords {
		var srcIdx int
		if dir == Clockwise {
			srcIdx = (i - 1 + n) % n
		} else {
			// CounterClockwise
			srcIdx = (i + 1) % n
		}
		r, c := rc[0], rc[1]
		b.Set(r, c, vals[srcIdx])
	}
}
