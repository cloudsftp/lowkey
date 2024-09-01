# syntax=docker/dockerfile:1

################################################################################
# Build stage

ARG RUST_VERSION=1.80.0
ARG APP_NAME=lowkey
FROM rust:${RUST_VERSION} AS build
ARG APP_NAME
WORKDIR /app

RUN --mount=type=bind,source=src,target=src \
    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
    --mount=type=cache,target=/app/target/ \
    --mount=type=cache,target=/usr/local/cargo/registry/ \
    <<EOF
set -e
cargo build --locked --release
cp ./target/release/$APP_NAME /bin/server
EOF

################################################################################
# Run stage
#FROM ubuntu:latest AS final
FROM rust:${RUST_VERSION} AS final

#ARG UID=10001
#RUN adduser \
#    --disabled-password \
#    --gecos "" \
#    --home "/nonexistent" \
#    --shell "/sbin/nologin" \
#    --no-create-home \
#    --uid "${UID}" \
#    appuser
#USER appuser

COPY --from=build /bin/server /bin/

EXPOSE 6670

CMD ["/bin/server"]
