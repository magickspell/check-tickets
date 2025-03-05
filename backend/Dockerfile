# build
FROM golang:1.24

WORKDIR /go-app-dir

COPY . .

RUN apt-get update && apt-get install -y bash

RUN go install github.com/air-verse/air@latest

RUN go mod tidy

# DEV
CMD ["air", "-c", ".air.toml"]

# PROD
# CMD ["go", "run", "main.go"]
# PROD с тестами
# CMD ["sh", "run.sh"]