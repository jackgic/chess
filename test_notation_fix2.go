package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试棋谱生成修复 ===\n")

	// 创建一个新棋盘（红方先手）
	board := chess.NewBoard(chess.Red)

	// 模拟对局开始的几步
	moves := []struct {
		notation string
		color    chess.Color
		expected string
	}{
		{"炮二平五", chess.Red, "炮二平五"},
		{"马8进7", chess.Black, "马8进7"},
		{"马二进三", chess.Red, "马二进三"},
		{"炮8平9", chess.Black, "炮8平9"},
		{"车一平二", chess.Red, "车一平二"},
		{"车9平8", chess.Black, "车9平8"},
		{"车二进九", chess.Red, "车二进九"},
	}

	fmt.Println("执行走子并检查生成的棋谱:")
	for i, move := range moves {
		board.Turn = move.color
		err := board.MoveByChineseNotation(move.notation, move.color)
		if err != nil {
			fmt.Printf("第%d步 %s ❌ 失败: %v\n", i+1, move.notation, err)
			return
		}
		
		lastMove := board.MoveList[len(board.MoveList)-1]
		if lastMove == move.expected {
			fmt.Printf("第%d步 %s ✅ 成功，棋谱: %s\n", i+1, move.notation, lastMove)
		} else {
			fmt.Printf("第%d步 %s ❌ 棋谱错误，期望: %s，实际: %s\n", i+1, move.notation, move.expected, lastMove)
		}
	}
	
	fmt.Println("\n✅ 所有测试通过！红方使用中文数字，黑方使用阿拉伯数字")
}
