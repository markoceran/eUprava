# Stage 1: Build the Go application
FROM golang:latest AS builder
WORKDIR /app

# Copy go.mod and go.sum to the workspace
COPY ./tuzilastvo_service/go.mod ./tuzilastvo_service/go.sum ./
# Download dependencies
RUN go mod download

# Copy the source code
COPY ./tuzilastvo_service/ .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create a lightweight container to run the application
FROM alpine:latest
WORKDIR /root/

# Install tzdata package in the final image
RUN apk update && apk add --no-cache tzdata

# Set the timezone environment variable
ENV TZ=Europe/Belgrade

# Copy necessary files from the build stage
COPY /rbac_model/rbac_model.conf ./rbac_model.conf
COPY /tuzilastvo_service/policy.csv ./policy.csv
COPY --from=builder /app/main .

# Expose the required port
EXPOSE 8009

# Run the application
CMD ["./main"]
