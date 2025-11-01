package lke

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"chinese-chess-ai/internal/config"
)

// Client LKE客户端
type Client struct {
	config *config.TencentCloudConfig
	sseURL string
}

// NewClient 创建LKE客户端
func NewClient(cfg *config.TencentCloudConfig) (*Client, error) {
	return &Client{
		config: cfg,
		sseURL: "https://wss.lke.cloud.tencent.com/v1/qbot/chat/sse",
	}, nil
}

// StartSession 通知智能体开始新对局
func (c *Client) StartSession(sessionID string) error {
	// 对于SSE模式，不需要单独的开始会话调用
	// 会话在第一次Chat调用时自动开始
	log.Printf("[LKE] 会话 %s 已准备就绪", sessionID)
	return nil
}

// ChatRequest 对话请求
type ChatRequest struct {
	Query      string   // 用户问题
	History    []string // 历史对话
	BoardState string   // 当前棋盘状态
}

// ChatResponse 对话响应
type ChatResponse struct {
	Answer string // AI回答
	Move   string // 走子指令（格式：0102-0304）
}

// SSERequest SSE请求参数
type SSERequest struct {
	Content         string                 `json:"content"`
	BotAppKey       string                 `json:"bot_app_key"`
	SessionID       string                 `json:"session_id"`
	VisitorBizID    string                 `json:"visitor_biz_id"`
	SearchNetwork   string                 `json:"search_network"`
	CustomVariables map[string]interface{} `json:"custom_variables"`
	Incremental     bool                   `json:"incremental"`
}

// SSEMessage SSE消息
type SSEMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ReplyPayload 回复消息载荷
type ReplyPayload struct {
	Content string `json:"content"`
}

// ThoughtPayload 思考过程载荷
type ThoughtPayload struct {
	Thought string `json:"thought"`
}

