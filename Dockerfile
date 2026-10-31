FROM golang:1.27 AS build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 go build -o /out/app ./cmd/app

FROM alpine:3.20
COPY --from=build /out/app /app
EXPOSE 8080
CMD ["/app"]