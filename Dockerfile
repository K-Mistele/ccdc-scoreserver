FROM golang:1.15.8
WORKDIR /app
COPY . .
EXPOSE 8080
CMD ["go", "get"]
ENTRYPOINT ["go", "run", "server.go"]