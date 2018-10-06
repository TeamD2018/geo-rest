FROM golang:1.11 as build
WORKDIR /go/src
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -tags=jsoniter -o geo-rest .

FROM alpine:3.8 AS runtime

COPY --from=build /go/src/geo-rest ./
RUN apk add --update ca-certificates

EXPOSE 8080/tcp
ENV GIN_MODE=release
ENTRYPOINT ["./geo-rest"]
