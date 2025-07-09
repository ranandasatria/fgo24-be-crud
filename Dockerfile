FROM golang:alpine as build

WORKDIR /buildapp
COPY . .
RUN go build -o goapp main.go

FROM alpine:3.22

ARG APP_PORT
ARG APP_SECRET
ARG PGHOST
ARG PGPORT
ARG PGDATABASE
ARG PGUSER
ARG PGPASSWORD
ARG RDADDRESS
ARG RDPASSWORD
ARG RDDB

WORKDIR /app

COPY .env.example .env

RUN sed -i "s|^APP_PORT=.*|APP_PORT=${APP_PORT}|" .env && \
    sed -i "s|^APP_SECRET=.*|APP_SECRET=${APP_SECRET}|" .env && \
    sed -i "s|^PGHOST=.*|PGHOST=${PGHOST}|" .env && \
    sed -i "s|^PGPORT=.*|PGPORT=${PGPORT}|" .env && \
    sed -i "s|^PGDATABASE=.*|PGDATABASE=${PGDATABASE}|" .env && \
    sed -i "s|^PGUSER=.*|PGUSER=${PGUSER}|" .env && \
    sed -i "s|^PGPASSWORD=.*|PGPASSWORD=${PGPASSWORD}|" .env && \
    sed -i "s|^RDADDRESS=.*|RDADDRESS=${RDADDRESS}|" .env && \
    sed -i "s|^RDPASSWORD=.*|RDPASSWORD=${RDPASSWORD}|" .env && \
    sed -i "s|^RDDB=.*|RDDB=${RDDB}|" .env

COPY --from=build /buildapp/goapp /app/goapp

ENTRYPOINT ["/app/goapp"]
