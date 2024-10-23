FROM golang:1.23-alpine AS binary_build
WORKDIR /app
COPY ./ ./
RUN go build -o /bin/notifier ./cmd/notifier

FROM alpine
COPY --from=binary_build /bin/notifier /bin/notifier
ENTRYPOINT ["/bin/notifier"]

