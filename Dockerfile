FROM alpine:3.20

WORKDIR /app/
COPY ./messaging-system ./
COPY ./.config ./.config

RUN chmod +x messaging-system

CMD [ "./messaging-system" ]