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

	// 中国象棋规则：红方总是先手
	firstMove := chess.Red

	game := &Game{
		ID:          gameID,
		Board:       chess.NewBoard(firstMove),
		LKEClient:   lkeClient,
		PlayerColor: playerColor,
		AIColor:     aiColor,
		Status:      StatusPlaying,
		SessionID:   sessionID, // 设置会话ID
	}

	// 添加调试日志
	fmt.Printf("[Game] 创建新游戏 %s: 玩家颜色=%d, AI颜色=%d, 棋盘当前轮次=%d\n", 
		gameID, int(playerColor), int(aiColor), int(game.Board.Turn))

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

	// 检查是否被将军
	isInCheck := g.Board.IsInCheck(g.PlayerColor)
	
	// 执行移动
	from := chess.Position{Row: fromRow, Col: fromCol}
	to := chess.Position{Row: toRow, Col: toCol}

	if err := g.Board.Move(from.ToInt(), to.ToInt()); err != nil {
		// 提供更明确的错误信息
		if isInCheck {
			return fmt.Errorf("被将军！必须应将：%v", err)
		}
		return fmt.Errorf("非法移动：%v", err)
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
	// 添加调试日志
	fmt.Printf("[Game] AIMove检查: 棋盘当前轮次=%d, AI颜色=%d\n", int(g.Board.Turn), int(g.AIColor))
	
	// 检查是否轮到AI走子
	if g.Board.Turn != g.AIColor {
		return "", fmt.Errorf("当前不是AI的回合，轮到%s走子", 
			map[chess.Color]string{chess.Red: "红方", chess.Black: "黑方"}[g.Board.Turn])
	}

	const maxRetries = 3
	retryCount := 0
	var lastAnswer string
	var lastError error

	for retryCount < maxRetries {
		var prompt string
		if len(g.Board.MoveList) == 0 {
			prompt = "开始"
		} else {
			lastMove := g.Board.MoveList[len(g.Board.MoveList)-1]
			prompt = lastMove
		}

		// 添加错误信息到提示（如果有）
		if lastError != nil {
			prompt = fmt.Sprintf("你刚才走了：%s，但执行失败：%v，请重新走子", lastAnswer, lastError)
		}

		// 使用会话ID与智能体交互
		fmt.Printf("[Game] 发送棋谱给AI: %s\n", prompt)
		response, err := g.LKEClient.ChatWithDetails(g.SessionID, prompt)
		if err != nil {
			fmt.Printf("[Game] LKE调用失败: %v\n", err)
			return "", fmt.Errorf("LKE调用失败: %v", err)
		}
		// 打印AI的思考过程
		fmt.Printf("[Game] AI思考过程: %s\n", response.ThoughtProcess)
		// 打印AI的完整回包内容
		fmt.Printf("[Game] AI完整回包: %s\n", response.FullResponse)
		// 打印AI分析内容
		fmt.Printf("[Game] AI分析: %s\n", response.AIAnalysis)
		// 打印走子指令
		fmt.Printf("[Game] 走子指令: %s\n", response.MoveInstruction)

		// 使用解析出的走子指令
		chineseMove := response.MoveInstruction
		if chineseMove == "" {
			// 如果没有解析到走子指令，尝试从完整回复中提取
			var err error
			chineseMove, err = g.LKEClient.ExtractChineseMove(response.FullResponse)
			if err != nil {
				fmt.Printf("[Game] 解析中文棋谱失败: %v\n", err)
				return response.AIAnalysis, fmt.Errorf("解析中文棋谱失败: %v，AI回复: %s", err, response.FullResponse)
			}
		}
		fmt.Printf("[Game] 解析到的中文棋谱: %s\n", chineseMove)

		// 尝试执行走子
		if err := g.Board.MoveByChineseNotation(chineseMove, g.AIColor); err != nil {
			fmt.Printf("[Game] 执行中文棋谱失败: %v\n", err)
			lastError = fmt.Errorf("执行中文棋谱失败: %v，AI回复: %s", err, response.FullResponse)
			lastAnswer = response.FullResponse
			retryCount++
			continue
		}

		fmt.Printf("[Game] 走子执行成功\n")
		// 切换回合
		if g.Board.Turn == chess.Red {
			g.Board.Turn = chess.Black
		} else {
			g.Board.Turn = chess.Red
		}

		// 检查游戏状态
		g.checkGameOver()
		// 返回AI分析内容，如果没有则返回思考过程
		if response.AIAnalysis != "" {
			return response.AIAnalysis, nil
		}
		return response.ThoughtProcess, nil
	}

	return "", fmt.Errorf("AI走子失败超过最大重试次数: %v", lastError)
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
		fmt.Printf("[Game] 游戏结束：红方将帅被吃，黑方获胜\n")
		return
	} else if !hasBlackKing {
		g.Status = StatusRedWin
		fmt.Printf("[Game] 游戏结束：黑方将帅被吃，红方获胜\n")
		return
	}

	// 检查是否无子可走（包括被将死的情况）
	legalMoves := g.Board.GetAllLegalMoves(g.Board.Turn)
	if len(legalMoves) == 0 {
		// 检查是否被将军
		isInCheck := g.Board.IsInCheck(g.Board.Turn)
		
		if isInCheck {
			// 被将死
			if g.Board.Turn == chess.Red {
				g.Status = StatusBlackWin
				fmt.Printf("[Game] 游戏结束：红方被将死，黑方获胜\n")
			} else {
				g.Status = StatusRedWin
				fmt.Printf("[Game] 游戏结束：黑方被将死，红方获胜\n")
			}
		} else {
			// 困毙（无子可走但未被将军）
			g.Status = StatusDraw
			fmt.Printf("[Game] 游戏结束：困毙，和棋\n")
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

	// 检查将军状态
	playerInCheck := g.Board.IsInCheck(g.PlayerColor)
	aiInCheck := g.Board.IsInCheck(g.AIColor)

	return map[string]interface{}{
		"id":          g.ID,
		"board":       board,
		"turn":        int(g.Board.Turn),
		"status":      string(g.Status),
		"playerColor": int(g.PlayerColor),
		"aiColor":     int(g.AIColor),
		"moveList":    g.Board.MoveList,
		"inCheck": map[string]bool{
			"player": playerInCheck,
			"ai":     aiInCheck,
		},
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
