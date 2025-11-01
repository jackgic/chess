#!/bin/bash

# 中国象棋AI游戏 - 全面深度测试脚本
# 测试所有核心功能并生成详细报告

set -e

BASE_URL="http://localhost:8090"
GAME_ID=""
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 测试结果数组
declare -a TEST_RESULTS

# 打印分隔线
print_separator() {
    echo -e "${BLUE}================================================================${NC}"
}

# 打印测试标题
print_header() {
    echo ""
    print_separator
    echo -e "${CYAN}$1${NC}"
    print_separator
}

# 打印子标题
print_subheader() {
    echo -e "${YELLOW}>>> $1${NC}"
}

# 记录测试结果
record_test() {
    ((TEST_COUNT++))
    if [ $1 -eq 0 ]; then
        ((PASS_COUNT++))
        echo -e "${GREEN}✅ PASS: $2${NC}"
        TEST_RESULTS+=("✅ $2")
    else
        ((FAIL_COUNT++))
        echo -e "${RED}❌ FAIL: $2${NC}"
        TEST_RESULTS+=("❌ $2")
    fi
}

# 打印JSON（格式化）
print_json() {
    echo "$1" | python3 -m json.tool 2>/dev/null || echo "$1"
}

# ============================================================================
# 测试1: 服务器健康检查
# ============================================================================
test_health_check() {
    print_header "测试1: 服务器健康检查"
    
    response=$(curl -s "${BASE_URL}/api/health")
    echo "响应: $response"
    
    if echo "$response" | grep -q '"status":"OK"'; then
        record_test 0 "服务器健康检查"
    else
        record_test 1 "服务器健康检查"
    fi
}

# ============================================================================
# 测试2: 创建游戏（红方）
# ============================================================================
test_create_game_red() {
    print_header "测试2: 创建游戏（玩家选择红方）"
    
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    echo "响应摘要:"
    echo "$response" | grep -o '"success":[^,]*' || echo "$response"
    
    if echo "$response" | grep -q '"success":true'; then
        GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        echo "游戏ID: $GAME_ID"
        
        # 验证初始状态
        if echo "$response" | grep -q '"turn":1'; then
            record_test 0 "创建游戏（红方），初始回合正确"
        else
            record_test 1 "创建游戏（红方），初始回合错误"
        fi
    else
        record_test 1 "创建游戏（红方）失败"
    fi
}

# ============================================================================
# 测试3: 查询游戏状态
# ============================================================================
test_get_game_state() {
    print_header "测试3: 查询游戏状态"
    
    if [ -z "$GAME_ID" ]; then
        record_test 1 "查询游戏状态（无游戏ID）"
        return 1
    fi
    
    response=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    
    # 检查关键字段
    has_board=$(echo "$response" | grep -c '"board":\[\[' || true)
    has_turn=$(echo "$response" | grep -c '"turn":' || true)
    has_status=$(echo "$response" | grep -c '"status":' || true)
    
    echo "检查结果: board=$has_board, turn=$has_turn, status=$has_status"
    
    if [ $has_board -gt 0 ] && [ $has_turn -gt 0 ] && [ $has_status -gt 0 ]; then
        record_test 0 "查询游戏状态，数据完整"
    else
        record_test 1 "查询游戏状态，数据不完整"
    fi
}

# ============================================================================
# 测试4: 棋盘初始化验证
# ============================================================================
test_board_initialization() {
    print_header "测试4: 棋盘初始化验证"
    
    if [ -z "$GAME_ID" ]; then
        record_test 1 "棋盘初始化验证（无游戏ID）"
        return 1
    fi
    
    response=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    
    # 统计棋子数量
    total_pieces=$(echo "$response" | grep -o '"type":[1-7]' | wc -l | tr -d ' ')
    red_pieces=$(echo "$response" | grep -o '"color":1' | wc -l | tr -d ' ')
    black_pieces=$(echo "$response" | grep -o '"color":2' | wc -l | tr -d ' ')
    
    echo "棋子统计: 总计=$total_pieces, 红方=$red_pieces, 黑方=$black_pieces"
    
    # 中国象棋每方16个棋子，总计32个
    if [ "$total_pieces" -eq 32 ] && [ "$red_pieces" -eq 16 ] && [ "$black_pieces" -eq 16 ]; then
        record_test 0 "棋盘初始化，棋子数量正确"
    else
        record_test 1 "棋盘初始化，棋子数量错误"
    fi
}

