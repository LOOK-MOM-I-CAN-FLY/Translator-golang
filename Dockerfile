FROM golang:1.21-alpine AS builder

RUN apk add --no-cache openjdk17-jre curl

ENV ANTLR_VERSION=4.13.2
ENV ANTLR_JAR=/usr/local/lib/antlr-${ANTLR_VERSION}-complete.jar

RUN curl -L --fail -o ${ANTLR_JAR} https://www.antlr.org/download/antlr-${ANTLR_VERSION}-complete.jar

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN java -jar ${ANTLR_JAR} -Dlanguage=Go -visitor -no-listener -package parser -o parser SimpleLexer.g4 SimpleParser.g4

RUN CGO_ENABLED=0 GOOS=linux go build -o translator .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/translator ./translator

ENTRYPOINT ["./translator"]
