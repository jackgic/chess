package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 重新分析马7退8 ===\n")

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
	fmt.Println("\n黑方第7列(col=6)的棋子:")
	for row := 0; row < 10; row++ {
		piece := board.Grid[row][6]
		if piece.Type != chess.Empty {
			pieceNames := map[chess.PieceType]string{
				chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
				chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
			}
			colorNames := map[chess.Color]string{
				chess.Red: "红方", chess.Black: "黑方",
			}
			fmt.Printf("  Row %d, Col 6: %s %s\n", row, colorNames[piece.Color], pieceNames[piece.Type])
		}
	}

	// 分析"马7退8"的含义
	fmt.Println("\n分析'马7退8':")
	fmt.Println("  - 起始列: 7 (黑方) → col = 6")
	fmt.Println("  - 动作: 退 (向己方方向移动，黑方向上，row减小)")
	fmt.Println("  - 目标列: 8 (黑方) → col = 7")
	fmt.Println("  - 黑方在上方(row 0-4)，'退'表示向上移动(row减小)")
	
	// 马在 (2, 6)，要退到 col 7
	// 退：向上移动，所以 toRow = 2 - 2 = 0
	fmt.Println("\n正确的计算:")
	fmt.Println("  - 起始位置: (2, 6)")
	fmt.Println("  - 目标列: 7")
	fmt.Println("  - 列差: |7 - 6| = 1")
	fmt.Println("  - 行差: 2 (马走日字)")
	fmt.Println("  - 退：向上移动，toRow = 2 - 2 = 0")
	fmt.Println("  - 目标位置: (0, 7)")
	fmt.Println("  - 马腿位置: (1, 6)")
	
	// 检查马腿位置
	fmt.Println("\n检查马腿位置 (1, 6):")
	piece := board.Grid[1][6]
	if piece.Type == chess.Empty {
		fmt.Println("  空位 ✅ 不蹩马腿")
	} else {
		pieceNames := map[chess.PieceType]string{
			chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
			chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
		}
		colorNames := map[chess.Color]string{
			chess.Red: "红方", chess.Black: "黑方",
		}
		fmt.Printf("  %s %s ❌ 蹩马腿\n", colorNames[piece.Color], pieceNames[piece.Type])
	}
	
	// 检查目标位置
	fmt.Println("\n检查目标位置 (0, 7):")
	targetPiece := board.Grid[0][7]
	if targetPiece.Type == chess.Empty {
		fmt.Println("  空位 ✅")
	} else {
		pieceNames := map[chess.PieceType]string{
			chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
			chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
		}
		colorNames := map[chess.Color]string{
			chess.Red: "红方", chess.Black: "黑方",
		}
		fmt.Printf("  %s %s\n", colorNames[targetPiece.Color], pieceNames[targetPiece.Type])
	}
	
	// 分析"马7进8"（错误的走法）
	fmt.Println("\n\n对比：'马7进8'（会被蹩马腿）:")
	fmt.Println("  - 起始位置: (2, 6)")
	fmt.Println("  - 进：向下移动，toRow = 2 + 2 = 4")
	fmt.Println("  - 目标位置: (4, 7)")
	fmt.Println("  - 马腿位置: (3, 6)")
	
	fmt.Println("\n检查马腿位置 (3, 6):")
	piece2 := board.Grid[3][6]
	if piece2.Type == chess.Empty {
		fmt.Println("  空位 ✅ 不蹩马腿")
	} else {
		pieceNames := map[chess.PieceType]string{
			chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
			chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
		}
		colorNames := map[chess.Color]string{
			chess.Red: "红方", chess.Black: "黑方",
		}
		fmt.Printf("  %s %s ❌ 蹩马腿\n", colorNames[piece2.Color], pieceNames[piece2.Type])
	}
}
