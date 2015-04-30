build instructions

    brew install libsodium
    brew install --HEAD --with-libsodium zmq
    go get github.com/pebbe/zmq4
    go build sub.go
    go build pub.go
