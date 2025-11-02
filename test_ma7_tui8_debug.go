package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试马7退8的详细调试 ===\n")

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

	// 打印当前棋盘状态
	fmt.Println("当前棋盘状态:")
	fmt.Println("黑方第7列(col=6)的棋子:")
	for row := 0; row < 10; row++ {
		piece := board.Grid[row][6]
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
	
	// 马在 (2, 6)，要退到 col 7
	// 马走日字：列差1，行差2；或列差2，行差1
	// 从 col 6 到 col 7，列差 = 1，所以行差应该是 2
	// 退：向下移动，所以 toRow = 2 + 2 = 4
	fmt.Println("\n预期计算:")
	fmt.Println("  - 起始位置: (2, 6)")
	fmt.Println("  - 目标列: 7")
	fmt.Println("  - 列差: |7 - 6| = 1")
	fmt.Println("  - 行差: 2 (马走日字)")
	fmt.Println("  - 退：向下移动，toRow = 2 + 2 = 4")
	fmt.Println("  - 目标位置: (4, 7)")
	
	// 检查目标位置
	fmt.Println("\n目标位置 (4, 7) 的棋子:")
	piece := board.Grid[4][7]
	if piece.Type == chess.Empty {
		fmt.Println("  空位 ✅")
	} else {
		pieceNames := map[chess.PieceType]string{
			chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
			chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
		}
		colorNames := map[chess.Color]string{
			chess.Red: "红方", chess.Black: "黑方",
		}
		fmt.Printf("  %s %s\n", colorNames[piece.Color], pieceNames[piece.Type])
	}
	
	// 检查马的合法移动
	fmt.Println("\n检查马从 (2, 6) 到 (4, 7) 是否合法:")
	from := chess.Position{Row: 2, Col: 6}
	to := chess.Position{Row: 4, Col: 7}
	if board.IsValidMove(from, to) {
		fmt.Println("  ✅ 合法移动")
	} else {
		fmt.Println("  ❌ 非法移动")
	}
	
	// 尝试执行
	fmt.Println("\n尝试执行 '马7退8':")
	board.Turn = chess.Black
	err := board.MoveByChineseNotation("马7退8", chess.Black)
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
	} else {
		fmt.Printf("  ✅ 成功\n")
	}
}
