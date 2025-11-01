#!/bin/bash

# 测试LKE API的不同参数格式

APP_ID="ZRmBjPVaArSCCoXuoZQGMxTYxJAJLGzpjjsRasrVRsvmyEBJmQrmjnstFeODPeqDfmymtWymgGUrfnAENRswnOHmsGDJRziSpklkIJHMfRxRTBhZEcXJUCCMFlwMtumv"
SESSION_ID="test_session_$(date +%s)"

echo "测试LKE API参数格式..."

# 测试格式1: app_id
echo "测试格式1: app_id"
curl -X POST "https://wss.lke.cloud.tencent.com/v1/qbot/chat/sse" \
  -H "Accept: text/event-stream" \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "'$APP_ID'",
    "session_id": "'$SESSION_ID'",
    "content": "你好",
    "visitor_biz_id": "test_user",
    "search_network": "disable",
    "incremental": true
  }' \
  --max-time 10

echo -e "\n\n"

# 测试格式2: bot_app_key
echo "测试格式2: bot_app_key"
curl -X POST "https://wss.lke.cloud.tencent.com/v1/qbot/chat/sse" \
  -H "Accept: text/event-stream" \
  -H "Content-Type: application/json" \
  -d '{
    "bot_app_key": "'$APP_ID'",
    "session_id": "'$SESSION_ID'",
    "content": "你好",
    "visitor_biz_id": "test_user",
    "search_network": "disable",
    "incremental": true
  }' \
  --max-time 10

echo -e "\n\n"

# 测试格式3: 添加更多参数
echo "测试格式3: 添加更多参数"
curl -X POST "https://wss.lke.cloud.tencent.com/v1/qbot/chat/sse" \
  -H "Accept: text/event-stream" \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "'$APP_ID'",
    "session_id": "'$SESSION_ID'",
    "content": "你好",
    "visitor_biz_id": "test_user",
    "search_network": "disable",
    "incremental": true,
    "bot_biz_id": "chinese_chess"
  }' \
  --max-time 10