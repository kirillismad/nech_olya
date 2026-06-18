# Middleware в Go HTTP-сервере через stdlib

Middleware в `net/http` - это обёртка вокруг `http.Handler`, которая выполняет код до вызова следующего обработчика, после него или в обоих местах.

Если упростить до одной мысли, middleware отвечает на вопрос: как добавить общую логику вокруг handler'а, не копируя её в каждый endpoint.

## Каноническая сигнатура

В стандартной библиотеке самый распространённый вид middleware такой:

```go
func(next http.Handler) http.Handler
```

Смысл сигнатуры прямой:

- middleware принимает следующий обработчик;
- возвращает новый обработчик;
- новый обработчик решает, что сделать до `next.ServeHTTP`, что после, а в некоторых случаях вызывать ли `next` вообще.

Канонический шаблон реализации выглядит так:

```go
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // код до next

        next.ServeHTTP(w, r)

        // код после next
    })
}
```

Это и есть базовая механика почти всех middleware в Go HTTP-серверах.

## Почему middleware вообще удобны

Представь, что у тебя есть несколько endpoint'ов, и для каждого нужно делать одно и то же:

- писать лог запроса;
- проверять авторизацию;
- восстанавливаться после `panic`;
- проставлять CORS-заголовки;
- добавлять `request ID`.

Если размазывать это по каждому handler'у, код быстро становится шумным. Middleware выносят такую логику в отдельные переиспользуемые обёртки.

Тогда сам business handler может остаться простым:

```go
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("hello"))
})
```

А общая инфраструктурная логика собирается снаружи.

## Пример из `main.go`

Цепочка собрана по шагам:

```go
handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(http.StatusOK)
    _, _ = writer.Write([]byte("hello"))
})

handlerAndLogger := loggingMiddleware(handler)
handlerAndLoggerAndAuth := authMiddleware(handlerAndLogger)

mux.Handle("GET /hello", recoveryMiddleware(handlerAndLoggerAndAuth))
```

Если записать это в одну строку, получится:

```go
recoveryMiddleware(authMiddleware(loggingMiddleware(handler)))
```

Это важная форма записи, потому что по ней сразу видно структуру вложенности.

## Почему middleware комбинируются снаружи внутрь

Когда ты пишешь цепочку вида:

```go
recover(logging(auth(mux)))
```

выполнение начинается с самого внешнего слоя.

То есть порядок входа такой:

1. `recover`
2. `logging`
3. `auth`
4. `mux` или конкретный handler

А потом стек разматывается обратно наружу.

Если у внутренних слоёв есть код после `next.ServeHTTP`, он выполнится в обратном порядке:

1. handler завершился
2. управление вернулось в `auth`
3. потом в `logging`
4. потом в `recover`

Это обычная вложенная композиция функций. Внешняя обёртка вызывает внутреннюю, та вызывает следующую, и так до самого handler'а.

Для интуиции полезно думать так:

- middleware не стоят рядом;
- middleware вложены друг в друга;
- каждый следующий слой находится внутри предыдущего.

### Схема вызова

```mermaid
flowchart LR
    A[1. recover] --> B[2. logging]
    B --> C[3. auth]
    C --> D[4. handler]
    D --> E[5. return to auth]
    E --> F[6. return to logging]
    F --> G[7. return to recover]
```

На входе запрос идёт слева направо. После выполнения handler'а управление возвращается справа налево.

## Что происходит в вашем `main.go`

Порядок именно такой:

```go
recovery(auth(logging(handler)))
```

Значит, при запросе к `/hello` сервер заходит в middleware в таком порядке:

1. `recoveryMiddleware`
2. `authMiddleware`
3. `loggingMiddleware`
4. сам handler

Почему это важно:

- `recovery` ставят снаружи, чтобы поймать `panic` из любой внутренней части цепочки;
- `auth` может остановить запрос раньше и не пустить его в handler;
- `logging` в такой конфигурации увидит только те запросы, которые уже прошли `auth`.

Это хороший пример того, что порядок влияет на результат.
Если ты хочешь логировать вообще все входящие запросы, включая `401 Unauthorized`, `logging` нужно ставить снаружи `auth`.

На практике порядок middleware - это не косметика, а часть поведения сервера.

## Разберём шаблон middleware на простом примере

Вот минимальный middleware, который пишет лог до и после вызова следующего обработчика:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("before: %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        log.Printf("after: %s %s", r.Method, r.URL.Path)
    })
}
```

Здесь важно понять две вещи:

- пока middleware не вызвал `next.ServeHTTP`, запрос не пошёл дальше;
- если middleware вообще не вызывает `next.ServeHTTP`, цепочка на нём заканчивается.

Второй случай как раз используется в аутентификации, rate limiting и других проверках доступа.

## Middleware может прервать цепочку

`authMiddleware` работает именно так: он проверяет заголовок `Authorization` и, если токен неверный, не вызывает следующий handler.

Упрощённо это выглядит так:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") != authToken {
            http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

Это принципиальный момент: middleware не обязан пропускать запрос дальше. Он может закончить обработку сам.

Поэтому auth, CORS preflight, rate limiting и некоторые валидации часто живут именно в middleware.

## `recover` нужен, чтобы не уронить весь сервер из-за одной паники

`recoveryMiddleware` использует `defer` и `recover()`:

```go
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

Смысл такой:

- если внутри handler'а или любого внутреннего middleware случится `panic`, управление перескочит в `defer`;
- `recover()` перехватит панику;
- сервер вернёт `500 Internal Server Error` вместо аварийного падения процесса.

Именно поэтому `recover` обычно ставят самым внешним слоем.

## Почему логированию мало обычного `ResponseWriter`

