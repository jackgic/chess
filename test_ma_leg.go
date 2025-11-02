package main

import (
	"chinese-chess-ai/internal/chess"
	"fmt"
)

func main() {
	fmt.Println("=== 测试马腿问题 ===\n")

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

	// 马在 (2, 6)，要移动到 (4, 7)
	fromRow, fromCol := 2, 6
	toRow, toCol := 4, 7
	
	fmt.Printf("马的位置: (%d, %d)\n", fromRow, fromCol)
	fmt.Printf("目标位置: (%d, %d)\n", toRow, toCol)
	
	rowDiff := toRow - fromRow
	colDiff := toCol - fromCol
	fmt.Printf("行差: %d, 列差: %d\n", rowDiff, colDiff)
	
	// 计算马腿位置
	var legRow, legCol int
	if abs(rowDiff) == 2 {
		legRow = (fromRow + toRow) / 2
		legCol = fromCol
		fmt.Printf("行差为2，马腿在: (%d, %d)\n", legRow, legCol)
	} else {
		legRow = fromRow
		legCol = (fromCol + toCol) / 2
		fmt.Printf("列差为2，马腿在: (%d, %d)\n", legRow, legCol)
	}
	
	// 检查马腿位置的棋子
	piece := board.Grid[legRow][legCol]
	if piece.Type == chess.Empty {
		fmt.Println("马腿位置: 空位 ✅")
	} else {
		pieceNames := map[chess.PieceType]string{
			chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
			chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
		}
		colorNames := map[chess.Color]string{
			chess.Red: "红方", chess.Black: "黑方",
		}
		fmt.Printf("马腿位置: %s %s ❌ (蹩马腿)\n", colorNames[piece.Color], pieceNames[piece.Type])
	}
	
	// 打印周围的棋子
	fmt.Println("\n周围的棋子:")
	for r := fromRow - 1; r <= fromRow + 3; r++ {
		if r < 0 || r >= 10 {
			continue
		}
		for c := fromCol - 1; c <= fromCol + 2; c++ {
			if c < 0 || c >= 9 {
				continue
			}
			p := board.Grid[r][c]
			if p.Type == chess.Empty {
				fmt.Printf("  (%d, %d): 空\n", r, c)
			} else {
				pieceNames := map[chess.PieceType]string{
					chess.King: "将", chess.Advisor: "士", chess.Elephant: "象",
					chess.Horse: "马", chess.Chariot: "车", chess.Cannon: "炮", chess.Pawn: "卒",
				}
				colorNames := map[chess.Color]string{
					chess.Red: "红方", chess.Black: "黑方",
				}
				fmt.Printf("  (%d, %d): %s %s\n", r, c, colorNames[p.Color], pieceNames[p.Type])
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
