package chess

import (
	"fmt"
	"strconv"
	"strings"
)

// PieceType 棋子类型
type PieceType int

const (
	Empty    PieceType = iota
	King               // 将/帅
	Advisor            // 士
	Elephant           // 象
	Horse              // 马
	Chariot            // 车
	Cannon             // 炮
	Pawn               // 兵/卒
)

// Color 棋子颜色
type Color int

const (
	None  Color = iota
	Red         // 红方
	Black       // 黑方
)

// Piece 棋子
type Piece struct {
	Type  PieceType
	Color Color
}

// Board 棋盘
type Board struct {
	Grid     [10][9]Piece // 10行9列
	Turn     Color        // 当前回合
	MoveList []string     // 移动历史
	Moves    []Move       // 移动列表
}

// Position 位置
type Position struct {
	Row int
	Col int
}

// ToInt 将Position转换为int表示 (row*10 + col)
func (p Position) ToInt() int {
	return p.Row*10 + p.Col
}

// Move 移动结构
type Move struct {
	FromX, FromY, ToX, ToY int
}

// NewBoard 创建新棋盘
func NewBoard() *Board {
	board := &Board{
		Turn:     Red,
		MoveList: []string{},
	}
	board.InitBoard()
	return board
}

// InitBoard 初始化棋盘
func (b *Board) InitBoard() {
	// 清空棋盘
	for i := 0; i < 10; i++ {
		for j := 0; j < 9; j++ {
			b.Grid[i][j] = Piece{Empty, None}
		}
	}

	// 黑方（上方）
	b.Grid[0][0] = Piece{Chariot, Black}
	b.Grid[0][1] = Piece{Horse, Black}
	b.Grid[0][2] = Piece{Elephant, Black}
	b.Grid[0][3] = Piece{Advisor, Black}
	b.Grid[0][4] = Piece{King, Black}
	b.Grid[0][5] = Piece{Advisor, Black}
	b.Grid[0][6] = Piece{Elephant, Black}
	b.Grid[0][7] = Piece{Horse, Black}
	b.Grid[0][8] = Piece{Chariot, Black}
	b.Grid[2][1] = Piece{Cannon, Black}
	b.Grid[2][7] = Piece{Cannon, Black}
	b.Grid[3][0] = Piece{Pawn, Black}
	b.Grid[3][2] = Piece{Pawn, Black}
	b.Grid[3][4] = Piece{Pawn, Black}
	b.Grid[3][6] = Piece{Pawn, Black}
	b.Grid[3][8] = Piece{Pawn, Black}

	// 红方（下方）
	b.Grid[9][0] = Piece{Chariot, Red}
	b.Grid[9][1] = Piece{Horse, Red}
	b.Grid[9][2] = Piece{Elephant, Red}
	b.Grid[9][3] = Piece{Advisor, Red}
	b.Grid[9][4] = Piece{King, Red}
	b.Grid[9][5] = Piece{Advisor, Red}
	b.Grid[9][6] = Piece{Elephant, Red}
	b.Grid[9][7] = Piece{Horse, Red}
	b.Grid[9][8] = Piece{Chariot, Red}
	b.Grid[7][1] = Piece{Cannon, Red}
	b.Grid[7][7] = Piece{Cannon, Red}
	b.Grid[6][0] = Piece{Pawn, Red}
	b.Grid[6][2] = Piece{Pawn, Red}
	b.Grid[6][4] = Piece{Pawn, Red}
	b.Grid[6][6] = Piece{Pawn, Red}
	b.Grid[6][8] = Piece{Pawn, Red}
}

// 将数字坐标转换为标准棋谱格式
func toStandardNotation(fromX, fromY, toX, toY int, piece Piece) string {
	// 列号转换
	colNames := [9]string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}

	var pieceName, fromCol, toCol string
	if piece.Color == Red {
		pieceName = pieceNamesRed[piece.Type]
		fromCol = colNames[fromY]
		toCol = colNames[toY]
	} else {
		pieceName = pieceNamesBlack[piece.Type]
		fromCol = strconv.Itoa(fromY + 1)
		toCol = strconv.Itoa(toY + 1)
	}

	// 移动方向判断
	moveType := "平"
	if fromX != toX {
		if (piece.Color == Red && toX > fromX) || (piece.Color == Black && toX < fromX) {
			moveType = "进"
		} else {
			moveType = "退"
		}
	}

	return fmt.Sprintf("%s%s%s%s", pieceName, fromCol, moveType, toCol)
}

