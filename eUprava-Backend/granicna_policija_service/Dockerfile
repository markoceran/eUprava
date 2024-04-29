FROM golang:latest AS builder
WORKDIR /app
COPY ./granicna_policija_service/go.mod ./granicna_policija_service/go.sum ./
RUN go mod download
COPY ./granicna_policija_service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /root/
COPY /rbac_model/rbac_model.conf/ .
COPY /granicna_policija_service/policy.csv .
COPY --from=builder /app/main .
EXPOSE 8009

CMD ["./main"]