#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================="
echo "🎮 中国象棋AI通信调试测试"
echo "==========================================${NC}"
echo ""

# 1. 创建新游戏
echo -e "${BLUE}[1/3] 创建新游戏...${NC}"
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
  -H "Content-Type: application/json" \
  -d '{"playerColor": 1}')

echo "$CREATE_RESPONSE" | jq '.' 2>/dev/null || echo "$CREATE_RESPONSE"

GAME_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id' 2>/dev/null)

if [ -z "$GAME_ID" ] || [ "$GAME_ID" = "null" ]; then
    echo -e "${RED}❌ 创建游戏失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 游戏创建成功，游戏ID: $GAME_ID${NC}"
echo ""

# 2. 玩家走子
echo -e "${BLUE}[2/3] 玩家走子（马二进三）...${NC}"
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d "{\"gameId\": \"$GAME_ID\", \"fromRow\": 9, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 2}")

echo "$MOVE_RESPONSE" | jq '.' 2>/dev/null || echo "$MOVE_RESPONSE"

if echo "$MOVE_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 玩家走子成功${NC}"
else
    echo -e "${RED}❌ 玩家走子失败${NC}"
    exit 1
fi
echo ""

# 3. AI走子（带超时控制）
echo -e "${BLUE}[3/3] 请求AI走子...${NC}"
echo -e "${YELLOW}⏳ 等待AI思考（最多60秒）...${NC}"

# 使用timeout命令，如果60秒内没有响应则超时
AI_RESPONSE=$(timeout 60 curl -s -X POST "http://localhost:8090/api/game/$GAME_ID/ai-move" \
  -H "Content-Type: application/json")

EXIT_CODE=$?

if [ $EXIT_CODE -eq 124 ]; then
    echo -e "${RED}❌ AI走子超时（60秒）${NC}"
    echo ""
    echo "可能的原因："
    echo "  1. LKE服务响应慢"
    echo "  2. 网络连接问题"
    echo "  3. 配置错误"
    echo ""
    echo "建议检查："
    echo "  - 环境变量 LKE_APP_ID 和 LKE_APP_KEY 是否正确"
    echo "  - 网络是否能访问 https://wss.lke.cloud.tencent.com"
    echo "  - 查看服务器日志获取详细错误信息"
    exit 1
fi

echo ""
echo -e "${BLUE}AI响应内容：${NC}"
echo "$AI_RESPONSE" | jq '.' 2>/dev/null || echo "$AI_RESPONSE"
echo ""

if echo "$AI_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ AI走子成功！${NC}"
    echo ""
    
    # 提取AI的分析
    AI_ANSWER=$(echo "$AI_RESPONSE" | jq -r '.data.answer' 2>/dev/null)
    if [ ! -z "$AI_ANSWER" ] && [ "$AI_ANSWER" != "null" ]; then
        echo -e "${BLUE}🤖 AI分析：${NC}"
        echo "$AI_ANSWER"
        echo ""
    fi
    
    # 显示当前棋盘状态
    echo -e "${BLUE}📋 当前游戏状态：${NC}"
    GAME_STATE=$(curl -s "http://localhost:8090/api/game/$GAME_ID")
    echo "$GAME_STATE" | jq '.data | {turn, status, moveList}' 2>/dev/null
    
else
    echo -e "${RED}❌ AI走子失败${NC}"
    echo ""
    
    # 提取错误信息
    ERROR_MSG=$(echo "$AI_RESPONSE" | jq -r '.error' 2>/dev/null)
    if [ ! -z "$ERROR_MSG" ] && [ "$ERROR_MSG" != "null" ]; then
        echo -e "${RED}错误信息：${NC}"
        echo "$ERROR_MSG"
    fi
    
    echo ""
    echo "调试建议："
    echo "  1. 检查 LKE_APP_ID 和 LKE_APP_KEY 配置"
    echo "  2. 确认智能体提示词配置正确"
    echo "  3. 查看服务器日志中的详细错误"
    echo "  4. 测试网络连接：curl -I https://wss.lke.cloud.tencent.com"
    exit 1
fi

echo ""
echo -e "${GREEN}=========================================="
echo "🎉 测试完成！"
echo "==========================================${NC}"
