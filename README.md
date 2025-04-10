# modo

**modo**(call /ˈmoʊ.doʊ/, same as mode) is a statically typed, functional programming language inspired by **Clojure** and **Haskell**, designed with simplicity and clarity in mind. It uses LLVM as a backend, allowing high-performance native code generation.

---

## ✨ Features
- **Functional-first** language design
- **Static type system** inspired by Haskell
- **Minimal, clean syntax** inspired by Clojure
- **LLVM backend** for native code generation

---

## 📦 Example: FizzBuzz in modo

``` clojure
(def fizzbuzz :: int => int => nil
  (fn [n max] 
    (let [fizz ::int 3
          buzz ::int 5]
      (if (= n max)
        nil
        (if (= 0 (mod n (* fizz buzz)))
          (prn "FizzBuzz")
          (if (= 0 (mod n fizz))
            (prn "Fizz")
            (if (= 0 (mod n buzz))
              (prn "Buzz")
              (prn n))))))
      (fizzbuzz (+ 1 n) max))))

(def main :: int
  (fn []
    (fizzbuzz 1 20)))

```

## 🛠️ Installation
To install modo, follow these steps:
(TBD)

## 📚 Documentation
For detailed documentation, visit here.

## 🚧 TODO
 - [ ] Support for Float type
 - [ ] Support for Vectors (flexible-length arrays)
 - [ ] Support for Structs (custom compound types)
 
## 🤝 Contributing
(TBD)

## 📄 License
modo is provided under the MIT License.
