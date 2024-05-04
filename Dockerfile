FROM golang:1.22rc2-alpine

RUN apk add --no-cache postgresql-dev gcc musl-dev

WORKDIR /app
COPY . .

RUN go build -o main .

CMD [".cmd/main"]