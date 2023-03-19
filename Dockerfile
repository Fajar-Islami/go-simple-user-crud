# Start by building the application.
FROM golang:1.19.2 as build
LABEL stage=dockerbuilder
WORKDIR /app/example
COPY . .

# Build the binary
RUN go build -o app cmd/main.go

# Now copy it into our base image.
FROM alpine:3.9

# Copy bin file
COPY --from=build /app/example/app /app/app
COPY docs .
RUN mkdir /logs

EXPOSE 8000
ENTRYPOINT ["/app/app"]