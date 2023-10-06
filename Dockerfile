FROM golang:1.21.2-bookworm

WORKDIR /code

COPY . /code

RUN go install

RUN go build -o main .

CMD ["/code/main"]