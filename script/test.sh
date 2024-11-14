#!/bin/bash

dir="generated/test/$(date +%s)"
msg="\033[0;32mOK\033[0m"
code=0
total_tests=0
passed_count=0
failed_count=0

assert() {
  expected="$1"
  input="$2"

  ./generated/test/modo run -o "$dir/out" --exec "$input"

  actual="$?"

  ((total_tests++))

  if [ "$actual" = "$expected" ]; then
    ((passed_count++))
    echo -e "\033[0;32m$input => $actual\033[0m"
  else
    ((failed_count++))
    echo -e "\033[0;31m$input => $expected expected, but got $actual\033[0m"
    msg="\033[0;31mNG\033[0m"
    code=-1
  fi
}

build-compiler(){
  mkdir -p "$dir"
  go build -o ./generated/test/modo ./cmd/modo 
  echo -e "\033[0;32mcompiled!\033[0m"
}

summary(){
  echo -e "\n------------------------"
  echo -e "summary: $msg, total: $total_tests, passed: $passed_count, failed: $failed_count"
  echo -e "------------------------"
}

testit(){
  assert 17 '17'
  assert 17 '(+ 4 13)'
  assert 6 '(+ 1 2 3)'
  assert 20 '(+ 1 2 3 4 10)'
  assert 35 '(+ 1 2 3 4 5 20)'
  assert 10 '(+ 1 2 (+ 3 4))'
  assert 10 '(+ (+ 1 2) (+ 3 4))'
  assert 21 '(+ (+ 1 2) (+ (+ 9 5) 4))'
  assert 39 '(+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8))'
}

build-compiler
testit
summary

exit $code

