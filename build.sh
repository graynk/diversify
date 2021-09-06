#!/usr/bin/env bash
cd reader
go build -o ../output/writer
cd ../writer
go build -o ../output/writer