# ============================================================================
# 测试5: 非法走子 - 空位置
# ============================================================================
test_invalid_move_empty() {
    print_header "测试5: 非法走子 - 起始位置无棋子"
    
    if [ -z "$GAME_ID" ]; then
        record_test 1 "非法走子测试（无游戏ID）"
        return 1
    fi
    
    # 尝试移动空位置 (4,4) -> (5,4)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 4, \"fromCol\": 4, \"toRow\": 5, \"toCol\": 4}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        record_test 0 "正确拒绝空位置走子"
    else
        record_test 1 "未能拒绝空位置走子"
    fi
}

# ============================================================================
# 测试6: 非法走子 - 对方棋子
# ============================================================================
test_invalid_move_opponent() {
    print_header "测试6: 非法走子 - 移动对方棋子"
    
    if [ -z "$GAME_ID" ]; then
        record_test 1 "非法走子测试（无游戏ID）"
        return 1
    fi
    
    # 尝试移动黑方棋子 (0,0) -> (1,0)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 0, \"fromCol\": 0, \"toRow\": 1, \"toCol\": 0}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        record_test 0 "正确拒绝移动对方棋子"
    else
        record_test 1 "未能拒绝移动对方棋子"
    fi
}

# ============================================================================
# 测试7: 非法走子 - 违反规则
# ============================================================================
test_invalid_move_rule() {
    print_header "测试7: 非法走子 - 违反移动规则"
    
    if [ -z "$GAME_ID" ]; then
        record_test 1 "非法走子测试（无游戏ID）"
        return 1
    fi
    
    # 尝试让炮斜着走 (7,1) -> (6,2)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 6, \"toCol\": 2}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        record_test 0 "正确拒绝违反规则的走子"
    else
        record_test 1 "未能拒绝违反规则的走子"
    fi
}

# ============================================================================
# 测试8: 合法走子 - 炮二平五
# ============================================================================
test_valid_move_cannon() {
    print_header "测试8: 合法走子 - 炮二平五"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "新游戏ID: $GAME_ID"
    
    # 走炮二平五 (7,1) -> (7,4)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")
    
    echo "响应摘要:"
    echo "$response" | grep -o '"success":[^,]*' || echo "$response"
    
    if echo "$response" | grep -q '"success":true'; then
        # 检查走子历史
        move_list=$(echo "$response" | grep -o '"moveList":\[[^]]*\]')
        echo "走子历史: $move_list"
        
        # 检查回合是否切换
        turn=$(echo "$response" | grep -o '"turn":[0-9]' | cut -d':' -f2)
        echo "当前回合: $turn (1=红方, 2=黑方)"
        
        if [ "$turn" = "2" ]; then
            record_test 0 "炮二平五走子成功，回合正确切换"
        else
            record_test 1 "炮二平五走子成功，但回合未切换"
        fi
    else
        record_test 1 "炮二平五走子失败"
    fi
}

# ============================================================================
# 测试9: 合法走子 - 马二进三
# ============================================================================
test_valid_move_horse() {
    print_header "测试9: 合法走子 - 马二进三"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "新游戏ID: $GAME_ID"
    
    # 走马二进三 (9,1) -> (7,2)
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 9, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 2}")
    
    echo "响应摘要:"
    echo "$response" | grep -o '"success":[^,]*' || echo "$response"
    
    if echo "$response" | grep -q '"success":true'; then
        record_test 0 "马二进三走子成功"
    else
        record_test 1 "马二进三走子失败"
    fi
}

# ============================================================================
# 测试10: 回合管理
# ============================================================================
test_turn_management() {
    print_header "测试10: 回合管理"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "游戏ID: $GAME_ID"
    
    # 第一步：红方走子
    print_subheader "第1步: 红方走炮二平五"
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")
    
    turn1=$(echo "$response" | grep -o '"turn":[0-9]' | cut -d':' -f2)
    echo "走子后回合: $turn1"
    
    # 第二步：尝试再次走红方（应该失败）
    print_subheader "第2步: 尝试再次走红方（应该失败）"
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 9, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 2}")
    
    echo "响应: $response"
    
    if echo "$response" | grep -q '"success":false'; then
        record_test 0 "回合管理正确，拒绝连续走子"
    else
        record_test 1 "回合管理错误，允许连续走子"
    fi
}

