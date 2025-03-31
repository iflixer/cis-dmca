# build environment
FROM golang:1.22 as build-env
WORKDIR /server
COPY src/go.mod ./
RUN go mod download
COPY src src
WORKDIR /server/src
RUN CGO_ENABLED=0 GOOS=linux go build -o /server/build/httpserver .

FROM alpine:3.15
WORKDIR /app

COPY --from=build-env /server/build/httpserver /app/httpserver

#ENV GITHUB-SHA=<GITHUB-SHA>

ENTRYPOINT [ "/app/httpserver" ]
#ENTRYPOINT [ "ls", "-la", "/app/httpserver" ]
