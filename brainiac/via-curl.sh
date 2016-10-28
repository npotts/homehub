#!/bin/bash


curl --digest -u brainiac:brainiac -sw "Return Code: %{http_code}"  -X "PUT" -o curl.data --data-raw '{"a":"table", "data": {"field": 1.0}}' http://localhost:8080
echo 



#insert 100 lines
for i in {1..10}; do
	curl --digest -u brainiac:brainiac -sw "Return Code: %{http_code}" -X "POST"  -o curl.data --data-raw '{"a":"table", "data": {"field": 1.0}}' http://localhost:8080
	echo
done