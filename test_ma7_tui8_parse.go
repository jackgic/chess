package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试马7退8的解析 ===\n")

	// 创建一个新棋盘（红方先手）
	board := chess.NewBoard(chess.Red)

	// 模拟对局开始的几步
	moves := []struct {
		notation string
		color    chess.Color
	}{
		{"炮二平五", chess.Red},
		{"马8进7", chess.Black},
		{"马二进三", chess.Red},
		{"炮8平9", chess.Black},
		{"车一平二", chess.Red},
		{"车9平8", chess.Black},
		{"车二进九", chess.Red},
	}

	for _, move := range moves {
		board.Turn = move.color
		board.MoveByChineseNotation(move.notation, move.color)
	}

	fmt.Println("当前黑方马在位置: (2, 6)")
	fmt.Println("目标位置 (0, 7) 有红方车")
	fmt.Println()
	
	// 测试马7退8
	fmt.Println("测试: 马7退8")
	board.Turn = chess.Black
	err := board.MoveByChineseNotation("马7退8", chess.Black)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		fmt.Println("\n这是一个BUG！马7退8应该是合法的走法（吃红方车）")
	} else {
		fmt.Printf("✅ 成功\n")
		fmt.Println("马成功从 (2, 6) 移动到 (0, 7)，吃掉红方车")
	}
}