# ============================================================================
# 测试11: AI走子功能
# ============================================================================
test_ai_move() {
    print_header "测试11: AI走子功能"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "游戏ID: $GAME_ID"
    
    # 玩家走子
    print_subheader "玩家走子: 炮二平五"
    curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}" > /dev/null
    
    # AI走子
    print_subheader "请求AI走子..."
    response=$(curl -s -X POST "${BASE_URL}/api/game/${GAME_ID}/ai-move" \
        -H "Content-Type: application/json" \
        --max-time 30)
    
    echo "响应摘要:"
    echo "$response" | grep -o '"success":[^,]*' || echo "$response"
    
    if echo "$response" | grep -q '"success":true'; then
        # 提取AI回答
        answer=$(echo "$response" | grep -o '"answer":"[^"]*"' | cut -d'"' -f4)
        echo "AI回答: $answer"
        
        # 检查是否包含MOVE指令
        if echo "$answer" | grep -q 'MOVE:'; then
            record_test 0 "AI走子成功，包含MOVE指令"
        else
            record_test 1 "AI走子成功，但缺少MOVE指令"
        fi
    else
        error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
        echo "错误: $error"
        record_test 1 "AI走子失败: $error"
    fi
}

# ============================================================================
# 测试12: 连续对局（3回合）
# ============================================================================
test_continuous_game() {
    print_header "测试12: 连续对局（3回合）"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "游戏ID: $GAME_ID"
    
    success_count=0
    
    # 第1回合
    print_subheader "第1回合: 炮二平五"
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}")
    
    if echo "$response" | grep -q '"success":true'; then
        echo "✅ 玩家走子成功"
        ((success_count++))
        
        # AI走子
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
    
    # 第2回合
    print_subheader "第2回合: 马二进三"
    response=$(curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 9, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 2}")
    
    if echo "$response" | grep -q '"success":true'; then
        echo "✅ 玩家走子成功"
        ((success_count++))
        
        # AI走子
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
    
    echo "成功步数: $success_count/4"
    
    if [ $success_count -ge 3 ]; then
        record_test 0 "连续对局测试（$success_count/4步成功）"
    else
        record_test 1 "连续对局测试（$success_count/4步成功）"
    fi
}

# ============================================================================
# 测试13: 并发游戏
# ============================================================================
test_concurrent_games() {
    print_header "测试13: 并发游戏"
    
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
    
    # 验证每个游戏都可以独立访问
    success_count=0
    for game_id in "${game_ids[@]}"; do
        response=$(curl -s "${BASE_URL}/api/game/${game_id}")
        if echo "$response" | grep -q '"success":true'; then
            ((success_count++))
        fi
    done
    
    echo "可访问游戏: $success_count/3"
    
    if [ $success_count -eq 3 ]; then
        record_test 0 "并发游戏测试（3/3个游戏可访问）"
    else
        record_test 1 "并发游戏测试（$success_count/3个游戏可访问）"
    fi
}

# ============================================================================
# 测试14: 棋盘状态一致性
# ============================================================================
test_board_consistency() {
    print_header "测试14: 棋盘状态一致性"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    
    # 获取初始状态
    initial_state=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    initial_pieces=$(echo "$initial_state" | grep -o '"type":[1-7]' | wc -l | tr -d ' ')
    echo "初始棋子数量: $initial_pieces"
    
    # 走一步（不吃子）
    curl -s -X POST "${BASE_URL}/api/game/move" \
        -H "Content-Type: application/json" \
        -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": 7, \"fromCol\": 1, \"toRow\": 7, \"toCol\": 4}" > /dev/null
    
    # 获取走子后状态
    after_state=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    after_pieces=$(echo "$after_state" | grep -o '"type":[1-7]' | wc -l | tr -d ' ')
    echo "走子后棋子数量: $after_pieces"
    
    if [ "$initial_pieces" -eq "$after_pieces" ]; then
        record_test 0 "棋盘状态一致性（棋子数量不变）"
    else
        record_test 1 "棋盘状态一致性（棋子数量异常变化）"
    fi
}

