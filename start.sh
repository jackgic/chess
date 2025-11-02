#!/bin/bash

# 中国象棋AI项目启动脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

# Go环境配置
# export PATH="/usr/local/go/bin:$PATH"
# export GOPATH="$HOME/go"
# export PATH="$GOPATH/bin:$PATH"

# 项目配置
PROJECT_NAME="chinese-chess-ai"
LOG_DIR="logs"
PID_FILE="${PROJECT_NAME}.pid"
LOG_FILE="${LOG_DIR}/${PROJECT_NAME}.log"
ERROR_LOG="${LOG_DIR}/${PROJECT_NAME}.error.log"

# 创建日志目录
mkdir -p ${LOG_DIR}

# 环境变量配置
export LKE_APP_ID="ZRmBjPVaArSCCoXuoZQGMxTYxJAJLGzpjjsRasrVRsvmyEBJmQrmjnstFeODPeqDfmymtWymgGUrfnAENRswnOHmsGDJRziSpklkIJHMfRxRTBhZEcXJUCCMFlwMtumv"
export LKE_APP_KEY="ZRmBjPVaArSCCoXuoZQGMxTYxJAJLGzpjjsRasrVRsvmyEBJmQrmjnstFeODPeqDfmymtWymgGUrfnAENRswnOHmsGDJRziSpklkIJHMfRxRTBhZEcXJUCCMFlwMtumv"
export LKE_BOT_BIZ_ID="chinese_chess"
export PORT="8090"

echo "=== 中国象棋AI项目启动 ===" | tee -a ${LOG_FILE}
echo "启动时间: $(date)" | tee -a ${LOG_FILE}

# 检查是否已经在运行
if [ -f ${PID_FILE} ]; then
    OLD_PID=$(cat ${PID_FILE})
    if ps -p ${OLD_PID} > /dev/null 2>&1; then
        echo "错误: 程序已经在运行 (PID: ${OLD_PID})" | tee -a ${LOG_FILE}
        echo "请先运行 ./stop.sh 停止程序" | tee -a ${LOG_FILE}
        exit 1
    else
        echo "清理旧的PID文件..." | tee -a ${LOG_FILE}
        rm -f ${PID_FILE}
    fi
fi

# 编译程序
echo "正在编译程序..." | tee -a ${LOG_FILE}
if go build -o ${PROJECT_NAME} main.go; then
    echo "编译成功!" | tee -a ${LOG_FILE}
else
    echo "编译失败!" | tee -a ${ERROR_LOG}
    exit 1
fi

# 启动程序
echo "正在启动程序..." | tee -a ${LOG_FILE}
echo "服务端口: ${PORT}" | tee -a ${LOG_FILE}
echo "日志文件: ${LOG_FILE}" | tee -a ${LOG_FILE}
echo "错误日志: ${ERROR_LOG}" | tee -a ${LOG_FILE}

# 在后台启动程序，输出到日志文件
nohup ./${PROJECT_NAME} >> ${LOG_FILE} 2>> ${ERROR_LOG} &
PID=$!

# 保存PID
echo ${PID} > ${PID_FILE}

# 等待程序启动
sleep 3

# 检查程序是否成功启动
if ps -p ${PID} > /dev/null 2>&1; then
    echo "程序启动成功!" | tee -a ${LOG_FILE}
    echo "进程ID: ${PID}" | tee -a ${LOG_FILE}
    echo "PID已保存到: ${PID_FILE}" | tee -a ${LOG_FILE}
    echo "" | tee -a ${LOG_FILE}
    echo "使用以下命令:" | tee -a ${LOG_FILE}
    echo "  查看日志: tail -f ${LOG_FILE}" | tee -a ${LOG_FILE}
    echo "  查看错误: tail -f ${ERROR_LOG}" | tee -a ${LOG_FILE}
    echo "  停止程序: ./stop.sh" | tee -a ${LOG_FILE}
    echo "  检查状态: ./status.sh" | tee -a ${LOG_FILE}
    echo "" | tee -a ${LOG_FILE}
    echo "API地址: http://localhost:${PORT}" | tee -a ${LOG_FILE}
else
    echo "程序启动失败!" | tee -a ${ERROR_LOG}
    rm -f ${PID_FILE}
    exit 1
fi
