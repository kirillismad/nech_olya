# Transport HTTP

## http.Transport (connection pool)

```golang
t := &http.Transport{
    MaxIdleConns: 500,
    IdleConnTimeout: 90 * time.Second,
    MaxIdleConnsPerHost: 100,
}

client := &http.Client{Transport: t}

r, err := client.Get("http://google.com/search?q=text")
if err != nil {
    log.Fatal(err)
}
defer r.Body.Close() // Body = io.ReadCloser

_ = r.Body
io.Copy(io.Discard, r.Body) // Чтение тела ответа и его закрытие

```

нагрузка = 100000 запросов в секунду
время запроса = 100 мс

GET /search?q=text

| HOST | NUM_IDLE_CONNS | TOTAL_CONNS |
| --- | --- | --- |
| google.com | 0 | 100 |
| yandex.ru | 0 | 100 |

