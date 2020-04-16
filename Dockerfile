FROM golang:1.14

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]

EXPOSE 8080

#build with command: docker build -t visense .
#run server with command: docker run -P -it --rm --name visenserunning visense:latest