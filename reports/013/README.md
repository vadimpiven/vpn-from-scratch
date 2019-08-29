# Web services development on Go part 2, week 4 (2019/08/26)

## Configs & Monitoring
- `flag` package could be used to parse command line flags and arguments, first defint all required flags with `flag.Bool("comments", false, "Enable comments after post")` where first argument is flag name, second is default value and third is a usage and then `flag.Parse()` must be called, to define flag of nonstandard type methods `String() string` and `Set(in string) error` should be provided
- 

## Sources
- [Разработка веб-сервисов на Go, часть 2. Лекция 3](golang-6.pdf)
- [Кросс-компиляция в Go](https://habr.com/ru/post/249449/)
- [С-вызовы в Go: принцип работы и производительность](https://habr.com/ru/company/intel/blog/275709/)
- [Essential Go](https://www.programming-books.io/essential/go/)

