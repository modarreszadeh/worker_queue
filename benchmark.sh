#!/bin/bash

url="http://127.0.0.1:8000/work"

data='{"title":"Task","message":"This is benchmark task","delay":1}'

concurrency=1000

total_requests=1000000

tmpfile=$(mktemp)

echo "$data" > $tmpfile

ab -n $total_requests -c $concurrency -T "application/json" -p $tmpfile -k $url

rm $tmpfile

