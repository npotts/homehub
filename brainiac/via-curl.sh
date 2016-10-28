#!/bin/bash


curl -sw "Return Code: %{http_code}" -u brainiac:brainiac -X "PUT" -o curl.data --data-raw '{"table":"a", "data": {"field": 1.0}}' http://localhost:8080
echo 



#insert 100 lines
for i in {1..1000}; do
	curl -sw "Return Code: %{http_code}" -X "POST"  -u brainiac:brainiac -o curl.data --data-raw '{"table":"a", "data": {"field": 1.0}}' http://localhost:8080
	echo
done

curl -sw "Return Code: %{http_code}" -X "GET"  -u brainiac:brainiac -o curl.data --data-raw '{"table":"a", "data": {"field": 1.0}}' http://localhost:8080
echo
cat  curl.data