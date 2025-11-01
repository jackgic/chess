package main

import (
	"fmt"
	"strings"
	"chinese-chess-ai/internal/chess"
)

func main() {
	fmt.Println("=== 测试中文棋谱修复 ===\n")
	
	board := chess.NewBoard()
	
	// 测试1: 红方炮二平五
	fmt.Println("测试1: 红方炮二平五")
	fmt.Println("预期: 炮从'二'路移动到'五'路（col 7 → col 4）")
	
	// 获取初始炮的位置
	fmt.Println("初始红方炮位置: (7, 1) 和 (7, 7)")
	fmt.Println("col 1 是'八'路，col 7 是'二'路")
	
	// 执行走子
	err := board.Move(77, 74) // 从(7,7)移动到(7,4)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	// 检查生成的棋谱
	if len(board.MoveList) > 0 {
		lastMove := board.MoveList[len(board.MoveList)-1]
		fmt.Printf("生成的棋谱: %s\n", lastMove)
		if lastMove == "炮二平五" {
			fmt.Println("✅ 正确！")
		} else {
			fmt.Printf("❌ 错误！应该是'炮二平五'，实际是'%s'\n", lastMove)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
	
	// 测试2: 黑方炮2平5
	fmt.Println("测试2: 黑方炮2平5")
	board.Turn = chess.Black
	
	err = board.Move(21, 24) // 从(2,1)移动到(2,4)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	if len(board.MoveList) > 1 {
		lastMove := board.MoveList[len(board.MoveList)-1]
		fmt.Printf("生成的棋谱: %s\n", lastMove)
		if lastMove == "炮2平5" {
			fmt.Println("✅ 正确！")
		} else {
			fmt.Printf("❌ 错误！应该是'炮2平5'，实际是'%s'\n", lastMove)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
	
	// 测试3: 红方马二进三
	fmt.Println("测试3: 红方马二进三")
	fmt.Println("预期: 马从'二'路进到'三'路（col 7 → col 6）")
	board.Turn = chess.Red
	
	err = board.Move(97, 76) // 从(9,7)移动到(7,6)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	if len(board.MoveList) > 2 {
		lastMove := board.MoveList[len(board.MoveList)-1]
		fmt.Printf("生成的棋谱: %s\n", lastMove)
		if lastMove == "马二进三" {
			fmt.Println("✅ 正确！")
		} else {
			fmt.Printf("❌ 错误！应该是'马二进三'，实际是'%s'\n", lastMove)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
	
	// 测试4: 列号对照表
	fmt.Println("测试4: 列号对照表")
	fmt.Println("\n红方视角（从右到左）:")
	fmt.Println("col 0 → 九")
	fmt.Println("col 1 → 八")
	fmt.Println("col 2 → 七")
	fmt.Println("col 3 → 六")
	fmt.Println("col 4 → 五")
	fmt.Println("col 5 → 四")
	fmt.Println("col 6 → 三")
	fmt.Println("col 7 → 二")
	fmt.Println("col 8 → 一")
	
	fmt.Println("\n黑方视角（从左到右）:")
	fmt.Println("col 0 → 1")
	fmt.Println("col 1 → 2")
	fmt.Println("col 2 → 3")
	fmt.Println("col 3 → 4")
	fmt.Println("col 4 → 5")
	fmt.Println("col 5 → 6")
	fmt.Println("col 6 → 7")
	fmt.Println("col 7 → 8")
	fmt.Println("col 8 → 9")
	
	fmt.Println("\n=== 测试完成 ===")
}
