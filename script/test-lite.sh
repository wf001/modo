#!/bin/bash

. ./script/test.sh

testfile(){
  echo "== 1.modo ==="
  assertfile "1.modo" "ans 456\\\ndone\\\n"
  echo "== 2.modo ==="
  assertfile "2.modo" "Another\\\n"
  echo "== 3.modo ==="
  assertfile "3.modo" "45\\\n"
  echo "== 4.modo ==="
  assertfile "4.modo" "f\\\n"
  echo "== fizzbuzz.modo ==="
  assertfile "fizzbuzz.modo" "1\\\n2\\\nFizz\\\n4\\\nBuzz\\\nFizz\\\n7\\\n8\\\nFizz\\\nBuzz\\\n11\\\nFizz\\\n13\\\n14\\\nFizzBuzz\\\n16\\\n17\\\nFizz\\\n19\\\n"
  echo "== vector_scalar.modo ==="
  assertfile "vector_scalar.modo" "54\\\n[423, 83, 90]\\\n"

}

build-compiler
testfile
summary

exit $code
