#!/bin/bash

# 中国象棋AI项目日志查看脚本
# 作者: AI Assistant

PROJECT_NAME="chinese-chess-ai"
LOG_DIR="logs"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}.log"
ERROR_LOG="${LOG_DIR}/${PROJECT_NAME}.error.log"

# 显示使用帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -f, --follow     实时跟踪日志 (类似 tail -f)"
    echo "  -e, --error      查看错误日志"
    echo "  -n, --lines NUM  显示最后 NUM 行 (默认50)"
    echo "  -c, --clear      清空日志文件"
    echo "  -h, --help       显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                查看最后50行日志"
    echo "  $0 -f             实时跟踪日志"
    echo "  $0 -e             查看错误日志"
    echo "  $0 -n 100         查看最后100行日志"
    echo "  $0 -c             清空所有日志"
}

# 默认参数
FOLLOW=false
ERROR_LOG_MODE=false
LINES=50
CLEAR_LOGS=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--follow)
            FOLLOW=true
            shift
            ;;
        -e|--error)
            ERROR_LOG_MODE=true
            shift
            ;;
        -n|--lines)
            LINES="$2"
            shift 2
            ;;
        -c|--clear)
            CLEAR_LOGS=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 清空日志
if [ "$CLEAR_LOGS" = true ]; then
    echo "正在清空日志文件..."
    if [ -f "$LOG_FILE" ]; then
        > "$LOG_FILE"
        echo "已清空: $LOG_FILE"
    fi
    if [ -f "$ERROR_LOG" ]; then
        > "$ERROR_LOG"
        echo "已清空: $ERROR_LOG"
    fi
    echo "日志文件已清空"
    exit 0
fi

# 选择要查看的日志文件
if [ "$ERROR_LOG_MODE" = true ]; then
    TARGET_LOG="$ERROR_LOG"
    LOG_TYPE="错误日志"
else
    TARGET_LOG="$LOG_FILE"
    LOG_TYPE="主日志"
fi

# 检查日志文件是否存在
if [ ! -f "$TARGET_LOG" ]; then
    echo "日志文件不存在: $TARGET_LOG"
    echo "程序可能还没有运行过，请先运行 ./start.sh"
    exit 1
fi

echo "=== 查看${LOG_TYPE}: $TARGET_LOG ==="
echo ""

# 实时跟踪或显示指定行数
if [ "$FOLLOW" = true ]; then
    echo "实时跟踪日志 (按 Ctrl+C 退出)..."
    echo ""
    tail -f "$TARGET_LOG"
else
    echo "显示最后 $LINES 行:"
    echo ""
    tail -n "$LINES" "$TARGET_LOG"
fi