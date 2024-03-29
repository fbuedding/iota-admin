FROM --platform=linux/amd64 golang:1.22-alpine
FROM cosmtrek/air:latest

WORKDIR /app

RUN apt-get update && apt-get install -y gcc 

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
