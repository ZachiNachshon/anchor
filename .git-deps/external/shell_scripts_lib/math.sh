#!/bin/bash

increment() {
  local value=$1
  ((value++))
  echo "${value}"
}

subtract() {
  local first=$1
  local second=$2
  echo $((${first} - ${second}))
}

divide() {
  local value=$1
  local divide_by=$2
  echo $(awk -v dividend=${value} -v divisor=${divide_by} 'BEGIN { print  ( dividend / divisor ) }')
}

is_number() {
  local input=$1
  local regex='^[0-9]+$'
  local result=$input
  if ! [[ ${input} =~ ${regex} ]]; then
    result=""
  fi
  echo ${result}
}
