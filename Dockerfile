FROM golang:tip-20260131-alpine3.23 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY "./src/go.mod" .
COPY "./src/go.sum" .
RUN ["go", "mod", "download"]

COPY ./src/ .

RUN ["go", "build", "-o", "./smolearl", "."]

FROM scratch

COPY --from=builder /app/smolearl /

CMD ["/smolearl"]