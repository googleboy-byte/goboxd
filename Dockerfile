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
    bash \
    curl \
    g++ \
    gcc \
    gfortran \
    iverilog \
    libcap2 \
    libnl-route-3-200 \
    libprotobuf32 \
    nodejs \
    npm \
    openjdk-17-jdk-headless \
    python3 \
    unzip \
    wget \
    xz-utils \
# Bonus languages (commented out for Stage 2)
    golang-go \
    kotlin \
    lua5.4 \
    mono-devel \
    ocaml \
    ruby \
    rustc \
    && rm -rf /var/lib/apt/lists/*

# Install Swift and Zig (commented out for Stage 2)
RUN wget -q https://download.swift.org/swift-6.0.2-release/debian12/swift-6.0.2-RELEASE/swift-6.0.2-RELEASE-debian12.tar.gz \
    && tar -xzf swift-6.0.2-RELEASE-debian12.tar.gz -C /usr/local --strip-components=2 \
    && rm swift-6.0.2-RELEASE-debian12.tar.gz
RUN wget -q https://ziglang.org/download/0.13.0/zig-linux-x86_64-0.13.0.tar.xz \
    && tar -xJf zig-linux-x86_64-0.13.0.tar.xz -C /usr/local \
    && ln -sf /usr/local/zig-linux-x86_64-0.13.0/zig /usr/local/bin/zig \
    && rm zig-linux-x86_64-0.13.0.tar.xz

# Install Dart SDK and TypeScript
RUN wget -q https://storage.googleapis.com/dart-archive/channels/stable/release/3.4.4/sdk/dartsdk-linux-x64-release.zip \
    && unzip -q dartsdk-linux-x64-release.zip -d /usr/local \
    && ln -sf /usr/local/dart-sdk/bin/dart /usr/local/bin/dart \
    && ln -sf /usr/local/dart-sdk/bin/dart /usr/bin/dart \
    && rm dartsdk-linux-x64-release.zip \
    && npm install -g typescript @types/node \
    && ln -sf /usr/local/bin/tsc /usr/bin/tsc

COPY --from=nsjail-builder /nsjail-src/nsjail /usr/sbin/nsjail
COPY --from=builder /build/goboxd /usr/local/bin/goboxd
COPY languages.yaml /etc/goboxd/languages.yaml

EXPOSE 8080

CMD ["/usr/local/bin/goboxd", "-port", "8080", "-config", "/etc/goboxd/languages.yaml"]
