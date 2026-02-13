FROM golang:1.25.7

WORKDIR /repo/examples/hello-sqlite

COPY . /repo

RUN apt-get update \
  && apt-get install -y --no-install-recommends gcc libsqlite3-dev \
  && rm -rf /var/lib/apt/lists/*

RUN go mod download

RUN useradd -m -u 10001 appuser \
  && chown -R appuser:appuser /repo /go

ENV GOCACHE=/tmp/go-cache

USER appuser

CMD ["go", "run", "./cmd/api"]
