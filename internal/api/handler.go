package api

import (
	"net/http"

	"chinese-chess-ai/internal/chess"
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/game"

	"github.com/gin-gonic/gin"
)

var gameManager *game.Manager

func init() {
	cfg := config.LoadConfig()
	gameManager = game.GetManager(cfg)
}

// NewGameRequest 新游戏请求
type NewGameRequest struct {
	PlayerColor int `json:"playerColor"` // 1=红方, 2=黑方
}

// NewGame 创建新游戏
func NewGame(c *gin.Context) {
	var req NewGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	playerColor := chess.Red
	if req.PlayerColor == 2 {
		playerColor = chess.Black
	}

	g, err := gameManager.NewGame(playerColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    g.GetState(),
	})
}

// MoveRequest 走子请求
type MoveRequest struct {
	GameID  string `json:"gameId"`
	FromRow int    `json:"fromRow"`
	FromCol int    `json:"fromCol"`
	ToRow   int    `json:"toRow"`
	ToCol   int    `json:"toCol"`
}

// PlayerMove 玩家走子
func PlayerMove(c *gin.Context) {
	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	g, err := gameManager.GetGame(req.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := g.PlayerMove(req.FromRow, req.FromCol, req.ToRow, req.ToCol); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    g.GetState(),
	})
}

// GetGameState 获取游戏状态
func GetGameState(c *gin.Context) {
	gameID := c.Param("id")

	g, err := gameManager.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    g.GetState(),
	})
}

// AIMove AI走子
func AIMove(c *gin.Context) {
	gameID := c.Param("id")

	g, err := gameManager.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	answer, err := g.AIMove()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"state":  g.GetState(),
			"answer": answer,
		},
	})
}
