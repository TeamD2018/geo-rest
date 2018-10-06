FROM golang:1.11 AS build
WORKDIR /go/src
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -tags=jsoniter -o geo-rest .

FROM scratch AS runtime
COPY --from=build /go/src/geo-rest ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
EXPOSE 8080/tcp
ENTRYPOINT ["./geo-rest"]
