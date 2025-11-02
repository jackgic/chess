package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试所有修复 ===\n")

	// 测试1: 红方步数使用中文数字
	fmt.Println("【测试1】红方步数使用中文数字")
	board1 := chess.NewBoard(chess.Red)
	board1.Turn = chess.Red
	board1.Move(90, 91) // 车一平二
	board1.Turn = chess.Red
	board1.Move(91, 1) // 车二进九
	
	if len(board1.MoveList) >= 2 {
		move1 := board1.MoveList[0]
		move2 := board1.MoveList[1]
		if move1 == "车一平二" && move2 == "车二进九" {
			fmt.Printf("✅ 成功: %s, %s\n", move1, move2)
		} else {
			fmt.Printf("❌ 失败: %s, %s (期望: 车一平二, 车二进九)\n", move1, move2)
		}
	}

	// 测试2: 黑方步数使用阿拉伯数字
	fmt.Println("\n【测试2】黑方步数使用阿拉伯数字")
	board2 := chess.NewBoard(chess.Red)
	board2.Turn = chess.Black
	board2.Move(8, 28) // 车9进2
	
	if len(board2.MoveList) >= 1 {
		move := board2.MoveList[0]
		if move == "车9进2" {
			fmt.Printf("✅ 成功: %s\n", move)
		} else {
			fmt.Printf("❌ 失败: %s (期望: 车9进2)\n", move)
		}
	}

	// 测试3: 马7退8（黑方马向上退）
	fmt.Println("\n【测试3】马7退8（黑方马向上退）")
	board3 := chess.NewBoard(chess.Red)
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
		{"马7退8", chess.Black}, // 吃红方车
	}

	success := true
	for i, move := range moves {
		board3.Turn = move.color
		err := board3.MoveByChineseNotation(move.notation, move.color)
		if err != nil {
			fmt.Printf("❌ 第%d步 %s 失败: %v\n", i+1, move.notation, err)
			success = false
			break
		}
	}
	if success {
		fmt.Println("✅ 成功: 马7退8正确执行（从(2,6)到(0,7)吃红方车）")
	}

	// 测试4: 马7进8（会被蹩马腿，应该失败）
	fmt.Println("\n【测试4】马7进8（会被蹩马腿，应该失败）")
	board4 := chess.NewBoard(chess.Red)
	moves2 := []struct {
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

	for _, move := range moves2 {
		board4.Turn = move.color
		board4.MoveByChineseNotation(move.notation, move.color)
	}

	board4.Turn = chess.Black
	err := board4.MoveByChineseNotation("马7进8", chess.Black)
	if err != nil {
		fmt.Printf("✅ 成功: 马7进8被正确拒绝（蹩马腿）: %v\n", err)
	} else {
		fmt.Println("❌ 失败: 马7进8应该被拒绝（蹩马腿）")
	}

	fmt.Println("\n=== 所有测试完成 ===")
}
