FROM golang:alpine3.11 AS base

WORKDIR /app/bookshelf

RUN apk update; \
    apk add make gcc g++

COPY . /app/bookshelf/
RUN go mod download

# --- #

FROM base AS bookshelf

RUN make build

ENTRYPOINT ["./bookshelf"]