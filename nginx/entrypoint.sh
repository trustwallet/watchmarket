# ENVs must be set
test -n $API_CHARTS_URI || (echo "SERVICE argument not set" && false)
test -n $API_INFO_URI || (echo "SERVICE argument not set" && false)
test -n $API_TICKERS_URI || (echo "SERVICE argument not set" && false)
test -n $API_SWAGGER_URI || (echo "SERVICE argument not set" && false)

NGINX_VARS='$API_CHARTS_URI $API_INFO_URI $API_TICKERS_URI $API_SWAGGER_URI'
envsubst "$NGINX_VARS" < /template/watchmarket.conf > /etc/nginx/conf.d/default.conf

# Cat configs
echo "API_CHARTS_URI = $API_CHARTS_URI"
echo "API_INFO_URI = $API_INFO_URI"
echo "API_TICKERS_URI = $API_TICKERS_URI"
echo "API_SWAGGER_URI = $API_SWAGGER_URI"
echo "----------------------------------"
cat /etc/nginx/conf.d/default.conf

nginx -c /etc/nginx/nginx.conf -t

nginx -g 'daemon off;'