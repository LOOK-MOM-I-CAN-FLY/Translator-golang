
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache openjdk11-jre curl

ENV ANTLR_VERSION=4.13.2
RUN curl -LO https://www.antlr.org{ANTLR_VERSION}-complete.jar && \
    mv antlr-${ANTLR_VERSION}-complete.jar /usr/local/lib/antlr4.jar

ENV CLASSPATH=".:/usr/local/lib/antlr4.jar:$CLASSPATH"
alias antlr4='java -jar /usr/local/lib/antlr4.jar'

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN java -jar /usr/local/lib/antlr4.jar -Dlanguage=Go -o parser MyGrammar.g4

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]
