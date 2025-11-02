package main

import (
	"chinese-chess-ai/internal/config"
	"chinese-chess-ai/internal/lke"
	"fmt"
)

func main() {
	cfg := &config.TencentCloudConfig{
		AppID: "test",
	}
	client, _ := lke.NewClient(cfg)

	// 测试Chat方法是否返回空字符串（因为没有真实的API调用）
	result, err := client.Chat("test-session", "test")
	fmt.Printf("Chat result: %q, error: %v\n", result, err)

	// 测试ChatWithDetails方法
	details, err := client.ChatWithDetails("test-session", "test")
	if err != nil {
		fmt.Printf("ChatWithDetails error: %v\n", err)
	} else {
		fmt.Printf("ThoughtProcess: %q\n", details.ThoughtProcess)
		fmt.Printf("FullResponse: %q\n", details.FullResponse)
	}
}
