FROM golang

RUN go version

WORKDIR /app

COPY . .

#install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

#run waiting for postgres
RUN chmod +x wait-for-postgres.sh

#build my bot
RUN go build -o telega ./main.go

CMD ["./telega"]