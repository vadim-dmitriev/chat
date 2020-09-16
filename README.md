# Мессенджер

## Генерация кода, на основе proto-файлов
```
protoc proto/user.proto proto/token.proto --go_out=$GOPATH/src
protoc -Iproto proto/authService.proto --go_out=plugins=grpc:$GOPATH/src
```