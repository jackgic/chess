package main

import (
	"fmt"
	"os"

	"chinese-chess-ai/internal/chess"
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/lke"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("🔍 LKE通信调试工具")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 加载配置
	fmt.Println("[1/5] 加载配置...")
	cfg := config.LoadConfig()
	
	fmt.Printf("  - LKE_APP_ID: %s\n", maskString(cfg.TencentCloud.AppID))
	fmt.Printf("  - LKE_APP_KEY: %s\n", maskString(cfg.TencentCloud.BotAppKey))
	fmt.Printf("  - LKE_BOT_BIZ_ID: %s\n", cfg.TencentCloud.BotBizID)
	fmt.Printf("  - LKE_ENDPOINT: %s\n", cfg.TencentCloud.Endpoint)
	fmt.Println()

	// 检查配置
	if cfg.TencentCloud.AppID == "" {
		fmt.Println("❌ 错误: LKE_APP_ID 未配置")
		fmt.Println("请设置环境变量: export LKE_APP_ID=your_app_id")
		os.Exit(1)
	}

	if cfg.TencentCloud.BotAppKey == "" {
		fmt.Println("❌ 错误: LKE_APP_KEY 未配置")
		fmt.Println("请设置环境变量: export LKE_APP_KEY=your_app_key")
		os.Exit(1)
	}

	// 2. 创建LKE客户端
	fmt.Println("[2/5] 创建LKE客户端...")
	client, err := lke.NewClient(&cfg.TencentCloud)
	if err != nil {
		fmt.Printf("❌ 创建客户端失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 客户端创建成功")
	fmt.Println()

	// 3. 创建测试棋盘
	fmt.Println("[3/5] 创建测试棋盘...")
	board := chess.NewBoard()
	
	// 模拟红方走了一步马
	from := chess.Position{Row: 9, Col: 1}
	to := chess.Position{Row: 7, Col: 2}
	if err := board.Move(from.ToInt(), to.ToInt()); err != nil {
		fmt.Printf("❌ 走子失败: %v\n", err)
		os.Exit(1)
	}
	board.Turn = chess.Black // 切换到黑方
	
	fmt.Println("✅ 棋盘创建成功")
	fmt.Printf("  - 当前回合: %s\n", colorToString(board.Turn))
	fmt.Printf("  - 历史走子: %v\n", board.MoveList)
	fmt.Println()

	// 4. 构建提示词
	fmt.Println("[4/5] 构建AI提示词...")
	prompt := fmt.Sprintf(`当前棋盘状态（FEN格式）：%s

当前轮到：%s
你是：黑方

历史走子：
1. 2退3

请分析当前局面并走子。

【重要】走子格式说明：
1. 必须使用坐标格式：MOVE: 起始行起始列-目标行目标列
2. 坐标范围：行(0-9)，列(0-8)
3. 红方在下方(行7-9)，黑方在上方(行0-2)
4. 格式示例：MOVE: XY-ZW，其中X、Y、Z、W都是0-9的数字

请立即给出你的走子（必须严格按照MOVE: XY-ZW格式）：`,
		board.ToFEN(),
		colorToString(board.Turn))

	fmt.Println("提示词内容:")
	fmt.Println("---")
	fmt.Println(prompt)
	fmt.Println("---")
	fmt.Println()

	// 5. 调用LKE
	fmt.Println("[5/5] 调用LKE智能体...")
	fmt.Println("⏳ 等待AI响应（可能需要10-30秒）...")
	fmt.Println()

	sessionID := "test_session_123"
	answer, err := client.Chat(sessionID, prompt)
	if err != nil {
		fmt.Printf("❌ LKE调用失败: %v\n", err)
		fmt.Println()
		fmt.Println("可能的原因:")
		fmt.Println("  1. 网络连接问题")
		fmt.Println("  2. LKE_APP_ID 或 LKE_APP_KEY 配置错误")
		fmt.Println("  3. LKE服务异常")
		fmt.Println()
		fmt.Println("调试建议:")
		fmt.Println("  - 检查网络: curl -I https://wss.lke.cloud.tencent.com")
		fmt.Println("  - 验证配置: echo $LKE_APP_ID")
		fmt.Println("  - 查看详细日志")
		os.Exit(1)
	}

	fmt.Println("✅ AI响应成功!")
	fmt.Println()
	fmt.Println("AI完整回答:")
	fmt.Println("---")
	fmt.Println(answer)
	fmt.Println("---")
	fmt.Println()

	// 6. 解析走子
	fmt.Println("[6/6] 解析AI走子...")
	move, err := client.ExtractMove(answer)
	if err != nil {
		fmt.Printf("❌ 解析走子失败: %v\n", err)
		fmt.Println()
		fmt.Println("可能的原因:")
		fmt.Println("  1. AI没有按照MOVE格式输出")
		fmt.Println("  2. 智能体提示词配置不正确")
		fmt.Println()
		fmt.Println("解决方案:")
		fmt.Println("  - 检查智能体的系统提示词")
		fmt.Println("  - 确保提示词中强调了MOVE格式")
		os.Exit(1)
	}

	fromRow, fromCol := move.From/10, move.From%10
	toRow, toCol := move.To/10, move.To%10

	fmt.Println("✅ 走子解析成功!")
	fmt.Printf("  - 起始位置: (%d, %d)\n", fromRow, fromCol)
	fmt.Printf("  - 目标位置: (%d, %d)\n", toRow, toCol)
	fmt.Printf("  - 走子编码: from=%d, to=%d\n", move.From, move.To)
	fmt.Println()

	// 7. 验证走子
	fmt.Println("[7/7] 验证走子合法性...")
	fromPos := chess.Position{Row: fromRow, Col: fromCol}
	toPos := chess.Position{Row: toRow, Col: toCol}

	piece := board.Grid[fromRow][fromCol]
	if piece.Type == chess.Empty {
		fmt.Printf("❌ 起始位置(%d,%d)没有棋子\n", fromRow, fromCol)
		os.Exit(1)
	}

	fmt.Printf("  - 起始位置棋子: %s\n", board.PieceToString(piece))

	if !board.IsValidMove(fromPos, toPos) {
		fmt.Printf("❌ 走子不合法: 从(%d,%d)到(%d,%d)\n", fromRow, fromCol, toRow, toCol)
		os.Exit(1)
	}

	fmt.Println("✅ 走子合法!")
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("🎉 测试完成！LKE通信正常！")
	fmt.Println("========================================")
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

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
