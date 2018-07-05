brew install pkg-config
brew install zmq
go get github.com/satori/go.uuid

we need to set up openssl include
export CGO_CFLAGS="-I/usr/local/opt/openssl/include"