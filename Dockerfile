FROM debian:latest

RUN apt update -y
RUN apt install -y libc6-dev

WORKDIR /app

COPY vanieth ./

CMD ./vanieth
