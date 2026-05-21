FROM golang:1.23-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-X main.version=0.1.0 -X main.commit=$(git rev-parse --short HEAD)" -o goboxd ./cmd/goboxd

FROM debian:bookworm-slim AS nsjail-builder

RUN apt-get update && apt-get install -y \
    bison \
    flex \
    g++ \
    gcc \
    git \
    libcap-dev \
    libnl-route-3-dev \
    libprotobuf-dev \
    make \
    pkg-config \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

COPY external/nsjail /nsjail-src
WORKDIR /nsjail-src
RUN make -j$(nproc)

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    g++ \
    libcap2 \
    libnl-route-3-200 \
    libprotobuf32 \
    python3 \
    rustc \
    && rm -rf /var/lib/apt/lists/*

COPY --from=nsjail-builder /nsjail-src/nsjail /usr/sbin/nsjail
COPY --from=builder /build/goboxd /usr/local/bin/goboxd
COPY languages.yaml /etc/goboxd/languages.yaml

EXPOSE 8080

CMD ["/usr/local/bin/goboxd", "-port", "8080", "-config", "/etc/goboxd/languages.yaml"]
