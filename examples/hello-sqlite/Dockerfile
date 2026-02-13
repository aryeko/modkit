FROM golang:1.25.7

WORKDIR /repo/examples/hello-sqlite

COPY . /repo

RUN apt-get update \
  && apt-get install -y --no-install-recommends gcc libsqlite3-dev \
  && rm -rf /var/lib/apt/lists/*

RUN go mod download

CMD ["go", "run", "./cmd/api"]
