# syntax=docker/dockerfile:1

FROM golang:1.19-buster

WORKDIR /app

# COPY ./objects ./
COPY . .
RUN go mod download

COPY *.go ./

# ENV DB_INFO="host=pgsql-dev user=postgres password=Rhd93 dbname=postgres sslmode=disable"
RUN go build -o /objects ./cmd/objects
RUN go build ./cmd/objects-cli

EXPOSE 8080

CMD [ "/objects" ]