// ErrorPayload 错误消息载荷
type ErrorPayload struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Chat 使用会话ID与智能体交互
func (c *Client) Chat(sessionID, prompt string) (string, error) {
	// 构建请求体
	reqBody := map[string]interface{}{
		"session_id":     sessionID,
		"content":        prompt,
		"bot_app_key":    c.config.AppID,
		"visitor_biz_id": "player_1",
		"search_network": "disable",
		"incremental":    true,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %v", err)
	}

	log.Printf("[LKE] 请求体: %s", string(reqBodyJSON))

	// 创建HTTP请求
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.sseURL, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("[LKE] 发送SSE请求到: %s", c.sseURL)

	// 发送请求
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("SSE error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	log.Printf("[LKE] SSE连接成功，开始接收流式响应")

	// 读取SSE流
	var answer strings.Builder
	var thought strings.Builder
	var lastError string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		// 解析data字段
		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimSpace(data)

		log.Printf("[LKE] 原始SSE数据: %s", data)

		var msg SSEMessage
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			log.Printf("[LKE] 解析SSE消息失败: %v, data: %s", err, data)
			continue
		}

		log.Printf("[LKE] 收到消息类型: %s", msg.Type)

		switch msg.Type {
		case "reply":
			var payload ReplyPayload
			if err := json.Unmarshal(msg.Payload, &payload); err == nil {
				answer.WriteString(payload.Content)
				log.Printf("[LKE] 回复内容: %s", payload.Content)
			}
		case "thought":
			var payload ThoughtPayload
			if err := json.Unmarshal(msg.Payload, &payload); err == nil {
				thought.WriteString(payload.Thought)
				log.Printf("[LKE] 思考过程: %s", payload.Thought)
			}
		case "error":
			if msg.Error != nil {
				lastError = msg.Error.Message
				log.Printf("[LKE] 错误 (code=%d): %s", msg.Error.Code, lastError)
			} else {
				log.Printf("[LKE] 收到error类型消息但无错误详情")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read SSE stream failed: %v", err)
	}

	if lastError != "" {
		return "", fmt.Errorf("LKE error: %s", lastError)
	}

	return answer.String(), nil
}

// buildPrompt 构建完整提示
func (c *Client) buildPrompt(req ChatRequest) string {
	var prompt strings.Builder

	// 添加棋盘状态
	prompt.WriteString("当前棋盘状态（FEN格式）：\n")
	prompt.WriteString(req.BoardState)
	prompt.WriteString("\n\n")

	// 添加历史对话
	if len(req.History) > 0 {
		prompt.WriteString("历史走子：\n")
		for _, move := range req.History {
			prompt.WriteString(move)
			prompt.WriteString("\n")
		}
		prompt.WriteString("\n")
	}

	// 添加用户问题
	prompt.WriteString(req.Query)

	return prompt.String()
}

// extractMove 从AI回答中提取走子指令
func (c *Client) extractMove(answer string) string {
	// 优先匹配中文象棋记谱格式 MOVE: 炮2平5
	chineseMovePattern := regexp.MustCompile(`MOVE:\s*([炮车马相象士将帅兵卒])(\d|一|二|三|四|五|六|七|八|九)([进退平])(\d|一|二|三|四|五|六|七|八|九)`)
	matches := chineseMovePattern.FindStringSubmatch(answer)
	if len(matches) >= 5 {
		// 转换中文记谱为坐标
		return c.convertChineseMoveToCoords(matches[1], matches[2], matches[3], matches[4])
	}

	// 匹配 MOVE: 数字格式
	movePattern := regexp.MustCompile(`MOVE:\s*(\d)(\d)-(\d)(\d)`)
	matches = movePattern.FindStringSubmatch(answer)
	if len(matches) >= 5 {
		return fmt.Sprintf("%s%s-%s%s", matches[1], matches[2], matches[3], matches[4])
	}

	// 匹配其他格式
	patterns := []string{
		`(\d)(\d)-(\d)(\d)`,           // 0102-0304
		`\((\d),(\d)\)-\((\d),(\d)\)`, // (0,1)-(0,3)
		`从(\d)(\d)到(\d)(\d)`,          // 从01到03
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(answer)
		if len(matches) >= 5 {
			return fmt.Sprintf("%s%s-%s%s", matches[1], matches[2], matches[3], matches[4])
		}
	}

	// 尝试提取JSON格式
	var moveData struct {
		From struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"from"`
		To struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"to"`
	}

	if err := json.Unmarshal([]byte(answer), &moveData); err == nil {
		return fmt.Sprintf("%d%d-%d%d",
			moveData.From.Row, moveData.From.Col,
			moveData.To.Row, moveData.To.Col)
	}

	return ""
}

// Move 移动结构体
type Move struct {
	From int
	To   int
}

// ExtractMove 从AI回答中提取移动指令，返回Move结构体
func (c *Client) ExtractMove(answer string) (*Move, error) {
	moveStr := c.extractMove(answer)
	if moveStr == "" {
		return nil, fmt.Errorf("无法从回答中提取移动指令")
	}

	// 解析移动字符串 "0102-0304" 格式
	parts := strings.Split(moveStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("移动格式错误: %s", moveStr)
	}

	// 解析起始位置
	if len(parts[0]) != 2 {
		return nil, fmt.Errorf("起始位置格式错误: %s", parts[0])
	}
	fromRow := int(parts[0][0] - '0')
	fromCol := int(parts[0][1] - '0')
	from := fromRow*10 + fromCol

	// 解析目标位置
	if len(parts[1]) != 2 {
		return nil, fmt.Errorf("目标位置格式错误: %s", parts[1])
	}
	toRow := int(parts[1][0] - '0')
	toCol := int(parts[1][1] - '0')
	to := toRow*10 + toCol

	return &Move{
		From: from,
		To:   to,
	}, nil
}

// GetSystemPrompt 获取系统提示词模板
func GetSystemPrompt() string {
	return `# 中国象棋AI助手

你是一个专业的中国象棋AI对弈助手。你的任务是根据当前棋盘状态，选择最佳走子并以规定格式输出。

## 角色定位
- 你扮演黑方（上方），与红方（下方）对弈
- 你需要遵守中国象棋的所有规则
- 你应该具备基本的对局策略和战术思维

## 棋盘表示
- 棋盘使用FEN格式表示（类似国际象棋）
- 坐标系统：行号0-9（从上到下），列号0-8（从左到右）
- 棋子符号：
  * 大写字母代表红方：K(帅) A(仕) E(相) H(马) R(车) C(炮) P(兵)
  * 小写字母代表黑方：k(将) a(士) e(象) h(马) r(车) c(炮) p(卒)
  * 数字代表连续的空格数

## 走子规则
1. **将/帅(K/k)**：只能在九宫格内移动，每次一格，不能离开九宫格
2. **士(A/a)**：只能在九宫格内斜走，每次一格
3. **象/相(E/e)**：走田字，不能过河，象眼不能被堵
4. **马(H/h)**：走日字，马腿不能被堵
5. **车(R/r)**：直线移动，路径不能有障碍
6. **炮(C/c)**：直线移动，吃子时必须隔一个棋子
7. **兵/卒(P/p)**：未过河只能向前，过河后可以左右移动

## 对局策略
1. **开局阶段**：
   - 优先出动车、马等强子
   - 控制中路，占据要道
   - 保护己方将帅安全

2. **中局阶段**：
   - 寻找战术机会（如双车错、马后炮等）
   - 兑子时要考虑局面优劣
   - 注意子力协调配合

3. **残局阶段**：
   - 利用子力优势进攻
   - 单车、马炮等残局定式
   - 注意将帅对面（白脸将）

## 输出格式要求
**极其重要**：你必须严格按照以下格式输出走子指令：

格式：MOVE: 起始行起始列-目标行目标列

示例：
- MOVE: 01-03 （将位置(0,1)的棋子移动到(0,3)）
- MOVE: 27-47 （将位置(2,7)的炮移动到(4,7)）

### 完整回答格式
[简要分析当前局面]

我选择的走子是：
MOVE: 起始行起始列-目标行目标列

[简要说明走子理由]

### 示例回答
当前局面红方车马已出，我方需要加强中路控制。

我选择的走子是：
MOVE: 27-47

这一步将炮移到中路，既可以控制中心，又为后续进攻做准备。

## 注意事项
1. **必须输出合法走子**：每次回答都必须包含一个合法的MOVE指令
2. **坐标格式严格**：必须是4位数字，前两位是起始位置，后两位是目标位置
3. **检查合法性**：确保走子符合该棋子的移动规则
4. **避免送子**：不要走出明显的送子步骤
5. **考虑对方威胁**：注意对方的进攻意图，做好防守

## 错误示例（禁止）
❌ "我觉得应该走马" （没有具体坐标）
❌ "MOVE: a2-a4" （使用了字母坐标）
❌ "移动炮到中路" （没有MOVE格式）
❌ "MOVE: 0,1-0,3" （使用了逗号分隔）

## 正确示例（推荐）
✅ "MOVE: 01-21"
✅ "MOVE: 77-75"
✅ "MOVE: 00-02"

记住：每次回答都必须包含一个明确的MOVE指令，格式为"MOVE: 数字数字-数字数字"。
`
}

// convertChineseMoveToCoords 将中文象棋记谱转换为坐标格式
func (c *Client) convertChineseMoveToCoords(piece, from, action, to string) string {
	// 数字转换映射
	numMap := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 7, "八": 8, "九": 9,
	}

	// 获取起始列（从字符串转换为数字）
	fromCol, exists := numMap[from]
	if !exists {
		return ""
	}

	// 获取目标列或行（从字符串转换为数字）
	toNum, exists := numMap[to]
	if !exists {
		return ""
	}

	// 简化处理：假设是黑方炮的移动
	// 炮2平5 表示从第2列平移到第5列
	if action == "平" {
		// 平移：行不变，列改变
		// 假设炮在第2行（黑方炮的初始位置）
		fromRow := 2
		toRow := fromRow
		toCol := toNum

		// 转换为坐标格式：行列都从0开始
		return fmt.Sprintf("%d%d-%d%d", fromRow, fromCol-1, toRow, toCol-1)
	} else if action == "进" {
		// 前进：行减少（向红方移动）
		fromRow := 2 // 假设当前在第2行
		toRow := fromRow + toNum
		toCol := fromCol

		return fmt.Sprintf("%d%d-%d%d", fromRow, fromCol-1, toRow, toCol-1)
	} else if action == "退" {
		// 后退：行增加（向黑方移动）
		fromRow := 2 // 假设当前在第2行
		toRow := fromRow - toNum
		toCol := fromCol

		return fmt.Sprintf("%d%d-%d%d", fromRow, fromCol-1, toRow, toCol-1)
	}

	return ""
}
