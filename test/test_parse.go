package main

import (
	"fmt"
	"regexp"
	"strings"
)

func extractMove(answer string) string {
	// 匹配 MOVE: 数字格式
	movePattern := regexp.MustCompile(`MOVE:\s*(\d)(\d)-(\d)(\d)`)
	matches := movePattern.FindStringSubmatch(answer)
	if len(matches) >= 5 {
		return fmt.Sprintf("%s%s-%s%s", matches[1], matches[2], matches[3], matches[4])
	}
	return ""
}

func parseMove(moveStr string) (from, to int, err error) {
	parts := strings.Split(moveStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("移动格式错误: %s", moveStr)
	}

	// 解析起始位置
	if len(parts[0]) != 2 {
		return 0, 0, fmt.Errorf("起始位置格式错误: %s", parts[0])
	}
	fromRow := int(parts[0][0] - '0')
	fromCol := int(parts[0][1] - '0')
	from = fromRow*10 + fromCol

	// 解析目标位置
	if len(parts[1]) != 2 {
		return 0, 0, fmt.Errorf("目标位置格式错误: %s", parts[1])
	}
	toRow := int(parts[1][0] - '0')
	toCol := int(parts[1][1] - '0')
	to = toRow*10 + toCol

	return from, to, nil
}

func main() {
	testCases := []string{
		"MOVE: 02-24",
		"MOVE: 01-21",
		"MOVE: 71-51",
		"分析当前局面：\n- 黑方刚走炮8平5\n\nMOVE: 02-24\n\n这是最佳走法",
	}

	for _, tc := range testCases {
		fmt.Printf("\n测试用例: %s\n", tc)
		moveStr := extractMove(tc)
		fmt.Printf("提取到的走子字符串: %s\n", moveStr)
		
		if moveStr != "" {
			from, to, err := parseMove(moveStr)
			if err != nil {
				fmt.Printf("解析错误: %v\n", err)
			} else {
				fromRow, fromCol := from/10, from%10
				toRow, toCol := to/10, to%10
				fmt.Printf("解析结果: from=%d, to=%d\n", from, to)
				fmt.Printf("坐标: 从(%d,%d)到(%d,%d)\n", fromRow, fromCol, toRow, toCol)
			}
		}
	}
}
