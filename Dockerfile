FROM golang:1.14.3-alpine3.11 as build

# prepare build dir
RUN mkdir /app
WORKDIR /app

COPY . .

RUN go build

# prepare release image
FROM alpine:3.11.6 AS app

RUN mkdir /app
WORKDIR /app

COPY --from=build /app/star-wars-planets ./
RUN chown -R nobody: /app
USER nobody

ARG MONGODB_HOST
ARG MONGODB_DATABASE
ARG REDIS_HOST
ARG REDIS_NETWORK
ARG PORT

ENV MONGODB_HOST=$MONGODB_HOST
ENV MONGODB_DATABASE=$MONGODB_DATABASE
ENV REDIS_HOST=$REDIS_HOST
ENV REDIS_NETWORK=$REDIS_NETWORK
ENV PORT=$PORT

ENV HOME=/app

CMD ./star-wars-planets