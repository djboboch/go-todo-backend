FROM golang:alpine AS build-env

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main main.go

WORKDIR /dist

RUN cp /build/main .

FROM alpine

RUN mkdir /app
COPY --from=build-env /dist/main /app
WORKDIR /app
ENTRYPOINT ["./main"]

EXPOSE 8000