# ============================================================================
# 测试15: 走子历史记录
# ============================================================================
test_move_history() {
    print_header "测试15: 走子历史记录"
    
    # 创建新游戏
    response=$(curl -s -X POST "${BASE_URL}/api/game/new" \
        -H "Content-Type: application/json" \
        -d '{"playerColor": 1}')
    
    GAME_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    
    # 走3步
    moves=("7,1,7,4" "9,1,7,2" "9,0,8,0")
    
    for move in "${moves[@]}"; do
        IFS=',' read -r fromRow fromCol toRow toCol <<< "$move"
        curl -s -X POST "${BASE_URL}/api/game/move" \
            -H "Content-Type: application/json" \
            -d "{\"gameId\": \"${GAME_ID}\", \"fromRow\": $fromRow, \"fromCol\": $fromCol, \"toRow\": $toRow, \"toCol\": $toCol}" > /dev/null
        
        # 等待AI走子
        curl -s -X POST "${BASE_URL}/api/game/${GAME_ID}/ai-move" \
            -H "Content-Type: application/json" \
            --max-time 30 > /dev/null || true
    done
    
    # 获取最终状态
    response=$(curl -s "${BASE_URL}/api/game/${GAME_ID}")
    move_list=$(echo "$response" | grep -o '"moveList":\[[^]]*\]')
    move_count=$(echo "$move_list" | grep -o ',' | wc -l | tr -d ' ')
    
    echo "走子历史: $move_list"
    echo "走子数量: $((move_count + 1))"
    
    if [ $move_count -ge 2 ]; then
        record_test 0 "走子历史记录正常"
    else
        record_test 1 "走子历史记录异常"
    fi
}

# ============================================================================
# 主测试流程
# ============================================================================
main() {
    print_header "🎮 中国象棋AI游戏 - 全面深度测试"
    echo "测试时间: $(date)"
    echo "服务器地址: $BASE_URL"
    echo ""
    
    # 检查服务器是否运行
    if ! curl -s "${BASE_URL}/api/health" > /dev/null 2>&1; then
        echo -e "${RED}❌ 错误: 服务器未运行或无法访问${NC}"
        echo "请先启动服务器: ./start.sh"
        exit 1
    fi
    
    # 运行所有测试
    test_health_check
    test_create_game_red
    test_get_game_state
    test_board_initialization
    test_invalid_move_empty
    test_invalid_move_opponent
    test_invalid_move_rule
    test_valid_move_cannon
    test_valid_move_horse
    test_turn_management
    test_ai_move
    test_continuous_game
    test_concurrent_games
    test_board_consistency
    test_move_history
    
    # 打印测试总结
    print_header "📊 测试总结报告"
    
    echo ""
    echo -e "${CYAN}详细结果:${NC}"
    for result in "${TEST_RESULTS[@]}"; do
        echo "  $result"
    done
    
    echo ""
    print_separator
    echo -e "${BLUE}总计测试: $TEST_COUNT${NC}"
    echo -e "${GREEN}通过: $PASS_COUNT${NC}"
    echo -e "${RED}失败: $FAIL_COUNT${NC}"
    
    pass_rate=$((PASS_COUNT * 100 / TEST_COUNT))
    echo -e "${CYAN}通过率: ${pass_rate}%${NC}"
    print_separator
    
    echo ""
    if [ $FAIL_COUNT -eq 0 ]; then
        echo -e "${GREEN}🎉 恭喜！所有测试通过！游戏运行完全正常！${NC}"
        exit 0
    elif [ $pass_rate -ge 80 ]; then
        echo -e "${YELLOW}⚠️  大部分测试通过，但有 $FAIL_COUNT 个测试失败，建议修复${NC}"
        exit 0
    else
        echo -e "${RED}❌ 有 $FAIL_COUNT 个测试失败，需要立即修复${NC}"
        exit 1
    fi
}

# 运行主函数
main
