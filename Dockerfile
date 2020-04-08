############################
# STEP 1 build executable binary
############################
FROM golang:latest AS builder
LABEL maintainer="art.frela@gmail.com" version="0.0.1"
ARG DRONE_BUILD_NUMBER
ENV DRONE_BUILD_NUMBER ${DRONE_BUILD_NUMBER}
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ADD ./go.* /go/src/chat/
ADD ./*.go /go/src/chat/
ADD ./domain/*.go /go/src/chat/domain/
ADD ./infra/*.go /go/src/chat/infra/
ADD ./config.yaml /go/src/chat/
ADD ./assets/* /go/src/chat/assets/
ADD ./.git /go/src/chat/.git
WORKDIR /go/src/chat/
RUN CGO_ENABLED=0 go build -ldflags="-X 'main.build=$DRONE_BUILD_NUMBER' -X 'main.githash=$(git rev-parse HEAD)'"
############################
# STEP 2 build a small image
############################
FROM alpine
RUN apk add --no-cache tzdata wget
ENV CHAT_HTTPD_PORT=8000
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY --from=builder /go/src/chat/chat /go/bin/chat/chat
COPY --from=builder /go/src/chat/config.yaml /go/bin/chat/config.yaml
COPY --from=builder /go/src/chat/assets/* /go/bin/chat/assets/
WORKDIR /go/bin/chat/
EXPOSE 8000
CMD ["./chat"]