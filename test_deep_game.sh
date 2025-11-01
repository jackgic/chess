#!/bin/bash

# 深度游戏测试脚本
# 测试游戏的所有核心功能

set -e

BASE_URL="http://localhost:8090"
GAME_ID=""
TEST_RESULTS=()

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印测试标题
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# 打印测试结果
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
        TEST_RESULTS+=("✅ $2")
    else
        echo -e "${RED}❌ $2${NC}"
        TEST_RESULTS+=("❌ $2")
    fi
}

# 测试1: 健康检查
test_health_check() {
    print_header "测试1: 健康检查"
    
    response=$(curl -s "${BASE_URL}/api/health")
    echo "响应: $response"
    
    if echo "$response" | grep -q "OK"; then
        print_result 0 "健康检查通过"
        return 0
    else
        print_result 1 "健康检查失败"
        return 1
    fi
}

# 测试2: 创建游戏（红方）
test_create_game_red() {
    print_header "测试2: 创建游戏（玩家选择红方）"
    
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        echo "游戏ID: $GAME_ID"
        print_result 0 "创建游戏成功（红方）"
        return 0
    else
        print_result 1 "创建游戏失败（红方）"
        return 1
    fi
}

# 测试3: 创建游戏（黑方）
test_create_game_black() {
    print_header "测试3: 创建游戏（玩家选择黑方）"
    
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 2}')
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        echo "游戏ID: $GAME_ID"
        
        # 检查是否是AI先手
        if echo "$response" | grep -q '"turn":1'; then
            echo "当前回合: 红方（AI先手）"
            print_result 0 "创建游戏成功（黑方），AI先手"
        else
            print_result 1 "创建游戏成功但回合不对"
        fi
        return 0
    else
        print_result 1 "创建游戏失败（黑方）"
        return 1
    fi
}

# 测试4: 查询游戏状态
test_get_game_state() {
    print_header "测试4: 查询游戏状态"
    
    if [ -z "$GAME_ID" ]; then
        print_result 1 "没有游戏ID，跳过测试"
        return 1
    fi
    
    response=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        # 检查棋盘数据
        if echo "$response" | grep -q '"board":\[\['; then
            print_result 0 "查询游戏状态成功，棋盘数据完整"
            return 0
        else
            print_result 1 "查询游戏状态成功但棋盘数据缺失"
            return 1
        fi
    else
        print_result 1 "查询游戏状态失败"
        return 1
    fi
}

# 测试5: 非法走子（起始位置无棋子）
test_invalid_move_empty() {
    print_header "测试5: 非法走子测试（起始位置无棋子）"
    
    if [ -z "$GAME_ID" ]; then
        print_result 1 "没有游戏ID，跳过测试"
        return 1
    fi
    
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 4, \"fromCol\": 4, \"toRow\": 5, \"toCol\": 4}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        print_result 0 "正确拒绝非法走子（空位置）"
        return 0
    else
        print_result 1 "未能拒绝非法走子（空位置）"
        return 1
    fi
}

# 测试6: 非法走子（移动对方棋子）
test_invalid_move_opponent() {
    print_header "测试6: 非法走子测试（移动对方棋子）"
    
    if [ -z "$GAME_ID" ]; then
        print_result 1 "没有游戏ID，跳过测试"
        return 1
    fi
    
    # 尝试移动黑方的棋子（假设玩家是红方）
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 0, \"fromCol\": 0, \"toRow\": 1, \"toCol\": 0}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        print_result 0 "正确拒绝非法走子（对方棋子）"
        return 0
    else
        print_result 1 "未能拒绝非法走子（对方棋子）"
        return 1
    fi
}

# 测试7: 合法走子（炮二平五）
test_valid_move_cannon() {
    print_header "测试7: 合法走子测试（炮二平五）"
    
    # 创建新游戏（红方）
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "新游戏ID: $GAME_ID"
    
    # 走炮二平五 (7,1) -> (7,4)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        # 检查走子历史
        if echo "$response" | grep -q '"moveList"'; then
            print_result 0 "合法走子成功（炮二平五）"
            return 0
        else
            print_result 1 "走子成功但历史记录缺失"
            return 1
        fi
    else
        print_result 1 "合法走子失败"
        return 1
    fi
}

# 测试8: 合法走子（马二进三）
test_valid_move_horse() {
    print_header "测试8: 合法走子测试（马二进三）"
    
    if [ -z "$GAME_ID" ]; then
        print_result 1 "没有游戏ID，跳过测试"
        return 1
    fi
    
    # 走马二进三 (9,1) -> (7,2)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 9, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 2}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        print_result 0 "合法走子成功（马二进三）"
        return 0
    else
        print_result 1 "合法走子失败（马二进三）"
        return 1
    fi
}

# 测试9: AI走子
test_ai_move() {
    print_header "测试9: AI走子测试"
    
    if [ -z "$GAME_ID" ]; then
        print_result 1 "没有游戏ID，跳过测试"
        return 1
    fi
    
    echo "请求AI走子..."
    response=$(curl -s -X POST "${BASE_URL}/api/game/${GAME_ID}/ai-move" \
        -H "Content-Type: application/json" \
        --max-time 30)
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":true'; then
        # 检查AI回答
        if echo "$response" | grep -q '"answer"'; then
            echo "AI回答存在"
            print_result 0 "AI走子成功"
            return 0
        else
            print_result 1 "AI走子成功但无回答"
            return 1
        fi
    else
        error_msg=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
        echo "错误信息: $error_msg"
        print_result 1 "AI走子失败: $error_msg"
        return 1
    fi
}

