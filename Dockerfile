FROM golang:latest as builder

WORKDIR /build
COPY . ./
RUN go install

FROM debian:11-slim as run
WORKDIR /run
RUN apt -y update && \
    apt -y install stockfish && \
    apt -y autoclean && \
    apt clean

ENV STOCKFISH_PATH /usr/games/stockfish
COPY --from=builder /go/bin/edi /run/edi 

CMD ["/run/edi"]