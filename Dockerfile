FROM golang:1.11 AS build
WORKDIR /go/src
COPY models ./go
COPY main.go .
COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -tags=jsoniter -o openapi .

FROM scratch AS runtime
COPY --from=build /models ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
