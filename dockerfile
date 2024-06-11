FROM golang:1.22 as build

RUN mkdir /app

WORKDIR /app

COPY . . 

RUN CGO_ENABLED=0 go build -o TelegramBot . 

FROM alpine

RUN mkdir /app

WORKDIR /app

COPY --from=build /app/TelegramBot  /app

CMD [ "/app/TelegramBot" ]