# Track Logic Chess

**4×4 旋转棋（控制台版）**

这个项目实现了一个简化版的“Track Logic Chess”：

* **棋盘**：4×4，共 16 格
* **玩家**：黑（Black，●）先手，白（White，○）后手
* **胜利条件**：落子并执行“内圈/外圈”固定方向旋转后，若整行/整列/对角线出现连续 4 子，即判胜
* **旋转规则**：

  * **外圈**（12 格）环移一格
  * **内圈**（4 格）顺时针/逆时针旋转一格
  * 启动时通过命令行参数一次性指定，不可在对局中修改

此外，代码内置了一个基于 Negamax+α-β 的简易 AI，可通过 `-ai` 参数让后手由 AI 操作。

---

## 功能特性

* 控制台实时输出棋盘（`. ● ○` 表示空/黑/白）
* 启动时指定外圈与内圈旋转方向（顺时针或逆时针）
* 支持“人人对战”与“人机对战”两种模式
* 内置 Negamax+α-β 剪枝搜索，AI 支持可调深度（默认深度 5）
* “双四连”判定可视为平局、先手胜或后手胜（可按需要修改判断逻辑）

---

## 先决条件

* **Go 版本**：1.18 及以上
* **操作系统**：Windows / macOS / Linux
* **环境**：已正确配置 `GOPATH` 或使用 Go Modules

---

## 快速开始

1. **克隆仓库**

   ```bash
   git clone <仓库地址>
   cd trackLogicChess
   ```

2. **编译可执行文件**

   ```bash
   go build -o tracklogicchess.exe ./cmd/tracklogicchess
   ```

   生成的可执行文件会在当前目录，名为 `tracklogicchess.exe`（Windows）或 `tracklogicchess`（macOS/Linux）。

3. **运行游戏（人人对战）**

   ```bash
   ./tracklogicchess -outer 0 -inner 0
   ```

   * `-outer 0` 表示外圈顺时针；`-outer 1` 表示外圈逆时针
   * `-inner 0` 表示内圈顺时针；`-inner 1` 表示内圈逆时针
   * 默认不加 `-ai`，双方玩家都手动输入落子坐标

4. **运行游戏（人机对战）**

   ```bash
   ./tracklogicchess -outer 0 -inner 0 -ai
   ```

   * 加上 `-ai` 后，**White**（后手）由 AI 操作，Black（先手）由你通过控制台输入
   * AI 默认使用深度 5 的 Negamax 搜索（大约 1–5 ms / 步）

5. **指定 AI 搜索深度**
   在 `main.go` 中将 AI 调用改为：

   ```go
   mv := game.FindBestMoveDeep(g, 7)
   ```

   或者传递深度参数（如果你自行修改了 `FindBestMoveDeep` 接口）。

   * 深度越大，AI 越强，但耗时增加
   * 建议深度 5\~7 已能较好对抗人类

---

## 控制说明

* **落子输入**：对局中只需输入两个数字——`row col`，行列都在 0 \~ 3 之间

  * 举例：

    ```
    轮到玩家 ●，请输入 (row col)： 1 2
    ```

    表示在第 2 行、第 3 列落子（从 0 开始计数）
* **旋转方向**：由启动参数一次性指定，对局中不可更改
* **对局结束**：满足连 4 或棋盘填满后自动结束并打印结果

  * 双方同时出现连 4 时判定为平局（可按需要修改为先手/后手优先）

---

## 示例对局

```text
> ./tracklogicchess -outer 0 -inner 0 -ai

=== Track Logic Chess (4×4 旋转棋) ===
Black (●) 先手，White (○) 后手。
外圈旋转：顺时针，内圈旋转：顺时针。
已启用 AI 对手 (AI 执 White)。
人类玩家请输入：row col（0–3）

当前棋盘：
. . . .
. . . .
. . . .
. . . .

轮到玩家 ●，请输入 (row col)： 1 1

--- 旋转前棋盘 ---
. . . .
. ● . .
. . . .
. . . .

落子 + 旋转 后的棋盘：
● . . .
. . . .
. . . .
. . . .

轮到玩家 ○，AI 正在思考...
AI 在 (2, 2) 下棋。

--- 旋转前棋盘 ---
● . . .
. . . .
. . ○ .
. . . .

落子 + 旋转 后的棋盘：
. ● . .
. . . .
. . . .
. . . .

轮到玩家 ●，请输入 (row col)： 0 1
...
```

（以下省略）

---

## 项目结构

```
trackLogicChess/
├── cmd/
│   └── tracklogicchess/
│       ├── main.go          # 程序入口与运行逻辑
│
├── game/                    # 核心逻辑包
│   ├── board.go             # 棋盘数据结构与打印
│   ├── ring.go              # 内/外圈旋转函数
│   ├── rules.go             # 胜负判定 & AI 搜索（Negamax + α-β）
│   └── game.go              # GameState、落子 & 回合流程
│
├── player/                  
│   └── player.go            # 定义棋子颜色枚举和值
│
├── ui/                      # （保留，可后续扩展 GUI）
│   ├── terminal/
│   │   └── terminal.go      # 终端版交互（目前主要通过 main.go 实现）
│   └── gui/
│       └── gui.go           # GUI 版占位，暂未实现
│
├── go.mod                   # Go Modules 配置
└── README.md                # 项目说明（本文件）
```

* **`cmd/tracklogicchess/main.go`**：解析命令行参数，创建 `GameState`，循环读取落子/AI 决策并更新棋盘
* **`game/board.go`**：`Board` 结构体及 `.String()` 用于控制台打印
* **`game/ring.go`**：实现“外圈环移”与“内圈旋转”两种固定方向函数
* **`game/rules.go`**：包含 `CheckWin` 胜负判定，以及 `FindBestMoveDeep`（Negamax + α-β）AI 实现
* **`game/game.go`**：`GameState` 维护当前棋局状态、落子 & 旋转顺序、结束判定
* **`player/player.go`**：`player.Color` 定义（`Empty`、`Black`、`White`）
* **`ui/`**：预留终端/GUI 包，用于后续界面扩展

---

## 设计思路

1. **数据结构**

   * `Board` 使用 `[4][4]player.Color` 存储格子状态
   * `GameState` 保存当前玩家、棋盘、固定旋转方向、胜者 & 结束标志

2. **落子与旋转流程**

   ```
   ApplyMove(r, c):
     1. 检查合法性（游戏是否已结束 / 目标格子是否空）
     2. 在 (r,c) 放置棋子
     3. 执行 g.DirOuter 的外圈环移 & g.DirInner 的内圈旋转
     4. 旋转后判断胜负：
        - 如果当前玩家连 4，GameOver = true；Winner = 当前玩家
        - 否则若棋盘已满，GameOver = true；Winner = Empty（平局）
        - 否则切换玩家
   ```

3. **AI 实现**

   * **Negamax + α-β 剪枝**，默认深度 5（可调整）
   * **启发式评分（`heuristicScore`）**：对 10 条直线（4 行 + 4 列 + 2 对角）进行评分，

     * 若某条线上只有己方有 N 子则加 `lineScores[N]`
     * 若只有对手有 N 子则减 `lineScores[N]`
     * 阻塞（双方都有子）则该条不计分
   * 在搜索中：

     * 若落子后直接 `CheckWin` ，立即返回极大值或极小值
     * 否则递归深入，配合 α/β 界限剪枝
     * 终止时返回启发式分数

---
