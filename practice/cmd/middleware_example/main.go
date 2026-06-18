package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// curl -X GET http://localhost:9090/hello -H "Authorization: super-secret-token"

// authToken задаёт ожидаемое значение заголовка Authorization,
// которое authMiddleware использует для простой проверки доступа к endpoint'у.
const authToken = "super-secret-token"

func main() {
	mux := http.NewServeMux()

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("hello"))
	})
	handlerAndAuth := authMiddleware(handler)
	handlerAndAuthAndLogging := loggingMiddleware(handlerAndAuth)

	mux.Handle("GET /hello", recoveryMiddleware(handlerAndAuthAndLogging))

	if err := http.ListenAndServe(":9090", mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// recoveryMiddleware оборачивает следующий HTTP-обработчик и защищает сервер от падения при panic.
//
// Если в цепочке обработки запроса возникает panic, middleware перехватывает его через recover,
// пишет в лог сообщение о восстановлении и возвращает клиенту ответ с кодом 500 Internal Server Error.
// Это позволяет не завершать весь сервер из-за ошибки в одном запросе и делает поведение приложения
// более устойчивым к неожиданным сбоям внутри handler'ов.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[recoveryMiddleware] before serve: %s %s", request.Method, request.URL.Path)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[recoveryMiddleware] panic recovered: %v", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, request)

		log.Printf("[recoveryMiddleware] after serve: %s %s", request.Method, request.URL.Path)
	})
}

// loggerWriter расширяет стандартный http.ResponseWriter и сохраняет
// метаданные ответа: HTTP-статус и количество отправленных байт.
//
// Используется в loggingMiddleware для сбора информации о результате
// обработки запроса без изменения логики основного handler'а.
type loggerWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

// WriteHeader перехватывает установку HTTP-статуса,
// сохраняет его во внутреннем поле и делегирует запись
// исходному ResponseWriter.
func (w *loggerWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write перехватывает запись тела ответа, при необходимости
// устанавливает статус 200 OK по умолчанию и накапливает
// количество отправленных байт для последующего логирования.
func (w *loggerWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	written, err := w.ResponseWriter.Write(body)
	if err != nil {
		return written, fmt.Errorf("failed to write response body: %w", err)
	}

	w.bytes += written
	return written, nil
}

// loggingMiddleware логирует обработку каждого HTTP-запроса до и после вызова следующего handler'а.
//
// До передачи запроса дальше middleware фиксирует метод и путь, а также запоминает время старта.
// После выполнения next handler она записывает итоговые метрики: HTTP-статус, объём ответа в байтах
// и длительность обработки запроса. Для этого используется loggerWriter, который перехватывает
// WriteHeader и Write у стандартного ResponseWriter.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[loggingMiddleware] before serve: %s %s", request.Method, request.URL.Path)
		startedAt := time.Now()
		loggedWriter := &loggerWriter{
			ResponseWriter: writer,
		}

		next.ServeHTTP(loggedWriter, request)

		log.Printf("[loggingMiddleware] after serve: %s %s, status=%d, bytes=%d, duration=%s", request.Method, request.URL.Path, loggedWriter.statusCode, loggedWriter.bytes, time.Since(startedAt))
	})
}

// authMiddleware проверяет, что в запросе есть корректный токен авторизации в заголовке Authorization.
//
// Middleware сравнивает значение заголовка Authorization с заранее заданной константой authToken.
// Если токен отсутствует или не совпадает с ожидаемым значением, запрос не передаётся дальше,
// а клиент получает ответ 401 Unauthorized. Если токен верный, middleware пропускает запрос
// к следующему handler'у в цепочке.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[authMiddleware] before serve: %s %s", request.Method, request.URL.Path)
		if request.Header.Get("Authorization") != authToken {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)

		log.Printf("[authMiddleware] after serve: %s %s", request.Method, request.URL.Path)
	})
}
