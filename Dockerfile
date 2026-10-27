FROM golang:1.24-alpine AS build

RUN apk add --no-cache build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/byggabat .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir /data \
    && chown app:app /data
WORKDIR /app
COPY --from=build --chown=app:app /out/byggabat ./byggabat
COPY --chown=app:app templates ./templates
COPY --chown=app:app static ./static
ENV LISTEN_ADDR=0.0.0.0:8080
ENV DATABASE_PATH=/data/byggabat.db
EXPOSE 8080
VOLUME ["/data"]
USER app
ENTRYPOINT ["./byggabat"]
