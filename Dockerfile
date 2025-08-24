FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /auth

FROM alpine:3.20
ENV PORT=8081 JWT_SECRET=dev-secret
EXPOSE 8081
COPY --from=build /auth /auth
USER 65532:65532
ENTRYPOINT ["/auth"]
