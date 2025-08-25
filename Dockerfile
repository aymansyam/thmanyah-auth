# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src

# Prime module cache (ok even if go.sum doesn't exist yet)
COPY go.mod ./
RUN go mod download

# Copy source and finalize deps (writes go.sum)
COPY . .
RUN go mod tidy

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /auth

# ---- minimal runtime ----
FROM gcr.io/distroless/static:nonroot
ENV PORT=8081
EXPOSE 8081
USER nonroot:nonroot
COPY --from=build /auth /auth
ENTRYPOINT ["/auth"]
