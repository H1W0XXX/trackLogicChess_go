package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"trackLogicChess/internal/game"
	"trackLogicChess/internal/player"
	ui "trackLogicChess/internal/ui/gui"
)

func main() {
	// 启动参数
	outerFlag := flag.Int("outer", 0, "外圈旋转方向（0=顺时针,1=逆时针）")
	innerFlag := flag.Int("inner", 0, "内圈旋转方向（0=顺时针,1=逆时针）")
	useAI := flag.Bool("ai", true, "是否启用 AI 对手（AI 执 White）")
	uiMode := flag.String("ui", "terminal", "terminal | gui")
	flag.Parse()

	// 参数校验
	if *outerFlag != 0 && *outerFlag != 1 {
		fmt.Println("outer 参数无效，只能是 0（顺时针）或 1（逆时针）。")
		return
	}
	if *innerFlag != 0 && *innerFlag != 1 {
		fmt.Println("inner 参数无效，只能是 0（顺时针）或 1（逆时针）。")
		return
	}

	// 转换为 Direction
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

	// 创建游戏状态
	gState := game.NewGame(dirOuter, dirInner)

	// 根据 ui 参数选择运行模式
	switch *uiMode {
	case "terminal":
		launchGUI(gState, *useAI)
	default:
		runTerminalLoop(gState, *useAI)
	}
}

// launchGUI 以 Ebiten 窗口模式启动游戏
func launchGUI(gs *game.GameState, ai bool) {
	app := ui.NewApp(gs, ai)
	ebiten.SetWindowTitle("Track Logic Chess")
	ebiten.SetWindowResizable(false)
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}

// runTerminalLoop 原生命令行模式
func runTerminalLoop(g *game.GameState, useAI bool) {
	fmt.Println("=== Track Logic Chess (4×4 旋转棋) ===")
	fmt.Printf("外圈旋转：%s，内圈旋转：%s。\n",
		directionString(g.DirOuter), directionString(g.DirInner))
	if useAI {
		fmt.Println("已启用 AI 对手 (AI 执 White)。")
	} else {
		fmt.Println("人人对战模式。")
	}
	fmt.Println("人类玩家请输入：row col （0–3）")
	fmt.Println()
	fmt.Println("当前棋盘：")
	fmt.Println(g.Board.String())
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for !g.IsGameOver() {
		current := g.CurrentPlayer
		// AI 回合
		if useAI && current == player.White {
			fmt.Println("AI 正在思考...")
			mv := game.FindBestMoveDeep(g, 6)
			_ = g.ApplyMove(mv.Row, mv.Col)
			fmt.Printf("AI 在 (%d,%d) 下棋。\n", mv.Row, mv.Col)
		} else {
			// 人类回合
			fmt.Printf("轮到玩家 %s，请输入 (row col)：", current.String())
			if !scanner.Scan() {
				fmt.Println("\n读取输入失败，程序退出。")
				return
			}
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Println("输入格式错误，请输入 2 个数字，例如：1 2")
				continue
			}
			r, err1 := strconv.Atoi(parts[0])
			c, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || r < 0 || r > 3 || c < 0 || c > 3 {
				fmt.Println("坐标必须在 0–3 之间，请重试。")
				continue
			}
			if err := g.ApplyMove(r, c); err != nil {
				fmt.Println("操作无效：", err)
				continue
			}
		}

		// 显示最新棋盘
		fmt.Println("\n落子 + 旋转 后的棋盘：")
		fmt.Println(g.Board.String())
		fmt.Println()
	}

	// 结束判定
	if winner := g.WinnerColor(); winner == player.Empty {
		fmt.Println("棋盘已满，平局结束。")
	} else {
		fmt.Printf("游戏结束！玩家 %s 获胜。\n", winner.String())
	}
}

// directionString 将 Direction 转为中文
func directionString(d game.Direction) string {
	if d == game.Clockwise {
		return "顺时针"
	}
	return "逆时针"
}
