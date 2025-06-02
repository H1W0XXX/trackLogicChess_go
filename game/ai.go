// File game/ai.go
package game

import (
	"math"
	"math/rand"
	"time"

	"trackLogicChess/player"
)

/* ---------- 基础结构 ---------- */

// Move 仅记录落子坐标，旋转方向固定由 GameState.DirOuter / DirInner 决定。
type Move struct {
	Row int
	Col int
}

// opposite 返回相反颜色。
func opposite(c player.Color) player.Color {
	if c == player.Black {
		return player.White
	}
	return player.Black
}

/* ---------- 启发式评估 ---------- */

// lineScores[i]：一条 4 格直线里，己方有 i 子、对手 0 子 的加分（反之减分）。
var lineScores = [5]int{0, 1, 8, 64, 1_000_000} // 4 连直接视为绝杀

// heuristicScore 对整盘局面做线型统计。
// - 10 条直线：4 行 + 4 列 + 2 对角。
// - 若同一条线上双方都有子，计 0 分；否则按己/敌子数累加或累减。
func heuristicScore(b *Board, me player.Color) int {
	op := opposite(me)
	score := 0

	lines := [][4][2]int{
		// 行
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{1, 0}, {1, 1}, {1, 2}, {1, 3}},
		{{2, 0}, {2, 1}, {2, 2}, {2, 3}},
		{{3, 0}, {3, 1}, {3, 2}, {3, 3}},
		// 列
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
		{{0, 2}, {1, 2}, {2, 2}, {3, 2}},
		{{0, 3}, {1, 3}, {2, 3}, {3, 3}},
		// 对角
		{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		{{0, 3}, {1, 2}, {2, 1}, {3, 0}},
	}

	for _, ln := range lines {
		myCnt, opCnt := 0, 0
		for _, rc := range ln {
			switch b.Get(rc[0], rc[1]) {
			case me:
				myCnt++
			case op:
				opCnt++
			}
		}
		if myCnt > 0 && opCnt == 0 {
			score += lineScores[myCnt]
		} else if opCnt > 0 && myCnt == 0 {
			score -= lineScores[opCnt]
		}
	}
	return score
}

/* ---------- Negamax + α-β 剪枝 ---------- */

const (
	defaultDepth = 5         // 默认搜索深度
	winScore     = 1_000_000 // 必胜分
	loseScore    = -winScore
)

// FindBestMoveDeep 使用 Negamax 搜索给出最佳着法；depth ≤0 时采用 defaultDepth。
func FindBestMoveDeep(g *GameState, depth int) Move {
	if depth <= 0 {
		depth = defaultDepth
	}

	moves := g.GenerateMoves()
	if len(moves) == 0 {
		return Move{-1, -1}
	}

	// 随机打乱，避免评分相同总走同一手
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })

	bestScore := math.MinInt
	bestMove := moves[0]

	for _, mv := range moves {
		sim := g.cloneGameState()
		sim.Board.Set(mv.Row, mv.Col, sim.CurrentPlayer)
		if CheckWin(sim.Board, sim.CurrentPlayer) { // 一步必杀
			return mv
		}
		RotateOuter(sim.Board, sim.DirOuter)
		RotateInner(sim.Board, sim.DirInner)
		sim.CurrentPlayer = opposite(sim.CurrentPlayer)

		score := -negamax(sim, depth-1, loseScore, winScore)

		if score > bestScore {
			bestScore, bestMove = score, mv
		}
	}
	return bestMove
}

// negamax 递归：当前 gs.CurrentPlayer 视角，返回局面评分。
func negamax(gs *GameState, depth, alpha, beta int) int {
	// 终局：上一手把对方连 4，则我方输
	if CheckWin(gs.Board, opposite(gs.CurrentPlayer)) {
		return loseScore + (defaultDepth - depth) // 越晚输分数越高（延迟被杀）
	}
	// 深度到 0 或平局
	if depth == 0 || gs.isBoardFull() {
		return heuristicScore(gs.Board, gs.CurrentPlayer)
	}

	for _, mv := range gs.GenerateMoves() {
		sim := gs.cloneGameState()
		sim.Board.Set(mv.Row, mv.Col, sim.CurrentPlayer)
		if CheckWin(sim.Board, sim.CurrentPlayer) {
			return winScore - (defaultDepth - depth) // 越早杀分越高
		}
		RotateOuter(sim.Board, sim.DirOuter)
		RotateInner(sim.Board, sim.DirInner)
		sim.CurrentPlayer = opposite(sim.CurrentPlayer)

		score := -negamax(sim, depth-1, -beta, -alpha)
		if score > alpha {
			alpha = score
			if alpha >= beta { // β 剪枝
				break
			}
		}
	}
	return alpha
}

/* ---------- 原有辅助 ---------- */

// GenerateMoves：列出所有空格
func (g *GameState) GenerateMoves() []Move {
	if g.GameOver {
		return nil
	}
	var mv []Move
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if g.Board.IsEmpty(r, c) {
				mv = append(mv, Move{r, c})
			}
		}
	}
	return mv
}

// cloneBoard 深拷贝棋盘
func (b *Board) cloneBoard() *Board {
	var nb Board
	for r := range b.cells {
		copy(nb.cells[r][:], b.cells[r][:])
	}
	return &nb
}

// cloneGameState 深拷贝局面（含固定方向）
func (g *GameState) cloneGameState() *GameState {
	return &GameState{
		Board:         g.Board.cloneBoard(),
		CurrentPlayer: g.CurrentPlayer,
		DirOuter:      g.DirOuter,
		DirInner:      g.DirInner,
		Winner:        g.Winner,
		GameOver:      g.GameOver,
	}
}
