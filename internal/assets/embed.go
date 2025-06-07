package assets

import "embed"

// Images 嵌入全部 PNG；路径必须相对 embed.go 所在目录
//
//go:embed images/*.png
var Images embed.FS
