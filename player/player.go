// File player/player.go
package player

// Color 表示棋子颜色或空状态。
// Empty 表示该格子为空，Black 表示黑子，White 表示白子。
type Color int

const (
	Empty Color = iota
	Black
	White
)

// String 返回 Color 对应的可读字符串，方便调试与打印。
func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Empty"
	}
}
