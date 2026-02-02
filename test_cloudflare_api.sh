#!/bin/bash

# Cloudflare API 测试脚本
# 用于诊断 IP4P Cloudflare API 配置问题

echo "=========================================="
echo "Cloudflare API 诊断工具"
echo "=========================================="
echo ""

# 读取配置
read -p "请输入 API Token: " API_TOKEN
read -p "请输入 Zone ID: " ZONE_ID
read -p "请输入域名 (例如: ip4p.tianai.edu.kg): " DOMAIN

echo ""
echo "=========================================="
echo "1. 测试 API Token 有效性"
echo "=========================================="

TOKEN_VERIFY=$(curl -s -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json")

echo "$TOKEN_VERIFY" | python3 -m json.tool 2>/dev/null || echo "$TOKEN_VERIFY"

if echo "$TOKEN_VERIFY" | grep -q '"success":true'; then
    echo "✅ API Token 有效"
else
    echo "❌ API Token 无效或已过期"
    exit 1
fi

echo ""
echo "=========================================="
echo "2. 测试 Zone 访问权限"
echo "=========================================="

ZONE_INFO=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$ZONE_ID" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json")

echo "$ZONE_INFO" | python3 -m json.tool 2>/dev/null || echo "$ZONE_INFO"

if echo "$ZONE_INFO" | grep -q '"success":true'; then
    echo "✅ Zone 访问正常"
    ZONE_NAME=$(echo "$ZONE_INFO" | grep -o '"name":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "   Zone 名称: $ZONE_NAME"
else
    echo "❌ 无法访问 Zone，请检查 Zone ID 是否正确"
    exit 1
fi

echo ""
echo "=========================================="
echo "3. 查询 TXT 记录"
echo "=========================================="

DNS_RECORDS=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records?type=TXT&name=$DOMAIN" \
  -H "Authorization: Bearer $API_TOKEN" \
  -H "Content-Type: application/json")

echo "$DNS_RECORDS" | python3 -m json.tool 2>/dev/null || echo "$DNS_RECORDS"

if echo "$DNS_RECORDS" | grep -q '"success":true'; then
    echo "✅ TXT 记录查询成功"

    RECORD_COUNT=$(echo "$DNS_RECORDS" | grep -o '"result":\[' | wc -l)
    if [ "$RECORD_COUNT" -gt 0 ]; then
        echo "   找到 TXT 记录"

        # 尝试解析 TXT 记录内容
        CONTENT=$(echo "$DNS_RECORDS" | grep -o '"content":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ -n "$CONTENT" ]; then
            echo "   TXT 记录内容: $CONTENT"

            # 尝试 base64 解码
            DECODED=$(echo "$CONTENT" | base64 -d 2>/dev/null)
            if [ $? -eq 0 ]; then
                echo "   解码后内容: $DECODED"
            fi
        fi
    else
        echo "⚠️  未找到 TXT 记录"
        echo "   请在 Cloudflare 中添加 TXT 记录"
    fi
else
    echo "❌ TXT 记录查询失败"
    exit 1
fi

echo ""
echo "=========================================="
echo "4. 测试 DNS 解析"
echo "=========================================="

echo "标准 DNS 查询:"
dig $DOMAIN TXT +short 2>/dev/null || nslookup -type=TXT $DOMAIN 2>/dev/null || echo "dig/nslookup 命令不可用"

echo ""
echo "=========================================="
echo "5. 生成配置"
echo "=========================================="

cat << EOF
建议的 frpc.toml 配置:

serverAddr = "$DOMAIN"

[ip4p]
mode = "api"
provider = "cloudflare"
apiKey = "$API_TOKEN"
zoneId = "$ZONE_ID"

[auth]
method = "token"
token = "your-auth-token"

[log]
to = "./frpc.log"
level = "debug"
EOF

echo ""
echo "=========================================="
echo "诊断完成"
echo "=========================================="
