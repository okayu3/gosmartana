set GOOS=windows
set GOARCH=386
go build -o smartana.exe ./cmd/smartana/main.go
go build -o smartanalogic.exe ./cmd/smartanalogic/main.go
set GOOS=darwin
set GOARCH=386
go build -o smartana ./cmd/smartana/main.go
go build -o smartanalogic ./cmd/smartanalogic/main.go
