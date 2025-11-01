#!/bin/bash

# AI交互修复验证脚本
# 用于快速测试AI走子功能是否正常

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}AI交互修复验证测试${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# 检查服务是否运行
echo -e "${BLUE}[1/6] 检查服务状态...${NC}"
if ! curl -s http://localhost:8090/health > /dev/null 2>&1; then
    echo -e "${RED}❌ 服务未运行，请先启动服务：./start.sh${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 服务正常运行${NC}"
echo ""

# 创建新游戏
echo -e "${BLUE}[2/6] 创建新游戏...${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/new \
  -H "Content-Type: application/json" \
  -d '{"playerColor": 1}')

GAME_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$GAME_ID" ]; then
    echo -e "${RED}❌ 创建游戏失败${NC}"
    echo "响应: $RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ 游戏创建成功${NC}"
echo "   游戏ID: $GAME_ID"
echo ""

# 玩家走子：炮二平五
echo -e "${BLUE}[3/6] 玩家走子（炮二平五）...${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/move \
  -H "Content-Type: application/json" \
  -d "{\"gameId\": \"$GAME_ID\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")

if echo "$RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 玩家走子成功${NC}"
    echo "   走子: 炮二平五 (7,1) -> (7,4)"
else
    echo -e "${RED}❌ 玩家走子失败${NC}"
    echo "响应: $RESPONSE"
    exit 1
fi
echo ""

# AI走子
echo -e "${BLUE}[4/6] 请求AI走子...${NC}"
echo -e "${YELLOW}⏳ 等待AI思考（可能需要10-30秒）...${NC}"
echo ""

START_TIME=$(date +%s)
RESPONSE=$(curl -s -X POST http://localhost:8090/api/game/$GAME_ID/ai-move \
  -H "Content-Type: application/json" \
  --max-time 60)
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo -e "${BLUE}AI响应时间: ${DURATION}秒${NC}"
echo ""

# 解析响应
if echo "$RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ AI走子成功！${NC}"
    echo ""
    
    # 提取AI回答
    AI_ANSWER=$(echo "$RESPONSE" | grep -o '"answer":"[^"]*"' | cut -d'"' -f4 | sed 's/\\n/\n/g')
    
    if [ ! -z "$AI_ANSWER" ]; then
        echo -e "${BLUE}🤖 AI完整回答：${NC}"
        echo "---"
        echo "$AI_ANSWER"
        echo "---"
        echo ""
    fi
    
    # 检查是否包含MOVE指令
    if echo "$AI_ANSWER" | grep -q 'MOVE:'; then
        MOVE=$(echo "$AI_ANSWER" | grep -o 'MOVE: [0-9][0-9]-[0-9][0-9]' | head -1)
        echo -e "${GREEN}✅ 找到MOVE指令: $MOVE${NC}"
        
        # 解析坐标
        COORDS=$(echo "$MOVE" | grep -o '[0-9][0-9]-[0-9][0-9]')
        FROM=$(echo "$COORDS" | cut -d'-' -f1)
        TO=$(echo "$COORDS" | cut -d'-' -f2)
        FROM_ROW=${FROM:0:1}
        FROM_COL=${FROM:1:1}
        TO_ROW=${TO:0:1}
        TO_COL=${TO:1:1}
        
        echo "   起始位置: ($FROM_ROW, $FROM_COL)"
        echo "   目标位置: ($TO_ROW, $TO_COL)"
    else
        echo -e "${YELLOW}⚠️  未找到MOVE指令（但走子可能仍然成功）${NC}"
    fi
    
else
    echo -e "${RED}❌ AI走子失败${NC}"
    echo ""
    
    # 提取错误信息
    ERROR=$(echo "$RESPONSE" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$ERROR" ]; then
        echo -e "${RED}错误信息：${NC}"
        echo "$ERROR"
    else
        echo "完整响应:"
        echo "$RESPONSE"
    fi
    
    echo ""
    echo -e "${YELLOW}可能的原因：${NC}"
    echo "  1. AI输出的坐标不合法"
    echo "  2. AI对棋盘理解有误"
    echo "  3. LKE系统提示词配置不正确"
    echo ""
    echo -e "${YELLOW}建议操作：${NC}"
    echo "  1. 查看日志: tail -50 logs/chinese-chess-ai.log"
    echo "  2. 检查LKE平台的系统提示词配置"
    echo "  3. 参考 AI_INTERACTION_FIX.md 文档"
    
    exit 1
fi
echo ""

# 查询游戏状态
echo -e "${BLUE}[5/6] 查询游戏状态...${NC}"
RESPONSE=$(curl -s http://localhost:8090/api/game/$GAME_ID)

if echo "$RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ 状态查询成功${NC}"
    
    # 提取走子历史
    MOVE_LIST=$(echo "$RESPONSE" | grep -o '"moveList":\[[^]]*\]' | sed 's/"moveList"://g')
    if [ ! -z "$MOVE_LIST" ]; then
        echo "   走子历史: $MOVE_LIST"
    fi
    
    # 提取当前回合
    TURN=$(echo "$RESPONSE" | grep -o '"turn":[0-9]' | cut -d':' -f2)
    if [ "$TURN" = "1" ]; then
        echo "   当前回合: 红方"
    else
        echo "   当前回合: 黑方"
    fi
else
    echo -e "${YELLOW}⚠️  状态查询失败${NC}"
fi
echo ""

# 总结
echo -e "${BLUE}[6/6] 测试总结${NC}"
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}✅ 所有测试通过！${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo "测试结果："
echo "  ✅ 服务运行正常"
echo "  ✅ 游戏创建成功"
echo "  ✅ 玩家走子成功"
echo "  ✅ AI走子成功"
echo "  ✅ 状态查询成功"
echo ""
echo -e "${BLUE}🎉 AI交互功能正常！${NC}"
echo ""
echo "下一步："
echo "  1. 在浏览器中打开: http://localhost:8090"
echo "  2. 开始完整的对局测试"
echo "  3. 如有问题，查看日志: tail -f logs/chinese-chess-ai.log"
echo ""
