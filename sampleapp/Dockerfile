FROM golang:1.15

WORKDIR /usr/src/app

COPY . .

RUN cd sampleapp && go build

CMD ["sampleapp/./sampleapp"]
