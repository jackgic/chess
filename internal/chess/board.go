package chess

import (
	"fmt"
	"sort"
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
	FirstMove Color        // 先手方颜色
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
func NewBoard(firstMove Color) *Board {
	board := &Board{
		Turn:     firstMove,
		FirstMove: firstMove,
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

	// 根据先手方设置棋子布局
	if b.FirstMove == Red {
		// 红方先手（玩家执红），红方在下方（靠近玩家）
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
	} else {
		// 黑方先手（AI执黑），黑方在下方（靠近玩家）
		// 黑方（下方）
		b.Grid[9][0] = Piece{Chariot, Black}
		b.Grid[9][1] = Piece{Horse, Black}
		b.Grid[9][2] = Piece{Elephant, Black}
		b.Grid[9][3] = Piece{Advisor, Black}
		b.Grid[9][4] = Piece{King, Black}
		b.Grid[9][5] = Piece{Advisor, Black}
		b.Grid[9][6] = Piece{Elephant, Black}
		b.Grid[9][7] = Piece{Horse, Black}
		b.Grid[9][8] = Piece{Chariot, Black}
		b.Grid[7][1] = Piece{Cannon, Black}
		b.Grid[7][7] = Piece{Cannon, Black}
		b.Grid[6][0] = Piece{Pawn, Black}
		b.Grid[6][2] = Piece{Pawn, Black}
		b.Grid[6][4] = Piece{Pawn, Black}
		b.Grid[6][6] = Piece{Pawn, Black}
		b.Grid[6][8] = Piece{Pawn, Black}

		// 红方（上方）
		b.Grid[0][0] = Piece{Chariot, Red}
		b.Grid[0][1] = Piece{Horse, Red}
		b.Grid[0][2] = Piece{Elephant, Red}
		b.Grid[0][3] = Piece{Advisor, Red}
		b.Grid[0][4] = Piece{King, Red}
		b.Grid[0][5] = Piece{Advisor, Red}
		b.Grid[0][6] = Piece{Elephant, Red}
		b.Grid[0][7] = Piece{Horse, Red}
		b.Grid[0][8] = Piece{Chariot, Red}
		b.Grid[2][1] = Piece{Cannon, Red}
		b.Grid[2][7] = Piece{Cannon, Red}
		b.Grid[3][0] = Piece{Pawn, Red}
		b.Grid[3][2] = Piece{Pawn, Red}
		b.Grid[3][4] = Piece{Pawn, Red}
		b.Grid[3][6] = Piece{Pawn, Red}
		b.Grid[3][8] = Piece{Pawn, Red}
	}
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
	// 转换为Position
	startRow, startCol := start/10, start%10
	endRow, endCol := end/10, end%10
	
	from := Position{Row: startRow, Col: startCol}
	to := Position{Row: endRow, Col: endCol}
	
	// 验证移动是否合法
	if !b.IsValidMove(from, to) {
		return fmt.Errorf("非法移动")
	}
	
	// 在移动之前生成棋谱（此时棋子还在原位置）
	moveStr := b.toStandardNotation(start, end)
	
	// 执行移动
	piece := b.Grid[from.Row][from.Col]
	b.Grid[to.Row][to.Col] = piece
	b.Grid[from.Row][from.Col] = Piece{Type: Empty, Color: None}
	
	// 记录移动
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
		// 红方：列从右到左为一到九（中文数字）
		// col 0 → 九, col 1 → 八, ..., col 8 → 一
		startPos = toChineseNumber(9 - startCol)
	} else {
		// 黑方：列从左到右为1-9（阿拉伯数字）
		startPos = strconv.Itoa(startCol + 1)
	}

	// 移动类型（根据棋子颜色判断方向）
	moveType := ""
	
	// 判断是否是斜向移动的棋子（马、相、士）
	// isDiagonalPiece := pieceType == Horse || pieceType == Elephant || pieceType == Advisor
	
	if startCol == endCol {
		// 纵向移动（车、炮、兵、将）
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
		// 添加移动步数（红方用中文数字，黑方用阿拉伯数字）
		steps := abs(startRow - endRow)
		if pieceColor == Red {
			moveType += toChineseNumber(steps)
		} else {
			moveType += strconv.Itoa(steps)
		}
	} else if startRow == endRow {
		// 纯横向移动（车、炮）
		moveType = "平"
		// 目标位置
		if pieceColor == Red {
			// col 0 → 九, col 8 → 一
			moveType += toChineseNumber(9 - endCol)
		} else {
			moveType += strconv.Itoa(endCol + 1)
		}
	} else {
		// 斜向移动（马、相、士）或兵的横向移动
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
		// 对于斜向移动，使用目标列位置
		if pieceColor == Red {
			// col 0 → 九, col 8 → 一
			moveType += toChineseNumber(9 - endCol)
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

// IsInCheck 检查指定颜色的将/帅是否被将军
func (b *Board) IsInCheck(color Color) bool {
	// 找到指定颜色的将/帅位置
	var kingPos Position
	found := false
	
	for row := 0; row < 10; row++ {
		for col := 0; col < 9; col++ {
			piece := b.Grid[row][col]
			if piece.Type == King && piece.Color == color {
				kingPos = Position{Row: row, Col: col}
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	
	if !found {
		return false // 将/帅已被吃，游戏应该已经结束
	}
	
	// 检查对方所有棋子是否能攻击到将/帅
	attackerColor := Red
	if color == Red {
		attackerColor = Black
	}
	
	for row := 0; row < 10; row++ {
		for col := 0; col < 9; col++ {
			piece := b.Grid[row][col]
			if piece.Color == attackerColor {
				from := Position{Row: row, Col: col}
				// 检查这个棋子是否能攻击到将/帅
				if b.canAttack(from, kingPos) {
					return true
				}
			}
		}
	}
	
	return false
}

// canAttack 检查棋子是否能攻击到目标位置
func (b *Board) canAttack(from, to Position) bool {
	// 边界检查
	if !b.IsInBoard(from) || !b.IsInBoard(to) {
		return false
	}

	piece := b.Grid[from.Row][from.Col]

	// 检查是否有棋子
	if piece.Type == Empty {
		return false
	}

	// 检查目标位置是否是己方棋子
	targetPiece := b.Grid[to.Row][to.Col]
	if targetPiece.Color == piece.Color {
		return false
	}

	// 根据棋子类型检查攻击规则
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
		if !b.IsValidKingMove(from, to, piece.Color) {
			return false
		}
	case Advisor:
		if !b.IsValidAdvisorMove(from, to, piece.Color) {
			return false
		}
	case Elephant:
		if !b.IsValidElephantMove(from, to, piece.Color) {
			return false
		}
	case Horse:
		if !b.IsValidHorseMove(from, to) {
			return false
		}
	case Chariot:
		if !b.IsValidChariotMove(from, to) {
			return false
		}
	case Cannon:
		if !b.IsValidCannonMove(from, to) {
			return false
		}
	case Pawn:
		if !b.IsValidPawnMove(from, to, piece.Color) {
			return false
		}
	default:
		return false
	}

	// 检查将军状态：如果当前被将军，移动必须解除将军状态
	if b.IsInCheck(b.Turn) {
		// 模拟移动，检查移动后是否仍然被将军
		originalPiece := b.Grid[to.Row][to.Col]
		b.Grid[to.Row][to.Col] = piece
		b.Grid[from.Row][from.Col] = Piece{Type: Empty, Color: None}
		
		stillInCheck := b.IsInCheck(b.Turn)
		
		// 恢复棋盘状态
		b.Grid[from.Row][from.Col] = piece
		b.Grid[to.Row][to.Col] = originalPiece
		
		// 如果移动后仍然被将军，则移动不合法
		if stillInCheck {
			return false
		}
	} else {
		// 如果不被将军，检查移动后是否会导致己方被将军
		originalPiece := b.Grid[to.Row][to.Col]
		b.Grid[to.Row][to.Col] = piece
		b.Grid[from.Row][from.Col] = Piece{Type: Empty, Color: None}
		
		inCheckAfterMove := b.IsInCheck(b.Turn)
		
		// 恢复棋盘状态
		b.Grid[from.Row][from.Col] = piece
		b.Grid[to.Row][to.Col] = originalPiece
		
		// 如果移动后会导致己方被将军，则移动不合法
		if inCheckAfterMove {
			return false
		}
	}

	return true
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

// ToString 将棋盘转换为可视化字符串
func (b *Board) ToString() string {
	var sb strings.Builder
	
	sb.WriteString("```\n")
	sb.WriteString("   列: 0 1 2 3 4 5 6 7 8\n")
	
	for row := 0; row < 10; row++ {
		sb.WriteString(fmt.Sprintf("行%d:  ", row))
		for col := 0; col < 9; col++ {
			piece := b.Grid[row][col]
			if piece.Type == Empty {
				sb.WriteString(". ")
			} else {
				pieceStr := b.PieceToFEN(piece)
				sb.WriteString(pieceStr + " ")
			}
		}
		
		// 添加行注释
		if row == 0 {
			sb.WriteString(" (黑方底线)")
		} else if row == 2 {
			sb.WriteString(" (黑方炮)")
		} else if row == 3 {
			sb.WriteString(" (黑方卒)")
		} else if row == 6 {
			sb.WriteString(" (红方兵)")
		} else if row == 7 {
			sb.WriteString(" (红方炮)")
		} else if row == 9 {
			sb.WriteString(" (红方底线)")
		}
		
		sb.WriteString("\n")
	}
	
	sb.WriteString("```")
	return sb.String()
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

// MoveByChineseNotation 根据中文棋谱执行走子
// 例如：炮2平5、马二进三、车9进1
func (b *Board) MoveByChineseNotation(notation string, color Color) error {
	if len(notation) < 4 {
		return fmt.Errorf("中文棋谱格式错误: %s", notation)
	}

	// 解析棋子名称
	pieceName := string([]rune(notation)[0])
	pieceType, err := b.parsePieceType(pieceName, color)
	if err != nil {
		return err
	}

	// 解析起始列
	fromColStr := string([]rune(notation)[1])
	fromCol, err := b.parseColumn(fromColStr, color)
	if err != nil {
		return err
	}

	// 解析动作
	action := string([]rune(notation)[2])

	// 解析目标列或步数
	targetStr := string([]rune(notation)[3])
	target, err := b.parseNumber(targetStr)
	if err != nil {
		return err
	}

	// 使用新的findPiece函数，传递完整棋谱字符串
	fromRow, err := b.findPiece(pieceType, color, fromCol, notation)
	if err != nil {
		return err
	}

	// 根据动作计算目标位置
	var toRow, toCol int
	switch action {
	case "平":
		// 平移：行不变，列改变
		toRow = fromRow
		toCol, err = b.parseColumn(targetStr, color)
		if err != nil {
			return err
		}
	case "进":
		// 前进
		// 对于直线移动的棋子（车、炮、兵、卒），列不变，数字表示步数
		if pieceType == Chariot || pieceType == Cannon || pieceType == Pawn {
			if color == Red {
				toRow = fromRow - target // 红方向上移动
			} else {
				toRow = fromRow + target // 黑方向下移动
			}
			toCol = fromCol
		} else {
			// 对于斜线移动的棋子（马、相、士），数字表示目标列
			toCol, err = b.parseColumn(targetStr, color)
			if err != nil {
				return err
			}
			// 根据目标列计算行（向前移动）
			if color == Red {
				// 红方向上移动，尝试2步或1步
				if pieceType == Horse {
					// 马走日字：如果列差为1，则行差为2；如果列差为2，则行差为1
					colDiff := abs(toCol - fromCol)
					if colDiff == 1 {
						toRow = fromRow - 2
					} else if colDiff == 2 {
						toRow = fromRow - 1
					} else {
						return fmt.Errorf("马的移动不符合日字规则")
					}
				} else if pieceType == Elephant {
					// 象走田字：列差必须为2
					colDiff := abs(toCol - fromCol)
					if colDiff != 2 {
						return fmt.Errorf("象的移动不符合田字规则，列差必须为2，实际为%d", colDiff)
					}
					toRow = fromRow - 2
				} else {
					// 士向前移动1步
					toRow = fromRow - 1
				}
			} else {
				// 黑方向下移动
				if pieceType == Horse {
					colDiff := abs(toCol - fromCol)
					if colDiff == 1 {
						toRow = fromRow + 2
					} else if colDiff == 2 {
						toRow = fromRow + 1
					} else {
						return fmt.Errorf("马的移动不符合日字规则")
					}
				} else if pieceType == Elephant {
					colDiff := abs(toCol - fromCol)
					if colDiff != 2 {
						return fmt.Errorf("象的移动不符合田字规则，列差必须为2，实际为%d", colDiff)
					}
					toRow = fromRow + 2
				} else {
					toRow = fromRow + 1
				}
			}
		}
	case "退":
		// 后退
		// 对于直线移动的棋子（车、炮、兵、卒），列不变，数字表示步数
		if pieceType == Chariot || pieceType == Cannon || pieceType == Pawn {
			if color == Red {
				toRow = fromRow + target // 红方向下移动（后退）
			} else {
				toRow = fromRow - target // 黑方向上移动（后退）
			}
			toCol = fromCol
		} else {
			// 对于斜线移动的棋子（马、相、士），数字表示目标列
			toCol, err = b.parseColumn(targetStr, color)
			if err != nil {
				return err
			}
			// 根据目标列计算行（向后移动）
			if pieceType == Horse {
				colDiff := abs(toCol - fromCol)
				if colDiff == 1 {
					if color == Red {
						toRow = fromRow + 2 // 红方向下移动
					} else {
						toRow = fromRow - 2 // 黑方向上移动
					}
				} else if colDiff == 2 {
					if color == Red {
						toRow = fromRow + 1 // 红方向下移动
					} else {
						toRow = fromRow - 1 // 黑方向上移动
					}
				} else {
					return fmt.Errorf("马的移动不符合日字规则")
				}
			} else if pieceType == Elephant {
				colDiff := abs(toCol - fromCol)
				if colDiff != 2 {
					return fmt.Errorf("象的移动不符合田字规则，列差必须为2，实际为%d", colDiff)
				}
				if color == Red {
					toRow = fromRow + 2 // 红方向下移动
				} else {
					toRow = fromRow - 2 // 黑方向上移动
				}
			} else {
				// 士向后移动1步
				if color == Red {
					toRow = fromRow + 1 // 红方向下移动
				} else {
					toRow = fromRow - 1 // 黑方向上移动
				}
			}
		}
	default:
		return fmt.Errorf("未知的动作: %s", action)
	}

	// 验证目标位置是否在棋盘内
	if toRow < 0 || toRow >= 10 || toCol < 0 || toCol >= 9 {
		return fmt.Errorf("目标位置超出棋盘范围: (%d,%d)", toRow, toCol)
	}

	// 执行走子
	from := fromRow*10 + fromCol
	to := toRow*10 + toCol
	return b.Move(from, to)
}

// parsePieceType 解析棋子类型
func (b *Board) parsePieceType(name string, color Color) (PieceType, error) {
	pieceMap := map[string]PieceType{
		"帅": King, "将": King,
		"仕": Advisor, "士": Advisor,
		"相": Elephant, "象": Elephant,
		"马": Horse,
		"车": Chariot,
		"炮": Cannon,
		"兵": Pawn, "卒": Pawn,
	}

	pieceType, ok := pieceMap[name]
	if !ok {
		return Empty, fmt.Errorf("未知的棋子名称: %s", name)
	}
	return pieceType, nil
}

// parseColumn 解析列号
func (b *Board) parseColumn(colStr string, color Color) (int, error) {
	// 数字映射
	numMap := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9,
	}

	num, ok := numMap[colStr]
	if !ok {
		return 0, fmt.Errorf("无效的列号: %s", colStr)
	}

	// 红方：列从右到左为1-9，需要转换为0-8（从左到右）
	// 黑方：列从左到右为1-9，需要转换为0-8
	if color == Red {
		return 9 - num, nil // 红方1对应列8，9对应列0
	}
	return num - 1, nil // 黑方1对应列0，9对应列8
}

// parseNumber 解析数字
func (b *Board) parseNumber(numStr string) (int, error) {
	numMap := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9,
	}

	num, ok := numMap[numStr]
	if !ok {
		return 0, fmt.Errorf("无效的数字: %s", numStr)
	}
	return num, nil
}

// findPiece 查找指定类型和颜色的棋子在指定列的位置
func (b *Board) findPiece(pieceType PieceType, color Color, col int, notation string) (int, error) {
	// 收集同一列所有符合条件的棋子
	var candidates []int
	for row := 0; row < 10; row++ {
		if b.Grid[row][col].Type == pieceType && b.Grid[row][col].Color == color {
			candidates = append(candidates, row)
		}
	}

	// 根据棋子数量处理
	switch len(candidates) {
	case 0:
		return -1, fmt.Errorf("未找到棋子: type=%d, color=%d, col=%d", pieceType, color, col)
	case 1:
		return candidates[0], nil
	default:
		// 同列有多个相同棋子，需要根据棋谱中的"前"、"后"定位
		return b.resolveAmbiguousPiece(candidates, notation, color)
	}
}

func (b *Board) resolveAmbiguousPiece(candidates []int, notation string, color Color) (int, error) {
	// 排序候选棋子位置（红方从上到下，黑方从下到上）
	if color == Red {
		sort.Ints(candidates)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(candidates)))
	}

	// 解析棋谱中的位置修饰词（前/后）
	if strings.Contains(notation, "前") {
		return candidates[0], nil // 最前面的棋子
	} else if strings.Contains(notation, "后") {
		return candidates[len(candidates)-1], nil // 最后面的棋子
	}

	// 默认返回第一个找到的棋子（兼容旧行为）
	return candidates[0], nil
}
