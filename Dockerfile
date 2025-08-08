FROM alpine:latest

WORKDIR /app
COPY bin/main .

RUN chmod +x main

EXPOSE 8080

ENTRYPOINT ["./main"]