// Move 移动棋子
func (b *Board) Move(start, end int) error {
	// TODO: 实现移动逻辑
	// 这里应该包含移动验证、执行移动等逻辑

	// 记录移动（修改为生成标准棋谱格式）
	moveStr := b.toStandardNotation(start, end)
	b.MoveList = append(b.MoveList, moveStr)

	return nil
}

// 新增函数：将数字坐标转换为标准棋谱表示
func (b *Board) toStandardNotation(start, end int) string {
	startRow, startCol := start/10, start%10
	endRow, endCol := end/10, end%10

	// 获取棋子
	piece := b.Grid[startRow][startCol]
	pieceType := piece.Type
	pieceColor := piece.Color

	// 棋子名称
	pieceName := ""
	switch pieceType {
	case King:
		if pieceColor == Red {
			pieceName = "帅"
		} else {
			pieceName = "将"
		}
	case Advisor:
		if pieceColor == Red {
			pieceName = "仕"
		} else {
			pieceName = "士"
		}
	case Elephant:
		if pieceColor == Red {
			pieceName = "相"
		} else {
			pieceName = "象"
		}
	case Horse:
		pieceName = "马"
	case Chariot:
		pieceName = "车"
	case Cannon:
		pieceName = "炮"
	case Pawn:
		if pieceColor == Red {
			pieceName = "兵"
		} else {
			pieceName = "卒"
		}
	}

	// 起始位置（红方用中文数字，黑方用阿拉伯数字）
	startPos := ""
	if pieceColor == Red {
		// 红方：列从右到左为1-9（中文数字）
		startPos = toChineseNumber(8 - startCol + 1)
	} else {
		// 黑方：列从左到右为1-9（阿拉伯数字）
		startPos = strconv.Itoa(startCol + 1)
	}

	// 移动类型（根据棋子颜色判断方向）
	moveType := ""
	if startCol == endCol {
		// 纵向移动
		if pieceColor == Red {
			// 红方：向上移动为"进"，向下移动为"退"
			if endRow < startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		} else {
			// 黑方：向下移动为"进"，向上移动为"退"
			if endRow > startRow {
				moveType = "进"
			} else {
				moveType = "退"
			}
		}
		// 添加移动步数
		moveType += strconv.Itoa(abs(startRow - endRow))
	} else {
		// 横向移动
		moveType = "平"
		// 目标位置
		if pieceColor == Red {
			moveType += toChineseNumber(8 - endCol + 1)
		} else {
			moveType += strconv.Itoa(endCol + 1)
		}
	}

	return pieceName + startPos + moveType
}

// 辅助函数：数字转中文
func toChineseNumber(num int) string {
	chineseNumbers := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
	if num >= 1 && num <= 9 {
		return chineseNumbers[num-1]
	}
	return strconv.Itoa(num)
}

// 辅助函数：绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IsValidMove 检查移动是否合法
func (b *Board) IsValidMove(from, to Position) bool {
	// 边界检查
	if !b.IsInBoard(from) || !b.IsInBoard(to) {
		return false
	}

	piece := b.Grid[from.Row][from.Col]

	// 检查是否有棋子
	if piece.Type == Empty {
		return false
	}

	// 检查是否是当前回合的棋子
	if piece.Color != b.Turn {
		return false
	}

	// 检查目标位置是否是己方棋子
	targetPiece := b.Grid[to.Row][to.Col]
	if targetPiece.Color == piece.Color {
		return false
	}

	// 根据棋子类型检查移动规则
	switch piece.Type {
	case King:
		return b.IsValidKingMove(from, to, piece.Color)
	case Advisor:
		return b.IsValidAdvisorMove(from, to, piece.Color)
	case Elephant:
		return b.IsValidElephantMove(from, to, piece.Color)
	case Horse:
		return b.IsValidHorseMove(from, to)
	case Chariot:
		return b.IsValidChariotMove(from, to)
	case Cannon:
		return b.IsValidCannonMove(from, to)
	case Pawn:
		return b.IsValidPawnMove(from, to, piece.Color)
	}

	return false
}

