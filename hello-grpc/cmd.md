# CMD

gRPCにおけるコマンド

```
protoc greeter.proto --go_out=plugins=grpc:.

//http2のデバッグモード
export GODEBUG=http2debug=2
```

ssl 

```
brew install openssl

openssl version

openssl genrsa 2048 > private.key
openssl req -new -key private.key > server.csr
openssl x509 -days 367 -req -signkey private.key < server.csr > server.crt
```