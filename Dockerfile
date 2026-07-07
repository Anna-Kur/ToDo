FROM golang:1.25.4
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app ./main.go
EXPOSE 7540
CMD ["/my_app"]
