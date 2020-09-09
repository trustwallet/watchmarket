FROM alpine:3.7
RUN apk --update add postgresql-client && rm -rf /var/cache/apk/*

ENV STORAGE_POSTGRES_URI=

COPY wait_until_pg_ready.sh /check/wait_until_pg_ready.sh

ENTRYPOINT [ "sh", "/check/wait_until_pg_ready.sh" ]