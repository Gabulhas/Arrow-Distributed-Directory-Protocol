FROM golang
WORKDIR /src
COPY . .
RUN echo "test: ${address}"
RUN go build -o node
CMD ./node --address=$address --type=$type --link=$link --requests=$requests --visualization=$VIS_ADDRESS