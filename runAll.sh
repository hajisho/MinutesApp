npx gulp
mv -f minutes.sqlite3 ./temp
go test -v *.go
mv -f ./temp/minutes.sqlite3 .
go build -x
./main
