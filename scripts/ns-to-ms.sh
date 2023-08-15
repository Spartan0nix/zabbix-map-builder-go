#!/bin/bash

if [[ -z $1 ]]
then
    echo "missing required argument : <value-in-ns>"
    exit 1
fi

re='-h'
if [[ $1 =~ $re ]]
then
    echo "Convert value in nanoseconds to milliseconds"
    echo "Usage : $0 <value-in-ns>"
    exit 0
fi

echo $(( $1 / 1000000 ))

exit 0