FROM alpine:3.19.1

WORKDIR /app

ARG EXECUTABLE_NAME

COPY ./${EXECUTABLE_NAME} /app/oma

RUN chmod +x /app/oma

CMD ["/app/oma"]