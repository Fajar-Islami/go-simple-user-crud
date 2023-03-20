# Start by building the application.
FROM golang:1.19.2 as build
LABEL stage=dockerbuilder
WORKDIR /app
COPY . .

# Make docs swagger
RUN make check-swag
RUN make docs

# Build the binary
RUN go build -o apps cmd/main.go

# Now copy it into our base image.
FROM alpine:3.9

# Copy bin file
WORKDIR /app
COPY --from=build /app/apps /app/apps
COPY --from=build /app/docs /app/docs
RUN mkdir /app/logs

EXPOSE 8080
ENTRYPOINT ["/app/apps"]