FROM golang:1.25.7

WORKDIR /repo/examples/hello-postgres

COPY . /repo

RUN go mod download

CMD ["go", "run", "./cmd/api"]
