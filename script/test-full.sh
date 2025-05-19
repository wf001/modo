#!/bin/bash

. ./script/test.sh

testexec(){

  echo "==================="
  echo "== integer type ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn 17)))' "17\\\n"

  echo "==================="
  echo "== double value ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn 1.2))))' "1.2\\\n"

  echo "==================="
  echo "== string type ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn "hello")))' "hello\\\n"

  echo "==================="
  echo "== bool type ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn true)))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn false)))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= false (= 1 2)))))' "true\\\n"

  echo "==================="
  echo "== nil type ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn nil)))' "nil\\\n"

  echo "==================="
  echo "== arithmetic operator ==="
  echo "==================="
  # addition
  # int
  assertexec '(def main ::int (fn [] (prn (+ 15 20))))' "35\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 (+ 3 6)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ 3 4)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ (+ 9 5) 4)))))' "21\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 (+ (+ 4 5) 8)))))' "18\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (+ 1.2 3.3))))' "4.5\\\n"

  # subtraction
  # int
  assertexec '(def main ::int (fn [] (prn (- 15 20))))' "-5\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (- 3.1 2.3))))' "0.8\\\n"

  # multiplication
  # int
  assertexec '(def main ::int (fn [] (prn (* 4 5))))' "20\\\n"
  assertexec '(def main ::int (fn [] (prn (* 9 5))))' "45\\\n"
  assertexec '(def main ::int (fn [] (prn (* 3 4))))' "12\\\n"
  assertexec '(def main ::int (fn [] (prn (* (+ 2 3) 5))))' "25\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (* 3.1 2.3))))' "7.13\\\n"

  # division
  # int
  assertexec '(def main ::int (fn [] (prn (/ 4 2))))' "2\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 5 2))))' "2\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 5 0))))' "0\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 0 5))))' "0\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (/ 3.1 2.3))))' "1.34783\\\n"
  assertexec '(def main ::int (fn [] (prn (/ 1.0 10000000.0))))' "1e-07\\\n"

  # modulo
  # int
  assertexec '(def main ::int (fn [] (prn (mod 20 5))))' "0\\\n"
  assertexec '(def main ::int (fn [] (prn (mod 17 5))))' "2\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (mod 3.3 1.2))))' "0.9\\\n"
  assertexec '(def main ::int (fn [] (prn (mod 3.2 1.6))))' "0\\\n"



  echo "==================="
  echo "== comparison operator ==="
  echo "==================="
  # int
  assertexec '(def main ::int (fn [] (prn (= 123 123))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 123 456))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= 5 (+ 3 2)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 3) (+ 3 2)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 -3) (+ 3 -2)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 0 (mod 20 5)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 0 (mod (+ 18 2) 5)))))' "true\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (= 3.2 3.2))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= 3.1 3.2))))' "false\\\n"
  # bool
  assertexec '(def main ::int (fn [] (prn (= true true))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= true false))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (= false false))))' "true\\\n"
  # string
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo"))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (= "www" "rrr"))))' "false\\\n"
  # it's unknown bug
  assertexec '(def main ::int (fn [] (prn (= "foo" "bar"))))' "true\\\n" # must be false

  # int
  assertexec '(def main ::int (fn [] (prn (> 8 2))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (> 1 2))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (> 2 2))))' "false\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (> 3.2 3.1))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (> 3.2 3.3))))' "false\\\n"

  # int
  assertexec '(def main ::int (fn [] (prn (< 8 2))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (< 1 2))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (< 2 2))))' "false\\\n"
  # double
  assertexec '(def main ::int (fn [] (prn (< 3.2 3.1))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (< 3.2 3.3))))' "true\\\n"

  echo "==================="
  echo "== logical operator ==="
  echo "==================="
  assertexec '(def main ::int (fn [] (prn (and false false))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (and true false))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (and true true))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (and (= 1 1) (= 1 0)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (and (= 1 1) (= 0 0)))))' "true\\\n"

  assertexec '(def main ::int (fn [] (prn (or false false))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (or true false))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (or true true))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (or (= 1 0) (= 1 0)))))' "false\\\n"
  assertexec '(def main ::int (fn [] (prn (or (= 1 1) (= 1 0)))))' "true\\\n"
  assertexec '(def main ::int (fn [] (prn (or (= 1 1) (= 1 1)))))' "true\\\n"


  echo "==================="
  echo "== global variable ==="
  echo "==================="
  # int
  assertexec '(def x ::int 3) (def main ::int (fn [] (prn x)))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x 2))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ 2 x))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ 2 3)))))' "6\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ (+ 2 3) 4)))))' "10\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (+ x y))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x x))))' "2\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 2))))' "false\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 1))))' "true\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (= x y))))' "false\\\n"

  # double
  assertexec '(def f :: double 1.32) (def main ::int (fn [] (prn f)))' "1.32\\\n"

  # string
  assertexec '(def x ::string "hello") (def main ::int (fn [] (prn x)))' "hello\\\n"

  echo "==================="
  echo "== binded variable ==="
  echo "==================="
  # int
  assertexec '(def main ::int (fn [] (let [x ::int 1] (prn x))))' "1\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1 y ::int (+ x 2)] (prn y))))' "3\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1] (let [y ::int (+ x 2)] (prn (+ x y))))))' "4\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1] (let [y ::int (+ x 2)] (prn y)))))' "3\\\n"
  # double
  assertexec '(def f :: double (fn [] (let [fl ::double 1.42] fl))) (def main ::int (fn [] (prn f)))' "1.42\\\n"
  assertexec '(def f :: double (fn [] (let [fl ::double 1.42 fm ::double (+ fl 3.2)] fm))) (def main ::int (fn [] (prn f)))' "4.62\\\n"
  # string
  assertexec '(def main ::int (fn [] (let [s ::string "hello"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn t))))' "world\\\n"

  echo "==================="
  echo "== global and binded variable ==="
  echo "==================="
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2] (prn (+ x y)))))' "6\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ y 3)] (prn (+ x z)))))' "9\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ y 3)] (prn (+ y z)))))' "7\\\n"

  echo "==================="
  echo "== if ==="
  echo "==================="
  # int
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

  # string
  assertexec '(def main ::int (fn [] (let [s ::string "hello"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn t))))' "world\\\n"

  assertexec '(def f ::int => string (fn [a] (if (= 1 a) "y" "n"))) (def main ::int (fn [] (prn (f 1))))' "y\\\n"
  assertexec '(def f ::int => string (fn [a] (if (= 1 a) "y" "n"))) (def main ::int (fn [] (prn (f 2))))' "n\\\n"


  echo "==================="
  echo "== function calling ==="
  echo "==================="
  # int
  assertexec '(def f :: int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 0 0) another no))))) (def main ::int (fn [] (prn (f 1))))' "123\\\n"
  assertexec '(def f :: int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 0 0) another no))))) (def main ::int (fn [] (prn (f 2))))' "456\\\n"
  assertexec '(def f :: int => int (fn [a] (let [yes ::int 123 another ::int 456 no ::int 789] (if (= 2 (+ 1 a)) yes (if (= 1 0) another no))))) (def main ::int (fn [] (prn (f 2))))' "789\\\n"
  assertexec '(def f :: int => nil (fn [a] (prn a))) (def main ::int (fn [] (f 4)))' "4\\\n"
  assertexec '(def f :: int => int => int (fn [a b] (+ a b))) (def main ::int (fn [] (prn (f 1 2))))' "3\\\n"
  # double
  assertexec '(def f :: double => nil (fn [a] (prn a))) (def main ::int (fn [] (f 1.23)))' "1.23\\\n"
  # string
  assertexec '(def f :: string => string (fn [a] a)) (def main ::int (fn [] (prn (f "hello"))))' "hello\\\n"
  assertexec '(def f :: string => string (fn [a] "modo")) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"
  assertexec '(def f :: string => string (fn [a] (let [s ::string "modo"] s))) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"
  assertexec '(def f :: int => string => nil (fn [a s] (prn a) (prn s))) (def main ::int (fn [] (f 123 "hello")))' "123\\\nhello\\\n"
  # bool
  assertexec '(def isEven ::int => bool (fn [a] (= 0 (mod a 2)))) (def main ::int (fn [] (prn (isEven 5))))' "false\\\n"
  # nil
  assertexec '(def f :: int => string => nil (fn [a s] (prn a) (prn s))) (def main ::int (fn [] (prn(f 123 "hello"))))' "123\\\nhello\\\nnil\\\n"

  echo "==================="
  echo "== loop ==="
  echo "==================="
  assertexec '(def f ::int => nil (fn [a] (prn a) (if (= 3 a) nil (f (+ a 1))))) (def main ::int (fn [] (f 1))))' "1\\\n2\\\n3\\\n"

  echo "==================="
  echo "== vector ==="
  echo "==================="
  # local int
  assertexec '(def main :: int (fn [] (let [vec :: [int] [41 28 239]] (prn vec))))'  "[41 28 239]\\\n"
  # local string
  assertexec '(def main :: int (fn [] (let [vec :: [string] ["Apple" "Banana"]] (prn vec))))' "[Apple Banana]\\\n"
  # local bool
  assertexec '(def main :: int (fn [] (let [vec :: [bool] [true false false true]] (prn vec))))' "[true false false true]\\\n"
  # local double
  assertexec '(def main :: int (fn [] (let [vec :: [double] [3.3 9.3]] (prn vec))))' "[3.3 9.3]\\\n"
  # global int
  assertexec '(def v :: [int] [233 842]) (def main :: int (fn [] (prn v)))' "[233 842]\\\n"
  # global string
  assertexec '(def v :: [string] ["Global" "Banana"]) (def main :: int (fn [] (prn v)))' "[Global Banana]\\\n"
  # global bool
  assertexec '(def v :: [bool] [false true true]) (def main :: int (fn [] (prn v)))' "[false true true]\\\n"

  # taking args and returing
  assertexec "(def f::[int] => nil (fn[v] (prn v))) (def main::int (fn[] (let[v::[int][42 64 90]] (f v))))" "[42 64 90]\\\n"
  assertexec "(def f::[int] (fn[] (let[v::[int] [42 64 90 84]]v))) (def main::int (fn[] (prn f)))" "[42 64 90 84]\\\n"
  assertexec "(def f :: [int] => [int] (fn [v] (conj v 81))) (def main :: int (fn [] (let [v :: [int] [42 64 90]] (prn (f v)))))" "[42 64 90 81]\\\n"

  # 3-dimension
  assertexec "(def main :: int (fn [] (let [vec :: [[[int]]] [[[11 12 13] [14 15]] [[21 22 23] [24]] [[31 32] [33 34 35]] [[41 42] [43 44 45]]]] (prn (nth (nth (nth vec 1) 1) 0)) )))" "24\\\n"

  # reference
  assertexec "(def main :: int (fn [] (let [vec1 :: [int] [42 23] vec2 :: [int] vec1] (prn (conj vec2 1)) (prn (conj vec1 3)) (prn vec1) (prn vec2))))" "[42 23 1]\\\n[42 23 3]\\\n[42 23]\\\n[42 23]\\\n"


  echo "==================="
  echo "== struct ==="
  echo "==================="
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool})(def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true}] (prn (get node :name)) (prn (get node :age)) (prn (get node :isMale))))))' "richard\\\n20\\\ntrue\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double})(def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true :height 5.8}] (prn (get node :name)) (prn (get node :age)) (prn (get node :isMale)) (prn (get node :height))))))'  "richard\\\n20\\\ntrue\\\n5.8\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool}) (def f :: Person (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true}] node))) (def main :: int (fn [] (prn (get f :age))))' "20\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true :height 5.8 :v [4 3 2]}] (prn (get node :name)) (prn (get node :age)) (prn (get node :isMale)) (prn (get node :height)) (prn (get node :v)) )))' "richard\\\n20\\\ntrue\\\n5.8\\\n[4 3 2]\\\n"
  assertexec '(defschema Person {:name ::string :next ::Person})(def main ::int (fn [] (let [node ::Person {:name "richard" :next {:name "rebecca" :next {:name "andre" :next {:name "bob"}}}}] (prn (get node :name)) (prn (get (get node :next) :name)) (prn (get (get (get node :next) :next) :name)) (prn (get (get (get (get node :next) :next) :next) :name)))))' "richard\\\nrebecca\\\nandre\\\nbob\\\n"
  assertexec '(defschema Country {:name ::string})(defschema Person {:name ::string :country ::Country})(def main ::int (fn [] (let [c ::Country {:name "deutsch"} person ::Person {:name "mike" :country c}] (prn (get (get person :country) :name)))))' "deutsch\\\n"
  # null access
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:name "richard" :isMale true :height 5.8 :v [4 3 2]}] (prn (get node :age)))))' "0\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:age 20 :isMale true :height 5.8 :v [4 3 2]}] (prn (get node :name)))))' "nil\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :height 5.8 :v [4 3 2]}] (prn (get node :isMale)))))' "false\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true :v [4 3 2]}] (prn (get node :height)))))' "0\\\n"
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool :height :: double :v :: [int]}) (def main :: int (fn [] (let [node :: Person {:age 20 :name "richard" :isMale true :height 5.8 }] (prn (get node :v)))))' "nil\\\n"
  assertexec '(defschema Country {:name ::string}) (defschema Person {:name ::string :country ::Country}) (def main ::int (fn [] (let [co ::Country {:name "deutsch"} person ::Person {:name "jimmy"}] (prn (get (get person :country) :name)))))' "nil\\\n"

  #############
  # Prelude functions
  #############

  echo "==================="
  echo "== prn ==="
  echo "==================="
  # multi-line
  assertexec '(def main :: int (fn [] (prn (+ 1 2)) (prn (+ 3 4))))' "3\\\n7\\\n"
  assertexec '(def x :: int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ x 3)] (prn (+ x z)) (prn (+ x y)))))' "11\\\n6\\\n"

  # multi-value
  assertexec '(def main :: int (fn [] (prn 1 2 "hello")))' "1 2 hello\\\n"
  assertexec '(def main :: int (fn [] (prn (+ 1 2) (= 2 4) "world")))' "3 false world\\\n"

  echo "==================="
  echo "== map ==="
  echo "==================="
  assertexec "(def inc :: int => int (fn [x] (+ x 1))) (def main :: int (fn [] (prn (map inc [32 13 99]))))" "[33 14 100]\\\n"
  assertexec "(def inc :: double => double (fn [x] (+ x 1.3))) (def main :: int (fn [] (prn (map inc [4.2 1.8]))))" "[5.5 3.1]\\\n"
  assertexec '(def str :: int => string (fn [x] "a")) (def main :: int (fn [] (prn (map str [32 13 99]))))' "[a a a]\\\n"
  assertexec '(def truely :: int => bool (fn [x] true)) (def main :: int (fn [] (prn (map truely [32 13 99]))))' "[true true true]\\\n"

  echo "==================="
  echo "== filter ==="
  echo "==================="
  assertexec '(def even :: int => bool (fn [x] (= 0 (mod x 2))))  (def main :: int (fn [] (prn (filter even [8 1 7 10 18 42 45]))))' "[8 10 18 42]\\\n"
  assertexec '(def isThree :: int => bool (fn [x] (= 3 x))) (def main :: int (fn [] (prn (filter isThree [8 1 7 10 3 42 3]))))' "[3 3]\\\n"
  assertexec '(def isThreeOne :: double => bool (fn [x] (= 3.1 x))) (def main :: int (fn [] (prn (filter isThreeOne [1.1 2.1 3.1]))))' "[3.1]\\\n"

  echo "==================="
  echo "== reduce ==="
  echo "==================="
  assertexec '(def f :: int => int => int (fn [x y] (+ x y))) (def main :: int (fn [] (prn (reduce f 0 [3 2 1 9]))))' "15\\\n"
  assertexec '(def f :: double => double => double (fn [x y] (+ x y))) (def main :: int (fn [] (prn (reduce f 0.1 [1.2 2.2 8.1 3.2])))))' "14.8\\\n"

  echo "==================="
  echo "== nth ==="
  echo "==================="
  # nth
  assertexec '(def main :: int (fn [] (let [vec :: [int] [41 28 239]] (prn (nth vec 1)))))'  "28\\\n"
  assertexec '(def v :: [int] [233 842]) (def main :: int (fn [] (prn (nth v 0))))' "233\\\n"
  assertexec '(def main :: int (fn [] (let [vec :: [int] [423 83 90]] (prn (nth vec -1)) (prn (nth vec 2)) (prn (nth vec 3)) )))' "nil\\\n90\\\nnil\\\n"

  echo "==================="
  echo "== conj ==="
  echo "==================="
  assertexec '(def main :: int (fn [] (let [vec :: [int] [423 83 90]] (prn (conj vec 27)) (prn vec))))' "[423 83 90 27]\\\n[423 83 90]\\\n"
  assertexec '(def v :: [int] [233 842]) (def main :: int (fn [] (prn (conj v 569))))' "[233 842 569]\\\n"
  assertexec "(def v :: [int] [233 842]) (def f :: [int] (fn [] (conj v 39))) (def main :: int (fn [] (prn f)))" "[233 842 39]\\\n"
  assertexec '(def f::[int] => [int] (fn[v] (conj v 78))) (def main ::int (fn[] (let [v ::[int] [42 64 90]] (prn (f v)))))' "[42 64 90 78]\\\n"
  assertexec "(def main :: int (fn [] (let [vec :: [[int]] [[12 34 5 6] [789]]] (prn (nth (conj vec [11 12]) 2)) (prn (nth vec 2)))))" "[11 12]\\\nnil\\\n"

  # echo "==================="
  # echo "== assoc ==="
  # echo "==================="
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [423, 83, 90]] (prn (assoc vec 1 27)) (prn vec))))' "[423, 27, 90]\\\n[423, 83, 90]\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (assoc v 1 388))))' "[233, 388]\\\n"

  # echo "==================="
  # echo "== pop ==="
  # echo "==================="
  # assertexec '(def main :: int (fn [] (let [vec :: [int] [423, 83, 90]] (prn (pop vec)) (prn vec))))'  "[423, 83]\\\n[423, 83, 90]\\\n"
  # assertexec '(def v :: [int] [233, 842]) (def main :: int (fn [] (prn (pop v))))' "[233]\\\n"
  # map

  echo "==================="
  echo "== internal ==="
  echo "==================="
  # error message
  assertexec '(def main ::int (fn [] (let [x ::int 1] (let [y ::int (+ x 2)] (prn z)))))' "level=error msg=\"syntax error: undefined: z\""
  assertexec '(def main ::int (fn [] (let [x ::int "hello"] (prn x))))' "level=error msg=\"syntax error: cannot use string type as x (int type)\""
  assertexec '(def main ::int (fn [] (let [x ::bool "true"] (prn x))))' "level=error msg=\"syntax error: cannot use string type as x (bool type)\""
  assertexec '(defschema Person {:name :: string :age :: int :isMale :: bool})(def main :: int (fn [] (let [node :: node {:age 20 :name "richard" :isMale true}] (prn (get node :name)) (prn (get node :age)) (prn (get node :isMale))))))' "level=error msg=\"syntax error: undefined type: node\""
  assertexec '(def main ::int (fn [] (let [x ::int 1 x ::int 2] (prn x))))' "level=error msg=\"syntax error: cannot use x, already used\""
}

build-compiler
testexec
summary

exit $code

