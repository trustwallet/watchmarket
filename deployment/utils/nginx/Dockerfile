FROM nginx:alpine

ENV API_CHARTS_URI=
ENV API_INFO_URI=
ENV API_TICKERS_URI=
ENV API_SWAGGER_URI=

RUN apk update && apk --no-cache add gettext

COPY watchmarket.conf /template/watchmarket.conf
COPY nginx.conf /etc/nginx/nginx.conf

COPY entrypoint.sh /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT [ "sh",  "/app/entrypoint.sh" ]