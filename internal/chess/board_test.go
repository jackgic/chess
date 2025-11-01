package chess

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// 测试初始化
	if board.Turn != Red {
		t.Error("初始回合应该是红方")
	}

	// 测试红方帅的位置
	if board.Grid[9][4].Type != King || board.Grid[9][4].Color != Red {
		t.Error("红方帅位置不正确")
	}

	// 测试黑方将的位置
	if board.Grid[0][4].Type != King || board.Grid[0][4].Color != Black {
		t.Error("黑方将位置不正确")
	}
}

func TestValidMove(t *testing.T) {
	board := NewBoard()

	// 测试兵的合法移动
	from := Position{Row: 6, Col: 0}
	to := Position{Row: 5, Col: 0}

	if !board.IsValidMove(from, to) {
		t.Error("兵向前移动应该是合法的")
	}

	// 测试非法移动（兵不能后退）
	from = Position{Row: 6, Col: 0}
	to = Position{Row: 7, Col: 0}

	if board.IsValidMove(from, to) {
		t.Error("兵不能后退")
	}
}

func TestMove(t *testing.T) {
	board := NewBoard()

	// 移动红方兵
	from := Position{Row: 6, Col: 0}
	to := Position{Row: 5, Col: 0}

	err := board.Move(from, to)
	if err != nil {
		t.Errorf("移动失败: %v", err)
	}

	// 检查移动后的状态
	if board.Grid[5][0].Type != Pawn {
		t.Error("移动后目标位置应该有兵")
	}

	if board.Grid[6][0].Type != Empty {
		t.Error("移动后起始位置应该为空")
	}

	// 检查回合切换
	if board.Turn != Black {
		t.Error("移动后应该切换到黑方")
	}
}

func TestChariotMove(t *testing.T) {
	board := NewBoard()

	// 清空路径，测试车的移动
	board.Grid[9][0] = Piece{Chariot, Red}
	board.Grid[8][0] = Piece{Empty, None}
	board.Grid[7][0] = Piece{Empty, None}
	board.Grid[6][0] = Piece{Empty, None}

	from := Position{Row: 9, Col: 0}
	to := Position{Row: 6, Col: 0}

	if !board.IsValidMove(from, to) {
		t.Error("车的直线移动应该是合法的")
	}
}

func TestHorseMove(t *testing.T) {
	board := NewBoard()

	// 测试马的移动（日字）
	from := Position{Row: 9, Col: 1}
	to := Position{Row: 7, Col: 2}

	if !board.IsValidMove(from, to) {
		t.Error("马的日字移动应该是合法的")
	}

	// 测试马腿被堵
	board.Grid[8][1] = Piece{Pawn, Red}
	if board.IsValidMove(from, to) {
		t.Error("马腿被堵时不应该能移动")
	}
}

func TestCannonMove(t *testing.T) {
	board := NewBoard()

	// 测试炮的移动（不吃子）
	from := Position{Row: 7, Col: 1}
	to := Position{Row: 5, Col: 1}

	if !board.IsValidMove(from, to) {
		t.Error("炮的直线移动应该是合法的")
	}

	// 测试炮吃子（需要隔一个棋子）
	board.Grid[5][1] = Piece{Pawn, Red}   // 炮架
	board.Grid[3][1] = Piece{Pawn, Black} // 目标

	from = Position{Row: 7, Col: 1}
	to = Position{Row: 3, Col: 1}

	if !board.IsValidMove(from, to) {
		t.Error("炮隔子吃子应该是合法的")
	}
}

func TestGetAllLegalMoves(t *testing.T) {
	board := NewBoard()

	moves := board.GetAllLegalMoves(Red)

	if len(moves) == 0 {
		t.Error("开局红方应该有合法移动")
	}

	// 红方开局应该有：2个兵向前，2个马跳
	if len(moves) < 4 {
		t.Errorf("红方开局至少应该有4个合法移动，实际有%d个", len(moves))
	}
}

func TestFEN(t *testing.T) {
	board := NewBoard()
	fen := board.ToFEN()

	if fen == "" {
		t.Error("FEN字符串不应该为空")
	}

	// 检查FEN格式
	if len(fen) < 10 {
		t.Error("FEN字符串格式不正确")
	}
}
