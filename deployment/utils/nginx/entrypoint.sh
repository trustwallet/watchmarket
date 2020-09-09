set -e
# ENVs must be set
test -n $API_CHARTS_URI || (echo "API_CHARTS_URI argument not set" && false)
test -n $API_INFO_URI || (echo "API_INFO_URI argument not set" && false)
test -n $API_TICKERS_URI || (echo "API_TICKERS_URI argument not set" && false)
test -n $API_SWAGGER_URI || (echo "API_SWAGGER_URI argument not set" && false)
test -n $API_RATES_URI || (echo "API_RATES_URI argument not set" && false)

NGINX_VARS='$API_CHARTS_URI $API_INFO_URI $API_TICKERS_URI $API_SWAGGER_URI $API_RATES_URI'
envsubst "$NGINX_VARS" < /template/watchmarket.conf > /etc/nginx/conf.d/default.conf

# Cat configs
echo "API_CHARTS_URI = $API_CHARTS_URI"
echo "API_INFO_URI = $API_INFO_URI"
echo "API_TICKERS_URI = $API_TICKERS_URI"
echo "API_SWAGGER_URI = $API_SWAGGER_URI"
echo "API_RATES_URI = $API_RATES_URI"
echo "----------------------------------"

nginx -c /etc/nginx/nginx.conf -t

echo "proxy started successfully"

nginx -g 'daemon off;'