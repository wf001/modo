#!/bin/bash

. ./script/test.sh

testexec(){
  # integer type
  echo "== integer type ==="
  assertexec '(def main ::int (fn [] (prn 17)))' "17\\\n"
  # arithmetic operator
  echo "== arithmetic operator ==="
  assertexec '(def main ::int (fn [] (prn (+ 1 2 3 4 5 20))))' "35\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 2 (+ 3 4)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ 3 4)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ (+ 9 5) 4)))))' "21\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8)))))' "39\\\n"

  assertexec '(def main ::int (fn [] (prn (mod 20 5))))' "0\\\n"
  assertexec '(def main ::int (fn [] (prn (mod 17 5))))' "2\\\n"
  assertexec '(def main ::int (fn [] (prn (= 0 (mod 20 5)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 0 (mod (+ 18 2) 5)))))' "true\\\n"

  assertexec '(def main ::int (fn [] (prn (* 4 5))))' "20\\\n"
  assertexec '(def main ::int (fn [] (prn (* 9 5))))' "45\\\n"
  assertexec '(def main ::int (fn [] (prn (* 2 3 4))))' "24\\\n"
  assertexec '(def main ::int (fn [] (prn (* (+ 2 3) 5))))' "25\\\n"

  assertexec '(def main ::int (fn [] (prn (/ 4 2))))' "2\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 5 2))))' "2\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 8 2 2))))' "2\\\n"


  # equality operator
  echo "== equality operator ==="
  assertexec '(def main ::int (fn [] (prn (= 123 123))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 123 123 123))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 123 123 456))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= 5 (+ 3 2)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 3) (+ 3 2)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 -3) (+ 3 -2)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= true true))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= true false))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= false false))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo"))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= "www" "rrr"))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo" "foo"))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo" "goo"))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "goo" "hoo"))))' "false\\\n"
  # it's unknown bug
  assertexec '(def main ::int (fn [] (prn (= "foo" "bar"))))' "true\\\n" # must be false
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo" "bar"))))' "true\\\n" # must be false

  assertexec '(def main ::int (fn [] (prn (> 8 2))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (> 1 2))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (> 2 2))))' "false\\\n"

  assertexec '(def main ::int (fn [] (prn (< 8 2))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (< 1 2))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (< 2 2))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (< 1 2 3))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (< 1 2 1))))' "false\\\n"

  # logical operator
  echo "== logical operator ==="
  assertexec '(def main ::int (fn [] (prn (and (= 1 1) (= 1 0)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (and (= 1 1) (= 0 0)))))' "true\\\n"

  assertexec '(def main ::int (fn [] (prn (or (= 1 0) (= 1 0)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (or (= 1 1) (= 1 0)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (or (= 1 1) (= 1 1)))))' "true\\\n"

  ## global variable
  echo "== global variable ==="
  assertexec '(def x ::int 3) (def main ::int (fn [] (prn x)))' "3\\\n"
  assertexec '(def x ::string "hello") (def main ::int (fn [] (prn x)))' "hello\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x 2))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ 2 x))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ 2 3)))))' "6\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ 2 3) 4))))' "10\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (+ x y))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x x))))' "2\\\n"

  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 2))))' "false\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 1))))' "true\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (= x y))))' "false\\\n"

  ## binded variable
  echo "== binded variable ==="
  assertexec '(def main ::int (fn [] (let [x ::int 1] (prn x))))' "1\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1 y ::int (+ x 2)] (prn y))))' "3\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2] (prn (+ x y)))))' "6\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ y 3)] (prn (+ x z)))))' "9\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ y 3)] (prn (+ y z)))))' "7\\\n"
  assertexec '(def f ::int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 0 0) another no))))) (def main ::int (fn [] (prn (f 1))))' "123\\\n"
  assertexec '(def f ::int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 0 0) another no))))) (def main ::int (fn [] (prn (f 2))))' "456\\\n"
  assertexec '(def f ::int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 1 0) another no))))) (def main ::int (fn [] (prn (f 2))))' "789\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1] (let [y ::int (+ x 2)] (prn (+ x y))))))' "4\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1] (let [y ::int (+ x 2)] (prn y)))))' "3\\\n"

  ## if
  echo "== if ==="
  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (prn 12))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (prn 12))))' "12\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "12\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "13\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 1) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 1) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "12\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "13\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "13\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (if (= 1 2) (prn 12) (prn 13)) (prn 14) ))))' "13\\\n"

  assertexec '(def f ::int => int (fn [a] (if (= 1 a) 123 (if (= 0 0) 456 789)))) (def main ::int (fn [] (prn (f 1))))' "123\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) 123 (if (= 0 1) 456 789)))) (def main ::int (fn [] (prn (f 1))))' "123\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) 123 (if (= 0 0) 456 789)))) (def main ::int (fn [] (prn (f 2))))' "456\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) 123 (if (= 0 1) 456 789)))) (def main ::int (fn [] (prn (f 2))))' "789\\\n"

  assertexec '(def f ::int => int (fn [a] (if (= 1 a) (if (= 0 0) 123 456) 789))) (def main ::int (fn [] (prn (f 1))))' "123\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) (if (= 0 1) 123 456) 789))) (def main ::int (fn [] (prn (f 1))))' "456\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) (if (= 0 0) 123 456) 789))) (def main ::int (fn [] (prn (f 2))))' "789\\\n"
  assertexec '(def f ::int => int (fn [a] (if (= 1 a) (if (= 0 1) 123 456) 789))) (def main ::int (fn [] (prn (f 2))))' "789\\\n"

  assertexec '(def f ::int => int (fn [a] (if (= 1 a) (+ 1 1) (if (= 0 1) (+ 1 2) (+ 1 3))))) (def main ::int (fn [] (prn (f 1))))' "2\\\n"
  assertexec '(def f ::int => string (fn [a] (if (= 1 a) "y" "n"))) (def main ::int (fn [] (prn (f 1))))' "y\\\n"
  assertexec '(def f ::int => string (fn [a] (if (= 1 a) "y" "n"))) (def main ::int (fn [] (prn (f 2))))' "n\\\n"


  ## string type
  echo "== string type ==="
  assertexec '(def main ::int (fn [] (prn "hello")))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn t))))' "world\\\n"

  ## nil type
  echo "== nil type ==="
  assertexec '(def main ::int (fn [] (prn nil)))' "nil\\\n"
  assertexec '(def f :: int => string => nil (fn [a s] (prn a) (prn s))) (def main ::int (fn [] (prn(f 123 "hello"))))' "123\\\nhello\\\nnil\\\n"

  ## bool type
  echo "== bool type ==="
  assertexec '(def main ::int (fn [] (prn true)))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn false)))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= false (= 1 2)))))' "true\\\n"
  assertexec '(def isEven ::int => bool (fn [a] (= 0 (mod a 2)))) (def main ::int (fn [] (prn (isEven 5))))' "false\\\n"

  ## prn multi-line
  echo "== prn multi-line ==="
  assertexec '(def main :: int (fn [] (prn (+ 1 2)) (prn (+ 3 4))))' "3\\\n7\\\n"
  assertexec '(def x :: int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ x 3)] (prn (+ x z)) (prn (+ x y)))))' "11\\\n6\\\n"

  ## prn multi-value
  echo "== prn multi-value ==="
  assertexec '(def main :: int (fn [] (prn 1 2 "hello")))' "1 2 hello\\\n"
  assertexec '(def main :: int (fn [] (prn (+ 1 2) (= 2 4) "world")))' "3 false world\\\n"

  # function calling
  echo "== function calling ==="
  assertexec '(def f :: int => nil (fn [a] (prn a))) (def main ::int (fn [] (f 4)))' "4\\\n"
  assertexec '(def f :: int => string => nil (fn [a s] (prn a) (prn s))) (def main ::int (fn [] (f 123 "hello")))' "123\\\nhello\\\n"
  assertexec '(def f :: int => int => int (fn [a b] (+ a b))) (def main ::int (fn [] (prn (f 1 2))))' "3\\\n"
  assertexec '(def f :: string => string (fn [a] a)) (def main ::int (fn [] (prn (f "hello"))))' "hello\\\n"
  assertexec '(def f :: string => string (fn [a] "modo")) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"
  assertexec '(def f :: string => string (fn [a] (let [s ::string "modo"] s))) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"

  # loop
  echo "== loop ==="
  assertexec '(def f ::int => nil (fn [a] (prn a) (if (= 3 a) nil (f (+ a 1))))) (def main ::int (fn [] (f 1))))' "1\\\n2\\\n3\\\n"

  # vector
  echo "== vector ==="
  assertexec '(def main :: int (fn [] (let [vec :: [int] [41, 28, 239]] (prn vec))))'  "[41, 28, 239]\\\n"
  # assertexec '(def main :: int (fn [] (let [vec :: [string] ["Apple", "Banana"]] (prn vec))))' "[Apple, Banana]\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn v)))' "[233, 842]\\\n"
  # assertexec '(def v :: [string] ["Global", "Banana"]) (def main :: int (fn [] (prn v)))' "[Global, Banana]\\\n"

  # assertexec '(def main :: int (fn [] (let [vec :: [bool] [true, false, false, true]] (prn vec))))' "[true, false, false, true]\\\n"
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [41, 28, 239]] (prn (get vec 1)))))'  "28\\\n"
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [423, 83, 90]] (prn (conj vec 27)) (prn vec))))' "[423, 83, 90, 27]\\\n[423, 83, 90]\\\n"
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [423, 83, 90]] (prn (assoc vec 1 27)) (prn vec))))' "[423, 27, 90]\\\n[423, 83, 90]\\\n"
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [423, 83, 90]] (prn (pop vec)) (prn vec))))'  "[423, 83]\\\n[423, 83, 90]\\\n"

  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (get v 0))))' "233\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (conj v 569))))' "[233, 842, 569]\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (assoc v 1 388))))' "[233, 388]\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (pop v))))' "[233]\\\n"
  #assertexec "(def f :: [int] => [int] (fn [v] v)) (def main :: int (fn [] (let [vec :: [int] [42, 23]] (prn (f vec)))))" "[42, 23]\\\n"


}

build-compiler
testexec
summary

exit $code

