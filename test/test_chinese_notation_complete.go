package main

import (
	"fmt"
	"chinese-chess-ai/internal/chess"
)

func main() {
	fmt.Println("=== 中文棋谱解析测试 ===\n")
	
	// 测试1: 红方炮二平五
	fmt.Println("【测试1】红方炮二平五")
	board1 := chess.NewBoard()
	board1.Turn = chess.Red
	err := board1.MoveByChineseNotation("炮二平五", chess.Red)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: %s\n", board1.MoveList[len(board1.MoveList)-1])
	}
	fmt.Println()
	
	// 测试2: 黑方炮2平5
	fmt.Println("【测试2】黑方炮2平5")
	board2 := chess.NewBoard()
	board2.Turn = chess.Black
	err = board2.MoveByChineseNotation("炮2平5", chess.Black)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: %s\n", board2.MoveList[len(board2.MoveList)-1])
	}
	fmt.Println()
	
	// 测试3: 黑方马8进7
	fmt.Println("【测试3】黑方马8进7")
	board3 := chess.NewBoard()
	board3.Turn = chess.Black
	err = board3.MoveByChineseNotation("马8进7", chess.Black)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: %s\n", board3.MoveList[len(board3.MoveList)-1])
	}
	fmt.Println()
	
	// 测试4: 红方马二进三
	fmt.Println("【测试4】红方马二进三")
	board4 := chess.NewBoard()
	board4.Turn = chess.Red
	err = board4.MoveByChineseNotation("马二进三", chess.Red)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功: %s\n", board4.MoveList[len(board4.MoveList)-1])
	}
	fmt.Println()
	
	// 测试5: 完整对局
	fmt.Println("【测试5】完整对局模拟")
	board5 := chess.NewBoard()
	
	moves := []struct {
		notation string
		color    chess.Color
		desc     string
	}{
		{"炮二平五", chess.Red, "红方炮二平五"},
		{"马8进7", chess.Black, "黑方马8进7"},
		{"马二进三", chess.Red, "红方马二进三"},
		{"车9平8", chess.Black, "黑方车9平8"},
	}
	
	for i, move := range moves {
		board5.Turn = move.color
		err := board5.MoveByChineseNotation(move.notation, move.color)
		if err != nil {
			fmt.Printf("第%d步 %s ❌ 失败: %v\n", i+1, move.desc, err)
			break
		} else {
			fmt.Printf("第%d步 %s ✅ 成功\n", i+1, move.desc)
		}
	}
	
	fmt.Println("\n棋谱记录:")
	for i, move := range board5.MoveList {
		fmt.Printf("%d. %s\n", i+1, move)
	}
	
	fmt.Println("\n=== 所有测试完成 ===")
}
