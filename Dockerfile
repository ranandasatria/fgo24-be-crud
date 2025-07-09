FROM golang:alpine AS build
WORKDIR /buildapp
COPY . .
RUN go build -o goapp main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=build /buildapp/goapp /app/goapp
COPY .env.example /app/.env
RUN apk add --no-cache bash sed

ARG APP_PORT=""
ARG APP_SECRET=""
ARG PGHOST=""
ARG PGPORT=""
ARG PGDATABASE=""
ARG PGUSER=""
ARG PGPASSWORD=""
ARG RDADDRESS=""
ARG RDPASSWORD=""
ARG RDDB=""

ENV APP_PORT=${APP_PORT}
ENV APP_SECRET=${APP_SECRET}
ENV PGHOST=${PGHOST}
ENV PGPORT=${PGPORT}
ENV PGDATABASE=${PGDATABASE}
ENV PGUSER=${PGUSER}
ENV PGPASSWORD=${PGPASSWORD}
ENV RDADDRESS=${RDADDRESS}
ENV RDPASSWORD=${RDPASSWORD}
ENV RDDB=${RDDB}

RUN sed -i "s|^APP_PORT=.*|APP_PORT=$APP_PORT|" /app/.env && \
    sed -i "s|^APP_SECRET=.*|APP_SECRET=$APP_SECRET|" /app/.env && \
    sed -i "s|^PGHOST=.*|PGHOST=$PGHOST|" /app/.env && \
    sed -i "s|^PGPORT=.*|PGPORT=$PGPORT|" /app/.env && \
    sed -i "s|^PGDATABASE=.*|PGDATABASE=$PGDATABASE|" /app/.env && \
    sed -i "s|^PGUSER=.*|PGUSER=$PGUSER|" /app/.env && \
    sed -i "s|^PGPASSWORD=.*|PGPASSWORD=$PGPASSWORD|" /app/.env && \
    sed -i "s|^RDADDRESS=.*|RDADDRESS=$RDADDRESS|" /app/.env && \
    sed -i "s|^RDPASSWORD=.*|RDPASSWORD=$RDPASSWORD|" /app/.env && \
    sed -i "s|^RDDB=.*|RDDB=$RDDB|" /app/.env

ENTRYPOINT ["/app/goapp"]
