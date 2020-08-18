FROM alpine:3.7
RUN apk --update add postgresql-client && rm -rf /var/cache/apk/*

ENV STORAGE_POSTGRES_URI=

COPY seed.sh /seed/seed.sh

COPY watchmarket_public_rates.sql /seed/rates.sql
COPY watchmarket_public_tickers.sql /seed/tickers.sql

ENTRYPOINT [ "sh", "/seed/seed.sh" ]