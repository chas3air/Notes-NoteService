
FROM --platform=$BUILDPLATFORM golang:1.24 AS build

WORKDIR /src

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

# Сборка бинарника с CGO
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cli ./cmd/app

RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /cli ./cmd/app

FROM debian:bookworm-slim

WORKDIR /

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=build /cli /cli

COPY --from=build /src/migrations ./migrations
COPY config.yaml config.yaml

EXPOSE 8080
EXPOSE 50051
EXPOSE 9090

CMD ["./cli", "--config=config.yaml"]
