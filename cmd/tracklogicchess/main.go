// 文件：cmd/tracklogicchess/main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"trackLogicChess/game"
	"trackLogicChess/player"
)

func main() {
	// 1. 启动时通过 flag 固定“外圈/内圈”方向
	outerFlag := flag.Int("outer", 0, "外圈旋转方向（0=顺时针, 1=逆时针）")
	innerFlag := flag.Int("inner", 0, "内圈旋转方向（0=顺时针, 1=逆时针）")
	useAI := flag.Bool("ai", false, "是否启用 AI 对手（AI 执 White）")
	flag.Parse()

	if *outerFlag != 0 && *outerFlag != 1 {
		fmt.Println("outer 参数无效，只能是 0（顺时针）或 1（逆时针）。")
		return
	}
	if *innerFlag != 0 && *innerFlag != 1 {
		fmt.Println("inner 参数无效，只能是 0（顺时针）或 1（逆时针）。")
		return
	}

	// 2. 将启动参数转换为 Direction
	var dirOuter, dirInner game.Direction
	if *outerFlag == 0 {
		dirOuter = game.Clockwise
	} else {
		dirOuter = game.CounterClockwise
	}
	if *innerFlag == 0 {
		dirInner = game.Clockwise
	} else {
		dirInner = game.CounterClockwise
	}

	// 3. NewGame 时传入这两个固定方向
	g := game.NewGame(dirOuter, dirInner)

	// 4. 打印欢迎与说明
	fmt.Println("=== Track Logic Chess (4×4 旋转棋) ===")
	fmt.Println("Black (○) 先手，White (●) 后手。")
	fmt.Printf("外圈旋转：%s，内圈旋转：%s。\n",
		directionString(dirOuter), directionString(dirInner))
	if *useAI {
		fmt.Println("已启用 AI 对手 (AI 执 White)。")
	} else {
		fmt.Println("人人对战模式。")
	}
	fmt.Println("人类玩家请输入：row col （0–3）")
	fmt.Println()

	// 5. 打印初始空棋盘
	fmt.Println("当前棋盘：")
	fmt.Println(g.Board.String())
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if g.IsGameOver() {
			break
		}

		current := g.CurrentPlayer

		// 如果启用 AI 并且轮到 White，就让 AI 走
		if *useAI && current == player.White {
			fmt.Println("AI 正在思考...")
			mv := game.FindBestMoveDeep(g, 6)
			// AI 只传落子位置；旋转方向已包含在 g.DirOuter/g.DirInner 中
			_ = g.ApplyMove(mv.Row, mv.Col)
			fmt.Printf("AI 在 (%d,%d) 下棋。\n", mv.Row, mv.Col)

		} else {
			// 人类回合：只读 row col
			fmt.Printf("轮到玩家 %s，请输入 (row col)：", current.String())
			if !scanner.Scan() {
				fmt.Println("\n读取输入失败，程序退出。")
				return
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Println("输入格式错误，请输入 2 个数字，例如：1 2")
				continue
			}
			r, err1 := strconv.Atoi(parts[0])
			c, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				fmt.Println("输入必须是 0–3 之间的整数，请重试。")
				continue
			}
			if r < 0 || r > 3 || c < 0 || c > 3 {
				fmt.Println("坐标越界，row, col 必须在 0–3 之间。")
				continue
			}
			// 只传落子位置，旋转直接用 g.DirOuter/g.DirInner
			err := g.ApplyMove(r, c)
			if err != nil {
				fmt.Println("操作无效：", err.Error())
				continue
			}
		}

		// 打印本回合结束后的棋盘
		fmt.Println("\n落子 + 旋转 后的棋盘：")
		fmt.Println(g.Board.String())
		fmt.Println()
	}

	// 6. 游戏结束，显示结果
	if winner := g.WinnerColor(); winner == player.Empty {
		fmt.Println("棋盘已满，平局结束。")
	} else {
		fmt.Printf("游戏结束！玩家 %s 获胜。\n", g.WinnerColor().String())
	}
}

// directionString 将 Direction 转为可读字符串
func directionString(d game.Direction) string {
	if d == game.Clockwise {
		return "顺时针"
	}
	return "逆时针"
}
