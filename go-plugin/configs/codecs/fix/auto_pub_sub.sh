#!/bin/bash

export SERVICE_ID="fix-provider@fix" # please change fix-provider to your service identity
export BACKEND_PORT=7755            # please change port 7755 to your java server port
export PROVIDER_APP=fix-provider

export MOCK_PUB_DATA="{\"protocolType\": \"fix\", \"providerMetaInfo\": { \"appName\": \"${PROVIDER_APP}\",\"properties\": {\"application\": \"${PROVIDER_APP}\",\"port\": \"${BACKEND_PORT}\" }},	\"serviceName\": \"${SERVICE_ID}\"}"

export MOCK_SUB_DATA="{\"protocolType\":\"fix\",\"serviceName\":\"${SERVICE_ID}\"}"

echo "publish service ${SERVICE_ID}"
echo "curl -d \"${MOCK_PUB_DATA}\" localhost:13330/services/publish"
curl -s -d "${MOCK_PUB_DATA}" localhost:13330/services/publish

sleep 2

echo
echo
echo "subscribe service ${SERVICE_ID}"
echo "curl -d \"${MOCK_SUB_DATA}\" localhost:13330/services/subscribe"
curl -s -d "${MOCK_SUB_DATA}" localhost:13330/services/subscribe
