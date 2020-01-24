#!/bin/bash

res=`dmesg -LT|grep "killed as a result of limit"|awk -F '/' '{print $NF}'|cut -c 4-| awk '{a[$0]++}END{for(i in a){print i,a[i]}}' `

echo -e "${res}\n"
