FROM golang:alpine as builder

LABEL maintainer="Denis Chernyavskiy"

RUN apk add --update --no-cache make

WORKDIR /app

COPY . /app

RUN go mod init something
RUN go mod tidy
RUN go build
RUN go install github.com/lib/pq
RUN go get -d github.com/gorilla/mux



RUN make build

FROM alpine:latest as runner

COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["./util"]
