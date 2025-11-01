#!/bin/bash

# 中国象棋AI项目停止脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

# 项目配置
PROJECT_NAME="chinese-chess-ai"
PID_FILE="${PROJECT_NAME}.pid"
LOG_DIR="logs"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}.log"

echo "=== 中国象棋AI项目停止 ==="
echo "停止时间: $(date)"

# 检查PID文件是否存在
if [ ! -f ${PID_FILE} ]; then
    echo "警告: PID文件不存在 (${PID_FILE})"
    echo "程序可能没有运行或已经停止"
    
    # 尝试通过进程名查找并停止
    PIDS=$(pgrep -f ${PROJECT_NAME} 2>/dev/null || true)
    if [ -n "${PIDS}" ]; then
        echo "发现运行中的进程，正在停止..."
        for pid in ${PIDS}; do
            echo "停止进程: ${pid}"
            kill ${pid} 2>/dev/null || true
        done
        sleep 2
        
        # 强制停止仍在运行的进程
        REMAINING_PIDS=$(pgrep -f ${PROJECT_NAME} 2>/dev/null || true)
        if [ -n "${REMAINING_PIDS}" ]; then
            echo "强制停止剩余进程..."
            for pid in ${REMAINING_PIDS}; do
                echo "强制停止进程: ${pid}"
                kill -9 ${pid} 2>/dev/null || true
            done
        fi
        echo "所有相关进程已停止"
    else
        echo "没有发现运行中的${PROJECT_NAME}进程"
    fi
    exit 0
fi

# 读取PID
PID=$(cat ${PID_FILE})
echo "读取到PID: ${PID}"

# 检查进程是否存在
if ! ps -p ${PID} > /dev/null 2>&1; then
    echo "警告: 进程 ${PID} 不存在"
    echo "清理PID文件..."
    rm -f ${PID_FILE}
    echo "PID文件已清理"
    exit 0
fi

# 优雅停止进程
echo "正在停止进程 ${PID}..."
kill ${PID}

# 等待进程停止
WAIT_TIME=0
MAX_WAIT=10
while ps -p ${PID} > /dev/null 2>&1 && [ ${WAIT_TIME} -lt ${MAX_WAIT} ]; do
    echo "等待进程停止... (${WAIT_TIME}/${MAX_WAIT})"
    sleep 1
    WAIT_TIME=$((WAIT_TIME + 1))
done

# 检查进程是否已停止
if ps -p ${PID} > /dev/null 2>&1; then
    echo "进程未能优雅停止，强制终止..."
    kill -9 ${PID}
    sleep 1
    
    if ps -p ${PID} > /dev/null 2>&1; then
        echo "错误: 无法停止进程 ${PID}"
        exit 1
    else
        echo "进程已强制停止"
    fi
else
    echo "进程已优雅停止"
fi

# 清理PID文件
rm -f ${PID_FILE}
echo "PID文件已清理"

# 记录停止信息到日志
if [ -f ${LOG_FILE} ]; then
    echo "程序停止时间: $(date)" >> ${LOG_FILE}
    echo "================================" >> ${LOG_FILE}
fi

echo "程序已成功停止!"