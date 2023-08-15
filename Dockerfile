FROM golang:bookworm as builder

WORKDIR /build
COPY . ./
RUN go install

FROM debian:12-slim as run
WORKDIR /run
RUN apt -y update && \
    apt -y install stockfish ca-certificates && \
    apt -y autoclean && \
    apt clean

ENV STOCKFISH_PATH /usr/games/stockfish
COPY --from=builder /go/bin/edi /bin/edi 

ENTRYPOINT ["/bin/edi"]