// IsInBoard 检查位置是否在棋盘内
func (b *Board) IsInBoard(pos Position) bool {
	return pos.Row >= 0 && pos.Row < 10 && pos.Col >= 0 && pos.Col < 9
}

// IsValidKingMove 将/帅移动规则
func (b *Board) IsValidKingMove(from, to Position, color Color) bool {
	// 只能在九宫格内移动
	if color == Red {
		if to.Row < 7 || to.Row > 9 || to.Col < 3 || to.Col > 5 {
			return false
		}
	} else {
		if to.Row < 0 || to.Row > 2 || to.Col < 3 || to.Col > 5 {
			return false
		}
	}

	// 只能移动一格
	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)
	return (rowDiff == 1 && colDiff == 0) || (rowDiff == 0 && colDiff == 1)
}

// IsValidAdvisorMove 士移动规则
func (b *Board) IsValidAdvisorMove(from, to Position, color Color) bool {
	// 只能在九宫格内移动
	if color == Red {
		if to.Row < 7 || to.Row > 9 || to.Col < 3 || to.Col > 5 {
			return false
		}
	} else {
		if to.Row < 0 || to.Row > 2 || to.Col < 3 || to.Col > 5 {
			return false
		}
	}

	// 只能斜着走一格
	return abs(to.Row-from.Row) == 1 && abs(to.Col-from.Col) == 1
}

// IsValidElephantMove 象移动规则
func (b *Board) IsValidElephantMove(from, to Position, color Color) bool {
	// 不能过河
	if color == Red && to.Row < 5 {
		return false
	}
	if color == Black && to.Row > 4 {
		return false
	}

	// 必须走田字
	if abs(to.Row-from.Row) != 2 || abs(to.Col-from.Col) != 2 {
		return false
	}

	// 检查象眼
	eyeRow := (from.Row + to.Row) / 2
	eyeCol := (from.Col + to.Col) / 2
	return b.Grid[eyeRow][eyeCol].Type == Empty
}

// IsValidHorseMove 马移动规则
func (b *Board) IsValidHorseMove(from, to Position) bool {
	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)

	// 必须走日字
	if !((rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)) {
		return false
	}

	// 检查马腿
	var legRow, legCol int
	if rowDiff == 2 {
		legRow = (from.Row + to.Row) / 2
		legCol = from.Col
	} else {
		legRow = from.Row
		legCol = (from.Col + to.Col) / 2
	}

	return b.Grid[legRow][legCol].Type == Empty
}

// IsValidChariotMove 车移动规则
func (b *Board) IsValidChariotMove(from, to Position) bool {
	// 必须直线移动
	if from.Row != to.Row && from.Col != to.Col {
		return false
	}

	// 检查路径是否有障碍
	if from.Row == to.Row {
		start := min(from.Col, to.Col) + 1
		end := max(from.Col, to.Col)
		for col := start; col < end; col++ {
			if b.Grid[from.Row][col].Type != Empty {
				return false
			}
		}
	} else {
		start := min(from.Row, to.Row) + 1
		end := max(from.Row, to.Row)
		for row := start; row < end; row++ {
			if b.Grid[row][from.Col].Type != Empty {
				return false
			}
		}
	}

	return true
}

// IsValidCannonMove 炮移动规则
func (b *Board) IsValidCannonMove(from, to Position) bool {
	// 必须直线移动
	if from.Row != to.Row && from.Col != to.Col {
		return false
	}

	targetPiece := b.Grid[to.Row][to.Col]
	pieceCount := 0

	// 计算路径上的棋子数量
	if from.Row == to.Row {
		start := min(from.Col, to.Col) + 1
		end := max(from.Col, to.Col)
		for col := start; col < end; col++ {
			if b.Grid[from.Row][col].Type != Empty {
				pieceCount++
			}
		}
	} else {
		start := min(from.Row, to.Row) + 1
		end := max(from.Row, to.Row)
		for row := start; row < end; row++ {
			if b.Grid[row][from.Col].Type != Empty {
				pieceCount++
			}
		}
	}

	// 吃子时必须隔一个棋子，移动时路径必须为空
	if targetPiece.Type != Empty {
		return pieceCount == 1
	}
	return pieceCount == 0
}

