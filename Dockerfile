FROM ubuntu:latest as build

WORKDIR /theadversary

COPY backend/ .

RUN apt update && \
    apt -y install wget gcc

RUN wget -O- https://go.dev/dl/go1.17.6.linux-amd64.tar.gz | tar -C /usr/local/ -xzf - && \
    ln -s /usr/local/go/bin/go /usr/bin/

RUN go build .

FROM ubuntu:latest

WORKDIR /theadversary

COPY frontend/ ./frontend/
COPY database.sqlite3 .
COPY .env .

COPY --from=build /theadversary/TheAdversary .

EXPOSE 8080

CMD ["./TheAdversary"]
