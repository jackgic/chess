package main

import (
	"fmt"
	"log"
	"os"

	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/lke"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
	}

	// 加载配置
	cfg := config.LoadConfig()

	fmt.Println("========================================")
	fmt.Println("🧪 LKE智能体通信测试")
	fmt.Println("========================================")
	fmt.Println()

	// 显示配置信息
	fmt.Println("📋 当前配置:")
	fmt.Printf("   SecretID: %s\n", maskString(cfg.TencentCloud.SecretID))
	fmt.Printf("   SecretKey: %s\n", maskString(cfg.TencentCloud.SecretKey))
	fmt.Printf("   AppKey: %s\n", maskString(cfg.TencentCloud.AppKey))
	fmt.Printf("   Region: %s\n", cfg.TencentCloud.Region)
	fmt.Println()

	// 检查配置
	if cfg.TencentCloud.AppKey == "" || cfg.TencentCloud.AppKey == "your_app_key_here" {
		fmt.Println("❌ 错误: LKE_APP_KEY 未配置或使用默认值")
		fmt.Println()
		fmt.Println("请在 .env 文件中配置正确的 LKE_APP_KEY")
		fmt.Println("获取方式：")
		fmt.Println("1. 访问 https://adp.cloud.tencent.com/")
		fmt.Println("2. 选择你的智能体应用")
		fmt.Println("3. 进入 '应用配置' -> 'API密钥'")
		fmt.Println("4. 复制 AppKey")
		os.Exit(1)
	}

	// 创建LKE客户端
	fmt.Println("🔌 创建LKE客户端...")
	client, err := lke.NewClient(&cfg.TencentCloud)
	if err != nil {
		fmt.Printf("❌ 客户端创建失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 客户端创建成功")
	fmt.Println()

	// 测试简单对话
	fmt.Println("💬 测试1: 简单对话")
	fmt.Println("   问题: 你好，请介绍一下你自己")
	fmt.Println()

	response1, err := client.Chat(lke.ChatRequest{
		Query:      "你好，请介绍一下你自己",
		BoardState: "",
	})
	if err != nil {
		fmt.Printf("❌ 对话失败: %v\n", err)
		fmt.Println()
		fmt.Println("可能的原因：")
		fmt.Println("1. AppKey 配置错误")
		fmt.Println("2. 网络连接问题")
		fmt.Println("3. LKE服务异常")
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("✅ 对话成功！")
	fmt.Printf("   AI回复: %s\n", response1.Answer)
	fmt.Println()

	// 测试象棋走子
	fmt.Println("♟️  测试2: 象棋走子指令")
	fmt.Println("   问题: 当前是开局，请走第一步棋")
	fmt.Println()

	chessPrompt := `你是中国象棋AI助手。

当前棋盘状态（FEN格式）：
rnbakabnr/9/1c5c1/p1p1p1p1p/9/9/P1P1P1P1P/1C5C1/9/RNBAKABNR w - - 0 1

这是开局状态，现在轮到红方（你）走棋。

请分析局面并给出你的走子。

输出格式（必须严格遵守）：
MOVE: 起始行起始列-目标行目标列

示例：
MOVE: 71-74  （将炮从(7,1)移到(7,4)）

请给出你的走子：`

	response2, err := client.Chat(lke.ChatRequest{
		Query:      chessPrompt,
		BoardState: "rnbakabnr/9/1c5c1/p1p1p1p1p/9/9/P1P1P1P1P/1C5C1/9/RNBAKABNR w - - 0 1",
	})
	if err != nil {
		fmt.Printf("❌ 对话失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 对话成功！")
	fmt.Printf("   AI回复: %s\n", response2.Answer)
	fmt.Println()

	// 尝试提取走子指令
	fmt.Println("🔍 尝试提取走子指令...")
	move := response2.Move
	if move == "" {
		move = extractMove(response2.Answer)
	}
	if move != "" {
		fmt.Printf("✅ 成功提取走子: %s\n", move)
	} else {
		fmt.Println("⚠️  未能提取到走子指令")
		fmt.Println("   这可能意味着AI的回复格式不符合预期")
		fmt.Println("   需要优化系统提示词")
	}
	fmt.Println()

	// 总结
	fmt.Println("========================================")
	fmt.Println("📊 测试总结")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("✅ LKE客户端创建成功")
	fmt.Println("✅ 简单对话测试通过")
	fmt.Println("✅ 象棋对话测试通过")
	if move != "" {
		fmt.Println("✅ 走子指令提取成功")
		fmt.Println()
		fmt.Println("🎉 所有测试通过！LKE智能体通信正常！")
	} else {
		fmt.Println("⚠️  走子指令提取失败")
		fmt.Println()
		fmt.Println("💡 建议：")
		fmt.Println("   1. 在ADP平台配置系统提示词")
		fmt.Println("   2. 参考 PROMPT_TEMPLATE.md 文件")
		fmt.Println("   3. 确保AI理解输出格式要求")
	}
	fmt.Println()
	fmt.Println("========================================")
}

// maskString 遮蔽字符串中间部分
func maskString(s string) string {
	if s == "" || s == "your_secret_id_here" || s == "your_secret_key_here" || s == "your_app_key_here" {
		return "[未配置]"
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// extractMove 提取走子指令
func extractMove(response string) string {
	// 查找类似 "71-74" 的模式
	for i := 0; i < len(response)-4; i++ {
		if response[i] >= '0' && response[i] <= '9' &&
			response[i+1] >= '0' && response[i+1] <= '9' &&
			response[i+2] == '-' &&
			response[i+3] >= '0' && response[i+3] <= '9' &&
			i+4 < len(response) && response[i+4] >= '0' && response[i+4] <= '9' {
			return response[i : i+5]
		}
	}

	return ""
}
