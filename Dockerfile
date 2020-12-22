FROM golang:1.15.5-alpine3.12

COPY . /app
RUN mkdir /go/dist \
&& cd /app \
&& go build

ENTRYPOINT [ "/app/htn2hugo" ]