

FROM golang:1.19 as builder

# Install deps
RUN apt-get update && apt-get install -y \
  libssl-dev \
  ca-certificates \
  fuse

ENV SRC_DIR /hbnode
ENV OUT /hoodboot

WORKDIR $SRC_DIR

COPY go.mod $SRC_DIR
COPY go.sum $SRC_DIR

RUN go mod download

COPY . $SRC_DIR

RUN go build -o $OUT ./cmd/hoodboot




FROM alpine
ENV APP_DIR /hbnode
ENV OUT /hoodboot
WORKDIR $APP_DIR
COPY --from=builder $OUT $APP_DIR

# Swarm TCP; should be exposed to the public
EXPOSE 4002
# Swarm UDP; should be exposed to the public
EXPOSE 4002/udp

ENTRYPOINT ["./hoodboot"]
