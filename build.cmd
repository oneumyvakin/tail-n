set PKGNAME=github.com/oneumyvakin/tail-n
set LOCALPATH=%~dp0

go fmt %PKGNAME%
goimports.exe -w .

set GOOS=linux
set GOARCH=amd64
go build -o tail-n.%GOARCH% %PKGNAME%

set GOOS=windows
set GOARCH=amd64
go build -o tail-n.exe %PKGNAME%