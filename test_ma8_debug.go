package main

import (
	"fmt"
	"chinese-chess-ai/internal/chess"
)

func main() {
	board := chess.NewBoard()
	
	fmt.Println("=== 初始棋盘布局 ===")
	fmt.Println("黑方马的位置:")
	fmt.Println("- col 1 (黑方记谱'2'): row 0")
	fmt.Println("- col 7 (黑方记谱'8'): row 0")
	fmt.Println()
	
	// 测试解析"马8进7"
	fmt.Println("=== 测试: 马8进7 ===")
	fmt.Println("分析:")
	fmt.Println("- '马8' 表示黑方在第8列的马")
	fmt.Println("- 黑方第8列 = col 7")
	fmt.Println("- '进7' 对于马来说，表示移动到第7列")
	fmt.Println("- 黑方第7列 = col 6")
	fmt.Println("- 马应该从 (0,7) 移动到 (2,6)")
	fmt.Println()
	
	// 手动计算预期的移动
	fmt.Println("=== 手动计算 ===")
	fromRow := 0
	fromCol := 7
	toCol := 6  // 黑方"7" = col 6
	colDiff := 7 - 6  // = 1
	toRow := fromRow + 2  // 黑方向下，列差为1，行差为2
	fmt.Printf("预期移动: (%d,%d) -> (%d,%d)\n", fromRow, fromCol, toRow, toCol)
	fmt.Printf("列差: %d, 行差: %d\n", colDiff, toRow-fromRow)
	fmt.Printf("from = %d, to = %d\n", fromRow*10+fromCol, toRow*10+toCol)
	fmt.Println()
	
	// 先测试直接Move
	fmt.Println("=== 测试直接Move ===")
	board.Turn = chess.Black  // 确保是黑方回合
	err := board.Move(7, 26)  // from=7 (0,7), to=26 (2,6)
	if err != nil {
		fmt.Printf("❌ 直接Move失败: %v\n", err)
	} else {
		fmt.Println("✅ 直接Move成功")
	}
	fmt.Println()
	
	// 重新初始化棋盘
	board = chess.NewBoard()
	board.Turn = chess.Black
	
	fmt.Println("=== 测试MoveByChineseNotation ===")
	err = board.MoveByChineseNotation("马8进7", chess.Black)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		
		// 详细调试
		fmt.Println()
		fmt.Println("=== 详细调试 ===")
		
		// 检查col 7是否有黑马
		piece := board.Grid[0][7]
		fmt.Printf("Grid[0][7] = Type:%d, Color:%d\n", piece.Type, piece.Color)
		fmt.Printf("期望: Type:%d (Horse), Color:%d (Black)\n", chess.Horse, chess.Black)
		
		if piece.Type == chess.Horse && piece.Color == chess.Black {
			fmt.Println("✅ 找到了黑方的马在 col 7")
		} else {
			fmt.Println("❌ col 7 没有黑方的马")
		}
	} else {
		fmt.Println("✅ 成功执行走子")
		fmt.Printf("最后一步棋谱: %s\n", board.MoveList[len(board.MoveList)-1])
	}
}
