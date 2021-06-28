### build
FROM golang:1.16.5-alpine as build
LABEL MAINTAINER "Luis Rohten <luisaugustomelo@gmail.com>"

ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/klever
RUN cd /go/src/klever && go build .


FROM 1.16.5-alpine
WORKDIR /app
COPY --from=build /go/src/klever/klever /app
COPY .env /app

EXPOSE 9000

ENTRYPOINT [ "./klever" ]
# WORKDIR /app
# COPY . /app
# COPY .env /app
# RUN go get -u -v -f all

# EXPOSE 9000
# CMD ["go", "run", "./server/server.go"]
