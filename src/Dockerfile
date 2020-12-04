FROM golang
WORKDIR /src
COPY . .
RUN echo "test: ${address}"
RUN go build -o node
CMD ./node $address $type $link