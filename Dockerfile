FROM golang:1.21.0

COPY . /app
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/rubenv/sql-migrate/...@latest && \
    pwd && ls && \
    go mod tidy && \
    wire && \
    sql-migrate down && \
    sql-migrate up

