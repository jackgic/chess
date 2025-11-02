package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试炮二平五 ===")

	// 创建新棋盘，红方先手
	board := chess.NewBoard(chess.Red)

	// 打印初始棋盘状态
	fmt.Println("初始棋盘:")
	printBoard(board)

	// 测试炮二平五
	fmt.Println("\n测试: 炮二平五")
	err := board.MoveByChineseNotation("炮二平五", chess.Red)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功\n")
		fmt.Println("移动后棋盘:")
		printBoard(board)
	}
}

func printBoard(board *chess.Board) {
	pieceNames := map[chess.PieceType]string{
		chess.Empty:    "·",
		chess.King:     "帅",
		chess.Advisor:  "士",
		chess.Elephant: "象",
		chess.Horse:    "马",
		chess.Chariot:  "车",
		chess.Cannon:   "炮",
		chess.Pawn:     "兵",
	}

	for row := 0; row < 10; row++ {
		for col := 0; col < 9; col++ {
			piece := board.Grid[row][col]
			if piece.Type == chess.Empty {
				fmt.Print("· ")
			} else {
				name := pieceNames[piece.Type]
				if piece.Color == chess.Black {
					name = "黑" + name
				} else {
					name = "红" + name
				}
				fmt.Printf("%s ", name)
			}
		}
		fmt.Println()
	}
}
