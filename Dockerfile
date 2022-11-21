FROM golang:1.19

# Install deps
RUN apt-get update && apt-get install -y \
  libssl-dev \
  ca-certificates \
  fuse

ENV SRC_DIR /hbnode

WORKDIR $SRC_DIR

COPY go.mod $SRC_DIR
COPY go.sum $SRC_DIR

RUN go mod download

COPY . $SRC_DIR

RUN go build -o hbnode

VOLUME $SRC_DIR/data
# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm UDP; should be exposed to the public
EXPOSE 4001/udp

CMD ["./hbnode"]
