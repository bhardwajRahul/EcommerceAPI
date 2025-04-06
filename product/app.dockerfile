FROM golang:1.23-alpine3.20 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/rasadov/EcommerceAPI
COPY go.mod go.sum ./
RUN go mod download
COPY product product
COPY pkg pkg
RUN GO111MODULE=on go build -mod mod -o /go/bin/app ./product/cmd/product

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]