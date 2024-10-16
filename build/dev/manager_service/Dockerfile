FROM golang:1.23-alpine AS binary_build
WORKDIR /app
COPY ./ ./
RUN go build -o /bin/manager_service ./cmd/manager_service/main.go

FROM alpine
COPY --from=binary_build /bin/manager_service /bin/manager_service
COPY --from=binary_build /app/pkg/manager-service/v1/manager-service.swagger.json pkg/manager-service/v1/manager-service.swagger.json
ENTRYPOINT ["/bin/manager_service"]

