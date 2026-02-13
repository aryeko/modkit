FROM golang:1.25.7

WORKDIR /repo/examples/hello-mysql

COPY . /repo

RUN go mod download

RUN useradd -m -u 10001 appuser \
  && chown -R appuser:appuser /repo /go

ENV GOCACHE=/tmp/go-cache

USER appuser

CMD ["sh", "-c", "go run ./cmd/migrate && go run ./cmd/seed && go run ./cmd/api"]
