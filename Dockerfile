FROM golang:1.15.6
ADD . /go/src/go-pyxy
# COPY . .
WORKDIR /go/src/go-pyxy
RUN go get gonum.org/v1/gonum
RUN go get github.com/rs/xid
RUN go get cloud.google.com/go/firestore
RUN go get cloud.google.com/go/storage
RUN go get google.golang.org/api/iterator
RUN go get google.golang.org/api/option
# RUN go install /go/src/go-pyxy
RUN go get go-pyxy
RUN go install
ENTRYPOINT ["/go/bin/go-pyxy"]
# RUN go run /go/src/go-pyxy/main.go
EXPOSE 8080