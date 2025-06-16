# syntax=docker/dockerfile:1

FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# what about images and templ files
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /frontend ./cmd/frontend

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /frontend /frontend

EXPOSE 4001

ENTRYPOINT ["/frontend"]
