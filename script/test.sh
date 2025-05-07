#!/bin/bash

dir="generated/test/$(date +%s)"
testdatadir="./script/testdata"
msg="\033[0;32mOK\033[0m"
code=0
total_tests=0
passed_count=0
failed_count=0

assertexec() {
  input="$1"
  expected="$2"
  assert "$input" "$expected" 
}

assertfile() {
  input=$(cat "$testdatadir/$1")
  expected="$2"
  assert "$input" "$expected" 
}

assert() {
  input="$1"
  expected="$2"
  actual_output=$(./generated/test/modo run --exec "$input" |gsed ':a;N;$!ba;s/\n/\\\\n/g')
  actual_exit_code="$?"

  ((total_tests++))

  diff_result=$(diff <(echo "$expected") <(echo "$actual_output"))
  if [ "$actual_exit_code" -eq 0 ] && [ -z "$diff_result" ]; then
    ((passed_count++))
    echo -e "$input => $actual_output \033[0;32mOK\033[0m"
  else
    ((failed_count++))
    echo -e "$input \033[0;31mNG\033[0m\\nwant: $expected\\nhave: $actual_output "
    msg="\033[0;31mNG\033[0m"
    code=1
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

