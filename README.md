<div align="center">
<picture align="center">
  <!-- ダークテーマ用 -->
  <source srcset="https://github.com/wf001/mkdirp.com/blob/dev/project/modo/static/logo_large_dark.png" media="(prefers-color-scheme: dark)" />
  
  <!-- ライトテーマ用（デフォルト） -->
  <img src="https://github.com/wf001/mkdirp.com/blob/dev/project/modo/static/logo_large_light.png" alt="Logo" />
</picture>
</div>

# modo, Programming language

**modo** is a statically typed, functional programming language inspired by Clojure and Haskell. 

Designed with simplicity and clarity in mind, it embraces the principles that,
- Explicit is better than implicit 
- there should be one — and preferably only one — obvious way to do it. 

modo uses LLVM as its backend, enabling high-performance native code generation.


## ✨ Features
- **Functional-first** language design
- **Static type system** inspired by Haskell
- **Minimal, clean syntax** inspired by Clojure
- **LLVM backend** for native code generation


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
(TBD)
 
## 🤝 Contributing
(TBD)

## 📄 License
modo is provided under the MIT License.
