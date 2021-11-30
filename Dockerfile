#TODO
# Set log path
# Set volumes

FROM golang:1.17.3-alpine3.14

RUN GOCACHE=OFF

#RUN go env -w GOPRIVATE=github.com/ereshzealous


WORKDIR /app
COPY . .


#RUN apk add git

#RUN git config --global url."https://xtile:89188e6ef3a334cc8d29bc857e6bf48a90dee192@github.com".insteadOf "https://github.com"

#RUN git clone https://github.com/xtile/gotest

#RUN make 
RUN go build -v ./...
EXPOSE 8080
CMD ["gotest"]






