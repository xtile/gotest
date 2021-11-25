#TODO
# Set log path
# Set volumes

FROM golang:1.17.3-alpine3.14
WORKDIR /app
COPY . .
#RUN make 
RUN go build -v .
EXPOSE 8080
CMD ["gotest"]
