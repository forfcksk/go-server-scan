#!/bin/bash
input=$1
port=$2
try=8
while IFS= read -r line
do
  echo "$line"
  open=$(echo >/dev/tcp/$line/$port)
    if [ "$line" == "" ]
        then
            curl -X POST -H "Content-Type: application/json" -d '{"ID":"'$try'","IP":"'$line'","Port":"'$port'"}' "http://127.0.0.1:8000/socks";
        else
            echo "Port is closed"
    fi
    try=$try+1
done < "$input"
