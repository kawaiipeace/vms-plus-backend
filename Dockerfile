FROM registry.pea.co.th/developer/vms-plus/api/backend/base:stable AS builder

USER root

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vms-plus main.go

FROM alpine:3.21.3

RUN apk --no-cache add ca-certificates

USER 65532

WORKDIR /app

COPY --chown=65532:65532 --from=builder /app/vms-plus . 

EXPOSE 8702

CMD ["sh", "-c", "exec ./vms-plus"]
