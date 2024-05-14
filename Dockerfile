FROM golang:alpine

WORKDIR /app

RUN apk update && \
    apk add --no-cache bash \
                       curl \
                       grep \
                       sed

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build

EXPOSE 8080

CMD ["/build"]
