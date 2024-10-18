FROM golang:1.23-alpine AS binary_build
WORKDIR /app
COPY ./ ./
RUN go build -o /bin/notifier ./cmd/notifier

FROM alpine
COPY --from=binary_build /bin/notifier /bin/notifier
COPY --from=binary_build /app/configs/notifier.yaml configs/notifier.yaml
ENTRYPOINT ["/bin/notifier"]

