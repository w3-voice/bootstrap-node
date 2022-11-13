FROM golang:1.19

# Install deps
RUN apt-get update && apt-get install -y \
  libssl-dev \
  ca-certificates \
  fuse

ENV SRC_DIR /hbnode

RUN cd $SRC_DIR \
  && go build -o hbnode


# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm UDP; should be exposed to the public
EXPOSE 4001/udp

CMD [$SRC_DIR/hbnode]

