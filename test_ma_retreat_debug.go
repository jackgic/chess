package main

import (
	"fmt"
	"chinese-chess-ai/internal/chess"
)

func main() {
	fmt.Println("=== 测试马7退8的问题（详细调试）===\n")
	
	// 创建一个新棋盘（红方先手）
	board := chess.NewBoard(chess.Red)
	
	// 模拟对局开始的几步
	moves := []struct {
		notation string
		color    chess.Color
		desc     string
	}{
		{"炮二平五", chess.Red, "红方炮二平五"},
		{"马8进7", chess.Black, "黑方马8进7"},
		{"马二进三", chess.Red, "红方马二进三"},
		{"炮8平9", chess.Black, "黑方炮8平9"},
		{"车一平二", chess.Red, "红方车一平二"},
		{"车9平8", chess.Black, "黑方车9平8"},
		{"车二进9", chess.Red, "红方车二进9"},
	}
	
	fmt.Println("执行前面的走子:")
	for i, move := range moves {
		board.Turn = move.color
		err := board.MoveByChineseNotation(move.notation, move.color)
		if err != nil {
			fmt.Printf("第%d步 %s ❌ 失败: %v\n", i+1, move.desc, err)
			return
		}
		fmt.Printf("第%d步 %s ✅ 成功\n", i+1, move.desc)
	}
	
	// 现在测试马7退8
	fmt.Println("\n【关键测试】黑方马7退8")
	board.Turn = chess.Black
	
	// 检查黑方第7列的所有棋子
	fmt.Println("\n黑方第7列的棋子:")
	for row := 0; row < 10; row++ {
		piece := board.Grid[row][6] // 黑方7列对应col=6
		if piece.Type != chess.Empty && piece.Color == chess.Black {
			pieceNames := map[chess.PieceType]string{
				chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
				chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
			}
			fmt.Printf("  Row %d, Col 6: %s (黑方)\n", row, pieceNames[piece.Type])
		}
	}
	
	// 分析"马7退8"的含义
	fmt.Println("\n分析'马7退8':")
	fmt.Println("  - 起始列: 7 (黑方) → col = 6")
	fmt.Println("  - 动作: 退 (向己方方向移动)")
	fmt.Println("  - 目标列: 8 (黑方) → col = 7")
	fmt.Println("  - 黑方在上方(row 0-4)，'退'表示向下移动(row增大)")
	
	// 如果第7列有多个马，需要指定是哪个
	fmt.Println("\n问题分析:")
	fmt.Println("  第7列有多个黑方棋子，'马7退8'无法确定是哪个马")
	fmt.Println("  标准棋谱应该使用'前马'或'后马'来区分")
	
	// 尝试执行
	err := board.MoveByChineseNotation("马7退8", chess.Black)
	if err != nil {
		fmt.Printf("\n执行结果: ❌ 失败: %v\n", err)
	} else {
		fmt.Printf("\n执行结果: ✅ 成功\n")
	}
	
	// 测试可能的正确走法
	fmt.Println("\n\n=== 测试可能的正确走法 ===")
	
	// 重新创建棋盘
	board2 := chess.NewBoard(chess.Red)
	for i, move := range moves {
		board2.Turn = move.color
		err := board2.MoveByChineseNotation(move.notation, move.color)
		if err != nil {
			fmt.Printf("第%d步失败: %v\n", i+1, err)
			return
		}
	}
	
	// 测试：前马退8（row 2的马）
	fmt.Println("\n测试1: 前马7退8")
	board2.Turn = chess.Black
	// 手动构造移动：从(2,6)到(4,7)
	// 马走日字：列差1，行差2
	from := 2*10 + 6  // row=2, col=6
	to := 4*10 + 7    // row=4, col=7
	err = board2.Move(from, to)
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
	} else {
		fmt.Printf("  ✅ 成功: 马从(2,6)移动到(4,7)\n")
	}
}