Новички часто ожидают, что из `http.ResponseWriter` можно просто прочитать итоговый HTTP-статус. Но стандартный интерфейс этого не даёт.

У `ResponseWriter` есть методы для записи ответа, но нет метода вроде `StatusCode() int`.

Из-за этого middleware логирования не знает автоматически:

- какой статус реально вернул handler;
- был ли вызван `WriteHeader`;
- сколько байт ушло в ответ.

Чтобы это узнать, нужно обернуть `ResponseWriter` своим типом.

## Как сделать обёртку `ResponseWriter`

Для этого используется тип со встроенным `http.ResponseWriter`:

```go
type loggerWriter struct {
    http.ResponseWriter
    statusCode int
    bytes      int
}
```

Это встраивание даёт важный эффект:

- все методы исходного `ResponseWriter` по умолчанию остаются доступны;
- ты можешь переопределить только те методы, поведение которых хочешь перехватить.

Минимально для статуса нужно переопределить `WriteHeader`:

```go
func (w *loggerWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
```

Что здесь происходит:

- middleware запоминает статус во внутреннем поле `statusCode`;
- затем пробрасывает вызов в настоящий `ResponseWriter`, чтобы ответ реально ушёл клиенту.

## Почему часто переопределяют ещё и `Write`

Есть деталь протокола `net/http`: если handler не вызвал `WriteHeader` явно, первый вызов `Write` неявно означает `200 OK`.

Поэтому только `WriteHeader` недостаточно. Нужно ещё обработать `Write`:

```go
func (w *loggerWriter) Write(body []byte) (int, error) {
    if w.statusCode == 0 {
        w.statusCode = http.StatusOK
    }

    written, err := w.ResponseWriter.Write(body)
    if err != nil {
        return written, err
    }

    w.bytes += written
    return written, nil
}
```

Теперь middleware знает:

- статус ответа;
- размер ответа в байтах.

Это и делает логирование полезным, а не формальным.

## Как логирование использует эту обёртку

После создания `loggerWriter` middleware передаёт уже его, а не исходный `ResponseWriter`, дальше по цепочке:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        startedAt := time.Now()
        loggedWriter := &loggerWriter{
            ResponseWriter: w,
        }

        next.ServeHTTP(loggedWriter, r)

        log.Printf(
            "status=%d bytes=%d duration=%s",
            loggedWriter.statusCode,
            loggedWriter.bytes,
            time.Since(startedAt),
        )
    })
}
```

Это ключевой паттерн: middleware подменяет writer на совместимый wrapper, собирает метрики и не ломает интерфейс для следующего handler'а.

## Типичные middleware в Go HTTP-сервере

Самые частые middleware такие:

- логирование запросов и ответов;
- `recover` от паник;
- аутентификация и авторизация;
- CORS;
- rate limiting;
- `request ID` или trace ID;
- метрики;
- таймауты;
- audit logging.

Полезно различать их по роли:

- одни просто наблюдают за запросом, например логирование и метрики;
- другие могут остановить цепочку, например auth и rate limiting;
- третьи защищают серверную инфраструктуру, например `recover`.

## Частые ошибки

Когда только начинаешь работать с middleware, обычно путаются в нескольких местах.

### 1. Неправильно понимают порядок выполнения

Частая ошибка - читать цепочку слева направо как список независимых шагов.

Правильнее думать так:

- внешний middleware входит первым;
- внутренний handler вызывается последним;
- код после `next.ServeHTTP` выполняется при выходе назад.

### 2. Забывают вызвать `next.ServeHTTP`

Если middleware не вызывает `next`, запрос просто не дойдёт до обработчика. Иногда это нужно специально, но часто это просто баг.

### 3. Логируют статус, не оборачивая `ResponseWriter`

Без собственного wrapper'а middleware не знает реальный статус ответа. В итоге в логах появляются нули, guessed values или просто неправильные данные.

### 4. Ставят `recover` слишком глубоко

Если `recover` находится не снаружи, часть цепочки может остаться без защиты от `panic`.

### 5. Пытаются сложить в middleware всю бизнес-логику

Middleware хорош для общей поперечной логики. Если туда уезжает основная предметная логика endpoint'а, код становится трудно читать и тестировать.

## Практическое правило выбора порядка

Хороший стартовый принцип такой:

- снаружи ставь защитные и инфраструктурные слои вроде `recover` и `request ID`;
- ближе к handler'у ставь проверки доступа и узкоспециализированную логику;
- логирование располагай так, чтобы оно видело тот объём поведения, который тебе нужен.

Например:

```go
recover(requestID(logging(auth(mux))))
```

Такой порядок означает:

- `recover` ловит панику из любого внутреннего слоя;
- `request ID` доступен для всех следующих middleware;
- `logging` может писать ID запроса и итоговый статус;
- `auth` решает, пускать ли запрос в handler.

Конкретный порядок зависит от того, что именно ты хочешь наблюдать и защищать, но сам принцип композиции остаётся одним и тем же.

## Коротко

Middleware в Go через stdlib - это функция вида `func(next http.Handler) http.Handler`, которая оборачивает следующий handler и добавляет поведение вокруг `next.ServeHTTP`.

Нужно запомнить четыре практические идеи:

- middleware вложены друг в друга, а не просто перечислены рядом;
- самый внешний слой выполняется первым;
- middleware может как пропустить запрос дальше, так и остановить цепочку;
- для логирования статуса ответа нужен wrapper над `http.ResponseWriter`.

Если эта модель понятна, код перестаёт выглядеть как магия: это просто несколько вложенных обработчиков, каждый из которых отвечает за свою часть HTTP-пайплайна.
