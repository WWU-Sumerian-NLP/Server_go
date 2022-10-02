#Build Go
FROM golang:latest as builder 
ADD . /app 
WORKDIR /app/cmd
RUN go mod download 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .