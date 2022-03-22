FROM golang:1.17 AS builder
COPY . /spam
WORKDIR /spam
RUN make build

FROM debian:buster
COPY --from=builder /spam ./
CMD ./spam