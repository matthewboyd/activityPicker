from golang

run mkdir /app

ADD . /app

WORKDIR /app

RUN go build .

EXPOSE 8080

CMD ["/app/activityPicker"]