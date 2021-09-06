#!/usr/bin/env bash
cd reader
go build -o ../output/reader
cd ../writer
go build -o ../output/writer