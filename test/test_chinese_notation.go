package main

import (
	"fmt"
	"os"

	"chinese-chess-ai/internal/chess"
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/lke"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== 中文棋谱交互测试 ===")
	fmt.Println()

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Printf("警告: 无法加载.env文件: %v\n", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 创建LKE客户端
	lkeClient, err := lke.NewClient(&cfg.TencentCloud)
	if err != nil {
		fmt.Printf("❌ 创建LKE客户端失败: %v\n", err)
		os.Exit(1)
	}

	// 创建棋盘
	board := chess.NewBoard()
	sessionID := "test_session_123"

	fmt.Println("✅ 初始化成功")
	fmt.Println()

	// 测试1: 玩家走子（红方）
	fmt.Println("【测试1】玩家走子：炮二平五")
	if err := board.Move(71, 74); err != nil {
		fmt.Printf("❌ 走子失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 走子成功，棋谱记录: %s\n", board.MoveList[len(board.MoveList)-1])
	fmt.Println()

	// 测试2: 发送给AI（只发送中文棋谱）
	fmt.Println("【测试2】发送给AI的消息")
	lastMove := board.MoveList[len(board.MoveList)-1]
	fmt.Printf("发送内容: %s\n", lastMove)
	fmt.Println()

	// 测试3: 调用AI
	fmt.Println("【测试3】调用AI获取回复")
	answer, err := lkeClient.Chat(sessionID, lastMove)
	if err != nil {
		fmt.Printf("❌ AI调用失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("AI回复: %s\n", answer)
	fmt.Println()

	// 测试4: 解析AI返回的中文棋谱
	fmt.Println("【测试4】解析AI返回的中文棋谱")
	chineseMove, err := lkeClient.ExtractChineseMove(answer)
	if err != nil {
		fmt.Printf("❌ 解析失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 解析成功: %s\n", chineseMove)
	fmt.Println()

	// 测试5: 执行AI的走子
	fmt.Println("【测试5】执行AI的走子")
	if err := board.MoveByChineseNotation(chineseMove, chess.Black); err != nil {
		fmt.Printf("❌ 执行失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 执行成功，棋谱记录: %s\n", board.MoveList[len(board.MoveList)-1])
	fmt.Println()

	// 显示当前棋盘状态
	fmt.Println("【当前棋盘状态】")
	fmt.Println(board.ToString())
	fmt.Println()

	// 显示完整棋谱
	fmt.Println("【完整棋谱】")
	for i, move := range board.MoveList {
		fmt.Printf("%d. %s\n", i+1, move)
	}
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
}
