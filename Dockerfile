FROM quay.io/projectquay/golang:1.22 AS build-env
WORKDIR /go/src/app
COPY main.go .
RUN CGO_ENABLED=0 go build -o /go/bin/app main.go
FROM gcr.io/distroless/static-debian12
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app"]
