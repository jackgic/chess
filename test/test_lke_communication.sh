#!/bin/bash

echo "=========================================="
echo "🧪 中国象棋AI - LKE通信测试"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试服务器是否运行
echo -e "${BLUE}[1/5] 检查服务器状态...${NC}"
if curl -s http://localhost:8090 > /dev/null; then
    echo -e "${GREEN}✅ 服务器运行正常${NC}"
else
    echo -e "${RED}❌ 服务器未运行，请先启动服务器${NC}"
    exit 1
fi
echo ""

# 测试创建游戏
echo -e "${BLUE}[2/5] 创建新游戏...${NC}"
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
    -H "Content-Type: application/json" \
    -d '{"playerColor":1}')
GAME_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$GAME_ID" ]; then
    echo -e "${RED}❌ 创建游戏失败${NC}"
    echo "响应: $CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ 游戏创建成功${NC}"
echo "   游戏ID: $GAME_ID"
echo ""

# 测试玩家走子（红方先走）
echo -e "${BLUE}[3/5] 测试玩家走子（红方炮二平五）...${NC}"
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
    -H "Content-Type: application/json" \
    -d "{\"gameId\":\"$GAME_ID\",\"fromRow\":7,\"fromCol\":1,\"toRow\":7,\"toCol\":4}")

if echo "$MOVE_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 玩家走子成功${NC}"
    echo "   响应: $MOVE_RESPONSE"
else
    echo -e "${RED}❌ 玩家走子失败${NC}"
    echo "   响应: $MOVE_RESPONSE"
    exit 1
fi
echo ""

# 测试AI走子（关键测试）
echo -e "${BLUE}[4/5] 测试AI走子（与LKE智能体通信）...${NC}"
echo -e "${YELLOW}⏳ 正在等待AI响应（可能需要5-10秒）...${NC}"
echo ""

AI_RESPONSE=$(curl -s -X POST "http://localhost:8090/api/game/$GAME_ID/ai-move")

echo "AI完整响应:"
echo "$AI_RESPONSE" | jq '.' 2>/dev/null || echo "$AI_RESPONSE"
echo ""

# 检查AI响应
if echo "$AI_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ AI走子成功 - LKE通信正常！${NC}"
    
    # 提取走子信息
    AI_MOVE=$(echo "$AI_RESPONSE" | grep -o '"move":"[^"]*"' | cut -d'"' -f4)
    AI_MESSAGE=$(echo "$AI_RESPONSE" | grep -o '"message":"[^"]*"' | cut -d'"' -f4)
    
    echo ""
    echo "📊 AI走子详情:"
    echo "   走子: $AI_MOVE"
    echo "   说明: $AI_MESSAGE"
    
elif echo "$AI_RESPONSE" | grep -q "模拟AI"; then
    echo -e "${YELLOW}⚠️  当前使用模拟AI（未配置真实LKE）${NC}"
    echo ""
    echo "💡 如需使用真实LKE智能体，请："
    echo "   1. 确认 .env 文件中的 LKE_APP_KEY 已正确配置"
    echo "   2. 检查网络连接"
    echo "   3. 查看服务器日志了解详细错误"
    
else
    echo -e "${RED}❌ AI走子失败${NC}"
    echo ""
    echo "可能的原因："
    echo "   1. LKE_APP_KEY 配置错误"
    echo "   2. 网络连接问题"
    echo "   3. LKE服务异常"
    echo ""
    echo "建议："
    echo "   - 检查 .env 文件中的配置"
    echo "   - 查看服务器日志: tail -f /tmp/chinese-chess-ai.log"
fi
echo ""

# 查询游戏状态
echo -e "${BLUE}[5/5] 查询游戏状态...${NC}"
STATUS_RESPONSE=$(curl -s "http://localhost:8090/api/game/$GAME_ID")

if echo "$STATUS_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 状态查询成功${NC}"
    echo ""
    echo "📋 当前游戏状态:"
    echo "$STATUS_RESPONSE" | jq '.data' 2>/dev/null || echo "$STATUS_RESPONSE"
else
    echo -e "${RED}❌ 状态查询失败${NC}"
fi
echo ""

# 总结
echo "=========================================="
echo "📊 测试总结"
echo "=========================================="
echo ""

if echo "$AI_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}🎉 所有测试通过！LKE智能体通信正常！${NC}"
    echo ""
    echo "✅ 服务器运行正常"
    echo "✅ 游戏创建成功"
    echo "✅ 玩家走子成功"
    echo "✅ AI走子成功（LKE通信正常）"
    echo "✅ 状态查询成功"
    echo ""
    echo "🎮 你现在可以访问 http://localhost:8090 开始游戏了！"
else
    echo -e "${YELLOW}⚠️  基本功能正常，但LKE通信可能存在问题${NC}"
    echo ""
    echo "✅ 服务器运行正常"
    echo "✅ 游戏创建成功"
    echo "✅ 玩家走子成功"
    echo "⚠️  AI走子使用模拟模式"
    echo ""
    echo "💡 建议："
    echo "   1. 检查 .env 文件中的 LKE_APP_KEY"
    echo "   2. 查看服务器日志获取详细错误信息"
    echo "   3. 参考 ADP_SETUP_GUIDE.md 配置智能体"
fi

echo ""
echo "=========================================="
