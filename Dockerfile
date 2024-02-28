FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . . 
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/bracket.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bracket /bracket
COPY ./templates /templates

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bracket"]
