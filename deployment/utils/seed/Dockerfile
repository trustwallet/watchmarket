FROM alpine:3.7
RUN apk --update add postgresql-client && rm -rf /var/cache/apk/*

ENV STORAGE_POSTGRES_URI=

COPY deployment/utils/seed/seed.sh /seed/seed.sh

COPY seed/watchmarket_public_rates.sql /seed/rates.sql
COPY seed/watchmarket_public_tickers.sql /seed/tickers.sql

ENTRYPOINT [ "sh", "/seed/seed.sh" ]