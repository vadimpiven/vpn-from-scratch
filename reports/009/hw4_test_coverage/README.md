# HTTP-запросы и покрытие тестами

## Описание
Это комбинированное задание по тому, как отправлять запросы, получать ответы, работать с параметрами, хедерами, а так же писать тесты. Задание не сложное, основной объёма работы - написание разных условий и тестов, чтобы эти условия удовлетворить.

## Дано
У нас есть какой-то поисковый сервис:
- `SearchClient` - структура с методом `FindUsers`, который отправляет запрос во внешнюю систему и возвращает результат, немного преобразуя его
- `SearchServer` - своего рода внешняя система. Непосредственно занимается поиском данных в файле `dataset.xml`. В продакшене бы запускалась в виде отдельного веб-сервиса.

## Требуется
- Написать функцию `SearchServer` в файле `client_test.go`, который вы будете запускать через тестовый сервер (`httptest.NewServer`)
- Покрыть тестами метод `FindUsers`, чтобы покрытие было 100%. Тесты писать в `client_test.go`.
- Так же требуется сгенерировать HTML-отчет с покрытием.

## Дополнительно
- Данные для работы лежат в файле `dataset.xml`
- Параметр `query` ищет по полям `Name` и `About`
- Параметр `order_field` работает по полям `Id`, `Age`, `Name`, если пустой - то возвращаем по `Name`, если что-то другое - `SearchServer` ругается ошибкой. `Name` - это `first_name + last_name` из XML.
- Если `query` пустой, то делаем только сортировку, т.е. возвращаем все записи
- Код нужно писать в файле `client_test.go`. Там будут и ваши тесты, и функция `SearchServer`
- Как работать с XML смотрите в `xml/*`
- Запускать как `go test -cover`
- Построение покрытия: `go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`

## Советы
- [Документация](https://golang.org/pkg/net/http/) может помочь
- Не запихивайте всё в 1 тест, напишите много маленьких
- Вы можете не ограничиваться функцией `SearchServer` при тестировании, если вам надо проверить какой-то совсем отдельный хитрый кейс, вроде ошибки. Но таких случаев будет немного. В основном всё будет в `SearchServer`
- Для покрытия тестом одной из ошибок придётся залезть в исходники функции, которая возвращает эту ошибку, и посмотреть при каких условиях работы или входных данных это происходит
- Производительность, горутины и прочий асинхрон в этом задании не нужны
- Не пытайтесь реализовать таймаут подключением к неизвестному IP. В авто-грейдере вообще нет сети по соображениям безопасности и такого рода подключения сразу возвращают ошибку
