#!/bin/bash

go build -o client_app

mkdir app1 app2 app3

mkdir app1/downloads app2/downloads app3/downloads

cp client_app app1

cp client_app app2

cp client_app app3