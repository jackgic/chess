package game

import (
	"fmt"
	"sync"

	"chinese-chess-ai/internal/chess"
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/lke"

	"git.woa.com/gocdb/base/uuid"
)

// Game 游戏实例
type Game struct {
	ID          string
	Board       *chess.Board
	LKEClient   *lke.Client
	PlayerColor chess.Color
	AIColor     chess.Color
	Status      GameStatus
	SessionID   string // 新增会话ID字段
	mu          sync.RWMutex
}

// GameStatus 游戏状态
type GameStatus string

const (
	StatusPlaying  GameStatus = "playing"
	StatusRedWin   GameStatus = "red_win"
	StatusBlackWin GameStatus = "black_win"
	StatusDraw     GameStatus = "draw"
)

// Manager 游戏管理器
type Manager struct {
	games  map[string]*Game
	config *config.Config
	mu     sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取游戏管理器单例
func GetManager(cfg *config.Config) *Manager {
	once.Do(func() {
		instance = &Manager{
			games:  make(map[string]*Game),
			config: cfg,
		}
	})
	return instance
}

// NewGame 创建新游戏
func (m *Manager) NewGame(playerColor chess.Color) (*Game, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 创建LKE客户端
	lkeClient, err := lke.NewClient(&m.config.TencentCloud)
	if err != nil {
		return nil, fmt.Errorf("创建LKE客户端失败: %v", err)
	}

	// 生成游戏ID
	gameID := fmt.Sprintf("game_%d", len(m.games)+1)

	// 生成唯一会话ID (使用UUIDv4)
	sessionID := fmt.Sprintf("session_%s", uuid.New().String())

	// 确定AI颜色
	aiColor := chess.Black
	if playerColor == chess.Black {
		aiColor = chess.Red
	}

	game := &Game{
		ID:          gameID,
		Board:       chess.NewBoard(),
		LKEClient:   lkeClient,
		PlayerColor: playerColor,
		AIColor:     aiColor,
		Status:      StatusPlaying,
		SessionID:   sessionID, // 设置会话ID
	}

	// 通知智能体开始新对局
	if err := lkeClient.StartSession(sessionID); err != nil {
		return nil, fmt.Errorf("通知智能体开始新对局失败: %v", err)
	}

	m.games[gameID] = game

	return game, nil
}

// GetGame 获取游戏
func (m *Manager) GetGame(gameID string) (*Game, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.games[gameID]
	if !exists {
		return nil, fmt.Errorf("游戏不存在")
	}

	return game, nil
}

// PlayerMove 玩家走子
func (g *Game) PlayerMove(fromRow, fromCol, toRow, toCol int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusPlaying {
		return fmt.Errorf("游戏已结束")
	}

	// 检查是否是玩家回合
	if g.Board.Turn != g.PlayerColor {
		return fmt.Errorf("不是你的回合")
	}

	// 执行移动
	from := chess.Position{Row: fromRow, Col: fromCol}
	to := chess.Position{Row: toRow, Col: toCol}

	if err := g.Board.Move(from.ToInt(), to.ToInt()); err != nil {
		return err
	}

	// 切换回合
	if g.Board.Turn == chess.Red {
		g.Board.Turn = chess.Black
	} else {
		g.Board.Turn = chess.Red
	}

	// 检查游戏是否结束
	g.checkGameOver()

	return nil
}

// AIMove AI走子
func (g *Game) AIMove() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusPlaying {
		return "", fmt.Errorf("游戏已结束")
	}

	if g.Board.Turn != g.AIColor {
		return "", fmt.Errorf("不是AI的回合")
	}

	// 构建包含棋盘状态的提示词
	prompt := fmt.Sprintf(`当前棋盘状态（FEN格式）：%s

当前轮到：%s
你是：%s

历史走子：
%s

请分析当前局面并走子。必须严格按照以下格式输出：
MOVE: 棋子名称起始位置移动类型目标位置

例如：
- MOVE: 炮二平五
- MOVE: 马8进7
- MOVE: 车九进一

请立即给出你的走子：`,
		g.Board.ToFEN(),
		colorToString(g.Board.Turn),
		colorToString(g.AIColor),
		formatMoveHistory(g.Board.MoveList))

	// 使用会话ID与智能体交互
	answer, err := g.LKEClient.Chat(g.SessionID, prompt)
	if err != nil {
		return "", fmt.Errorf("LKE调用失败: %v", err)
	}

	// 解析AI走子指令
	move, err := g.LKEClient.ExtractMove(answer)
	if err != nil {
		return "", fmt.Errorf("解析走子失败: %v", err)
	}

	// 执行走子
	if err := g.Board.Move(move.From, move.To); err != nil {
		return "", fmt.Errorf("执行走子失败: %v", err)
	}

	// 切换回合
	if g.Board.Turn == chess.Red {
		g.Board.Turn = chess.Black
	} else {
		g.Board.Turn = chess.Red
	}

	// 检查游戏状态
	g.checkGameOver()

	return answer, nil
}

// checkGameOver 检查游戏是否结束
func (g *Game) checkGameOver() {
	// 检查是否有将帅被吃
	hasRedKing := false
	hasBlackKing := false

	for row := 0; row < 10; row++ {
		for col := 0; col < 9; col++ {
			piece := g.Board.Grid[row][col]
			if piece.Type == chess.King {
				if piece.Color == chess.Red {
					hasRedKing = true
				} else {
					hasBlackKing = true
				}
			}
		}
	}

	if !hasRedKing {
		g.Status = StatusBlackWin
	} else if !hasBlackKing {
		g.Status = StatusRedWin
	}

	// 检查是否无子可走
	legalMoves := g.Board.GetAllLegalMoves(g.Board.Turn)
	if len(legalMoves) == 0 {
		if g.Board.Turn == chess.Red {
			g.Status = StatusBlackWin
		} else {
			g.Status = StatusRedWin
		}
	}
}

// GetState 获取游戏状态
func (g *Game) GetState() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// 转换棋盘为前端格式
	board := make([][]map[string]interface{}, 10)
	for i := 0; i < 10; i++ {
		board[i] = make([]map[string]interface{}, 9)
		for j := 0; j < 9; j++ {
			piece := g.Board.Grid[i][j]
			board[i][j] = map[string]interface{}{
				"type":  int(piece.Type),
				"color": int(piece.Color),
				"name":  g.Board.PieceToString(piece),
			}
		}
	}

	return map[string]interface{}{
		"id":          g.ID,
		"board":       board,
		"turn":        int(g.Board.Turn),
		"status":      string(g.Status),
		"playerColor": int(g.PlayerColor),
		"aiColor":     int(g.AIColor),
		"moveList":    g.Board.MoveList,
	}
}

// colorToString 将颜色转换为字符串
func colorToString(color chess.Color) string {
	switch color {
	case chess.Red:
		return "红方"
	case chess.Black:
		return "黑方"
	default:
		return "未知"
	}
}

// formatMoveHistory 格式化走子历史
func formatMoveHistory(history []string) string {
	if len(history) == 0 {
		return "无"
	}

	result := ""
	for i, move := range history {
		if i > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("%d. %s", i+1, move)
	}
	return result
}
