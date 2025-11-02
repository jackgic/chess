package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试棋谱生成修复 ===\n")

	// 创建一个新棋盘（红方先手）
	board := chess.NewBoard(chess.Red)

	// 测试红方车二进九
	fmt.Println("测试1: 红方车一平二")
	board.Turn = chess.Red
	err := board.Move(90, 91) // 从(9,0)到(9,1)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		lastMove := board.MoveList[len(board.MoveList)-1]
		fmt.Printf("✅ 成功，生成的棋谱: %s\n", lastMove)
		if lastMove == "车一平二" {
			fmt.Println("✅ 棋谱格式正确")
		} else {
			fmt.Printf("❌ 棋谱格式错误，期望: 车一平二，实际: %s\n", lastMove)
		}
	}

	fmt.Println("\n测试2: 红方车二进九")
	board.Turn = chess.Red
	err = board.Move(91, 1) // 从(9,1)到(0,1)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		lastMove := board.MoveList[len(board.MoveList)-1]
		fmt.Printf("✅ 成功，生成的棋谱: %s\n", lastMove)
		if lastMove == "车二进九" {
			fmt.Println("✅ 棋谱格式正确（使用中文数字）")
		} else {
			fmt.Printf("❌ 棋谱格式错误，期望: 车二进九，实际: %s\n", lastMove)
		}
	}

	fmt.Println("\n测试3: 黑方车9进1")
	board2 := chess.NewBoard(chess.Red)
	board2.Turn = chess.Black
	err = board2.Move(8, 28) // 从(0,8)到(2,8)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		lastMove := board2.MoveList[len(board2.MoveList)-1]
		fmt.Printf("✅ 成功，生成的棋谱: %s\n", lastMove)
		if lastMove == "车9进2" {
			fmt.Println("✅ 棋谱格式正确（使用阿拉伯数字）")
		} else {
			fmt.Printf("❌ 棋谱格式错误，期望: 车9进2，实际: %s\n", lastMove)
		}
	}
}
