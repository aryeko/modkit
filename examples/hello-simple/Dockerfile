FROM golang:1.25.7

WORKDIR /repo/examples/hello-simple

COPY . /repo

RUN go mod download

RUN useradd -m -u 10001 appuser \
  && chown -R appuser:appuser /repo /go

ENV GOCACHE=/tmp/go-cache

USER appuser

CMD ["go", "run", "main.go"]
