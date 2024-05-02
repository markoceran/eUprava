FROM golang:latest AS builder
WORKDIR /app
COPY ./sud_service/go.mod ./sud_service/go.sum ./
RUN go mod download
COPY ./sud_service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /root/
COPY /rbac_model/rbac_model.conf/ .
COPY /sud_service/policy.csv .
COPY --from=builder /app/main .
EXPOSE 8009
CMD ["./main"]