# 测试10: 连续对局
test_continuous_game() {
    print_header "测试10: 连续对局测试（3回合）"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "游戏ID: $GAME_ID"
    
    # 定义3步走子
    moves=(
        "7,1,7,4"   # 炮二平五
        "9,1,7,2"   # 马二进三
        "9,0,8,0"   # 车一平二
    )
    
    move_names=(
        "炮二平五"
        "马二进三"
        "车一平二"
    )
    
    success_count=0
    
    for i in "${!moves[@]}"; do
        IFS=',' read -r fromRow fromCol toRow toCol <<< "${moves[$i]}"
        echo ""
        echo "第$((i+1))步: ${move_names[$i]} ($fromRow,$fromCol)->($toRow,$toCol)"
        
        # 玩家走子
        response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
            -H "Content-Type: application/json" \
            -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": $fromRow, \"fromCol\": $fromCol, \"toRow\": $toRow, \"toCol\": $toCol}")
        
        if echo "$response" | grep -q '"success":true'; then
            echo "✅ 玩家走子成功"
            ((success_count++))
            
            # AI走子
            echo "等待AI走子..."
            ai_response=$(curl -s -X POST "${BASE_URL}/api/game/${GAME_ID}/ai-move" \
                -H "Content-Type: application/json" \
                --max-time 30)
            
            if echo "$ai_response" | grep -q '"success":true'; then
                echo "✅ AI走子成功"
                ((success_count++))
            else
                echo "❌ AI走子失败"
            fi
        else
            echo "❌ 玩家走子失败"
        fi
        
        sleep 1
    done
    
    if [ $success_count -ge 4 ]; then
        print_result 0 "连续对局测试通过（$success_count/6步成功）"
        return 0
    else
        print_result 1 "连续对局测试失败（$success_count/6步成功）"
        return 1
    fi
}

# 测试11: 棋盘状态一致性
test_board_consistency() {
    print_header "测试11: 棋盘状态一致性测试"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    
    # 获取初始状态
    initial_state=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    initial_pieces=$(echo "$initial_state" | grep -o '"type":[0-9]' | wc -l)
    echo "初始棋子数量: $initial_pieces"
    
    # 走一步
    curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}" > /dev/null
    
    # 获取走子后状态
    after_state=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    after_pieces=$(echo "$after_state" | grep -o '"type":[0-9]' | wc -l)
    echo "走子后棋子数量: $after_pieces"
    
    if [ "$initial_pieces" -eq "$after_pieces" ]; then
        print_result 0 "棋盘状态一致性测试通过"
        return 0
    else
        print_result 1 "棋盘状态一致性测试失败（棋子数量变化）"
        return 1
    fi
}

# 测试12: 并发游戏
test_concurrent_games() {
    print_header "测试12: 并发游戏测试"
    
    game_ids=()
    
    # 创建3个游戏
    for i in {1..3}; do
        response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
            -H "Content-Type: application/json" \
            -d '{"playerColor": 1}')
        
        game_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        game_ids+=("$game_id")
        echo "创建游戏 $i: $game_id"
    done
    
    # 验证每个游戏都可以独立操作
    success_count=0
    for game_id in "${game_ids[@]}"; do
        response=$(curl -s "${BASE_URL}/api/game/${game_id}")
        if echo "$response" | grep -q '"success":true'; then
            ((success_count++))
        fi
    done
    
    if [ $success_count -eq 3 ]; then
        print_result 0 "并发游戏测试通过（3/3个游戏可访问）"
        return 0
    else
        print_result 1 "并发游戏测试失败（$success_count/3个游戏可访问）"
        return 1
    fi
}

# 主测试流程
main() {
    print_header "🎮 中国象棋AI游戏深度测试"
    echo "测试时间: $(date)"
    echo "服务器地址: $BASE_URL"
    
    # 运行所有测试
    test_health_check
    test_create_game_red
    test_get_game_state
    test_invalid_move_empty
    test_invalid_move_opponent
    test_valid_move_cannon
    test_valid_move_horse
    test_ai_move
    test_create_game_black
    test_continuous_game
    test_board_consistency
    test_concurrent_games
    
    # 打印测试总结
    print_header "📊 测试总结"
    
    passed=0
    failed=0
    
    for result in "${TEST_RESULTS[@]}"; do
        echo "$result"
        if [[ $result == ✅* ]]; then
            ((passed++))
        else
            ((failed++))
        fi
    done
    
    echo ""
    echo -e "${BLUE}总计: $((passed + failed)) 个测试${NC}"
    echo -e "${GREEN}通过: $passed${NC}"
    echo -e "${RED}失败: $failed${NC}"
    
    if [ $failed -eq 0 ]; then
        echo ""
        echo -e "${GREEN}🎉 所有测试通过！游戏运行正常！${NC}"
        exit 0
    else
        echo ""
        echo -e "${RED}⚠️  有 $failed 个测试失败，需要修复${NC}"
        exit 1
    fi
}

# 运行主函数
main
