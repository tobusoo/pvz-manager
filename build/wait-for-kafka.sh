#!/bin/sh

HOST="kafka0"
PORT="29092"
cmd="$@"

echo 'Waiting for' $HOST $PORT
while :; do
  nc -w 1 -z $HOST $PORT -v
  if [ $? -eq 0 ]; then
    break;
  fi
done

echo 'exec' $cmd
exec $cmd
