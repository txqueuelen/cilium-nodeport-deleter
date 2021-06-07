FROM alpine:latest

ARG OUTPUT_FOLDER=out
ARG BINARY_NAME=nodeport-deleter

WORKDIR /app

ADD ./$OUTPUT_FOLDER/$BINARY_NAME /app/nodeport-deleter

ENTRYPOINT [ "/app/nodeport-deleter" ]
CMD []
