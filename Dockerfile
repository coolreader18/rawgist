FROM golang:1.10-alpine

WORKDIR /go/src/rawgist
COPY . .

RUN apk add --no-cache git

RUN go get -d -v ./...
RUN go install -v .

RUN apk del git

EXPOSE 443

CMD ["rawgist"]
