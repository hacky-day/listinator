FROM node:24.6.0-bookworm-slim AS npm-build
WORKDIR /frontend
COPY frontend ./
RUN npm run build

FROM golang:1.24-bookworm AS go-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=npm-build /frontend/dist /app/frontend/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o listinator

FROM debian:bookworm-slim
COPY --from=go-build /app/listinator /
ENTRYPOINT ["/listinator"]
