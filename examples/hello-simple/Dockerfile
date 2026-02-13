FROM golang:1.25.7

WORKDIR /repo/examples/hello-simple

COPY . /repo

RUN go mod download

CMD ["go", "run", "main.go"]
