#!/bin/bash

res=` dmesg -LT|grep "killed as a result of limit"|awk -F ' ' '{print $8}'|awk -F '/' '{print $(NF-1)}'|cut -c 4- | awk '{a[$0]++}END{for(i in a){print i,a[i]}}' `

echo -e "${res}\n"
