FROM registry.pea.co.th/developer/vmsplus/api/backend/base:stable AS builder

USER root

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vmsplus main.go

FROM alpine:3.21.3

RUN apk --no-cache add ca-certificates

ARG APP_PATH=/app

USER 65532

WORKDIR $APP_PATH

COPY --chown=65532:65532 --from=builder /app/vmsplus . 

EXPOSE 8702

CMD ["sh", "-c", "exec ./vmsplus"]
