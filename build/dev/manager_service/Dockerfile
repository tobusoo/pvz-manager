FROM golang:1.23-alpine AS binary_build
WORKDIR /app
COPY ./ ./
RUN go build -o /bin/manager_service ./cmd/manager_service

FROM alpine
COPY --from=binary_build /bin/manager_service /bin/manager_service
ENTRYPOINT ["/bin/manager_service"]

