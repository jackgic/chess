package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试马7退8的问题 ===\n")

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
		{"车二进九", chess.Red, "红方车二进九"},
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

	// 打印当前棋盘状态
	fmt.Println("\n当前棋盘状态:")
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

	fmt.Println("\n黑方第8列(col=7)的棋子:")
	for row := 0; row < 10; row++ {
		piece := board.Grid[row][7]
		if piece.Type != chess.Empty && piece.Color == chess.Black {
			pieceNames := map[chess.PieceType]string{
				chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
				chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
			}
			fmt.Printf("  Row %d, Col 7: %s (黑方)\n", row, pieceNames[piece.Type])
		}
	}

	// 现在测试马7退8
	fmt.Println("\n【关键测试】黑方马7退8")
	board.Turn = chess.Black
	err := board.MoveByChineseNotation("马7退8", chess.Black)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功\n")
	}
}
