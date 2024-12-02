FROM golang:alpine

WORKDIR /project/url-shortener/

COPY go.* ./
COPY .env ./  
COPY . .

RUN go build -o build/myapp .

EXPOSE 8080
ENTRYPOINT ["./build/myapp"]