// IsValidPawnMove 兵/卒移动规则
func (b *Board) IsValidPawnMove(from, to Position, color Color) bool {
	rowDiff := to.Row - from.Row
	colDiff := abs(to.Col - from.Col)

	// 只能移动一格
	if abs(rowDiff) > 1 || colDiff > 1 || (abs(rowDiff) == 1 && colDiff == 1) {
		return false
	}

	// 红方兵
	if color == Red {
		// 未过河只能向前
		if from.Row > 4 {
			return rowDiff == -1 && colDiff == 0
		}
		// 过河后可以左右移动
		return (rowDiff == -1 && colDiff == 0) || (rowDiff == 0 && colDiff == 1)
	}

	// 黑方卒
	// 未过河只能向前
	if from.Row < 5 {
		return rowDiff == 1 && colDiff == 0
	}
	// 过河后可以左右移动
	return (rowDiff == 1 && colDiff == 0) || (rowDiff == 0 && colDiff == 1)
}

// GetAllLegalMoves 获取所有合法移动
func (b *Board) GetAllLegalMoves(color Color) []string {
	var moves []string

	for row := 0; row < 10; row++ {
		for col := 0; col < 9; col++ {
			piece := b.Grid[row][col]
			if piece.Color != color {
				continue
			}

			from := Position{row, col}
			for toRow := 0; toRow < 10; toRow++ {
				for toCol := 0; toCol < 9; toCol++ {
					to := Position{toRow, toCol}
					if b.IsValidMove(from, to) {
						moveStr := fmt.Sprintf("%d%d-%d%d", from.Row, from.Col, to.Row, to.Col)
						moves = append(moves, moveStr)
					}
				}
			}
		}
	}

	return moves
}

// ToFEN 转换为FEN格式（简化版）
func (b *Board) ToFEN() string {
	var fen strings.Builder

	for row := 0; row < 10; row++ {
		emptyCount := 0
		for col := 0; col < 9; col++ {
			piece := b.Grid[row][col]
			if piece.Type == Empty {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fen.WriteString(fmt.Sprintf("%d", emptyCount))
					emptyCount = 0
				}
				fen.WriteString(b.PieceToFEN(piece))
			}
		}
		if emptyCount > 0 {
			fen.WriteString(fmt.Sprintf("%d", emptyCount))
		}
		if row < 9 {
			fen.WriteString("/")
		}
	}

	if b.Turn == Red {
		fen.WriteString(" r")
	} else {
		fen.WriteString(" b")
	}

	return fen.String()
}

// PieceToFEN 棋子转FEN字符
func (b *Board) PieceToFEN(piece Piece) string {
	chars := map[PieceType]string{
		King:     "k",
		Advisor:  "a",
		Elephant: "e",
		Horse:    "h",
		Chariot:  "r",
		Cannon:   "c",
		Pawn:     "p",
	}

	char := chars[piece.Type]
	if piece.Color == Red {
		return strings.ToUpper(char)
	}
	return char
}

// PieceToString 棋子转中文
func (b *Board) PieceToString(piece Piece) string {
	if piece.Color == Red {
		switch piece.Type {
		case King:
			return "帅"
		case Advisor:
			return "仕"
		case Elephant:
			return "相"
		case Horse:
			return "马"
		case Chariot:
			return "车"
		case Cannon:
			return "炮"
		case Pawn:
			return "兵"
		}
	} else {
		switch piece.Type {
		case King:
			return "将"
		case Advisor:
			return "士"
		case Elephant:
			return "象"
		case Horse:
			return "马"
		case Chariot:
			return "车"
		case Cannon:
			return "炮"
		case Pawn:
			return "卒"
		}
	}
	return ""
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 添加棋子名称映射
var pieceNamesRed = []string{"帅", "仕", "相", "马", "车", "炮", "兵"}
var pieceNamesBlack = []string{"将", "士", "象", "马", "车", "炮", "卒"}

// 添加掩码常量
const (
	PieceMask = 0x0F
	ColorMask = 0xF0
)
