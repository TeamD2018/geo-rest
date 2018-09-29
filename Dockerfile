FROM golang:1.11 AS build
WORKDIR /go/src
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -tags=jsoniter -o openapi .

FROM scratch AS runtime
COPY --from=build . ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
