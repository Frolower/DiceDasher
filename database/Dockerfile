FROM postgres:18-alpine AS cron-builder

ARG PG_CRON_VERSION=1.6.7
RUN apk add --no-cache build-base curl \
    && curl -fsSL "https://github.com/citusdata/pg_cron/archive/refs/tags/v${PG_CRON_VERSION}.tar.gz" \
       | tar -xz -C /tmp \
    && cd "/tmp/pg_cron-${PG_CRON_VERSION}" \
    && make with_llvm=no \
    && make with_llvm=no install

FROM postgres:18-alpine
COPY --from=cron-builder /usr/local/lib/postgresql/pg_cron.so /usr/local/lib/postgresql/
COPY --from=cron-builder /usr/local/share/postgresql/extension/pg_cron* /usr/local/share/postgresql/extension/
