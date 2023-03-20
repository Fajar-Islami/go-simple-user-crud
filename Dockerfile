# Start by building the application.
FROM golang:1.19.2 as build
LABEL stage=dockerbuilder
WORKDIR /app
COPY . .

# Build the binary
RUN go build -o apps cmd/main.go

# Now copy it into our base image.
FROM alpine:3.9

# Copy bin file
WORKDIR /app
COPY --from=build /app/apps /app/apps
COPY docs .
RUN mkdir /logs

EXPOSE 8080
ENTRYPOINT ["/app/app"]