rm ./main
mv -f ./minutes.sqlite3 ./temp
go test -v *.go
rm minutes.sqlite3
mv -f ./temp/minutes.sqlite3 .
go build -v
