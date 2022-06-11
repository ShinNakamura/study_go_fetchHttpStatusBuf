#! /bin/bash
unalias -a

test -f tests/url_list.txt || exit 1
test -d ./bin || mkdir ./bin
go build -o bin/fetchHttpStatusBuf.exe
time <tests/url_list.txt ./bin/fetchHttpStatusBuf.exe
