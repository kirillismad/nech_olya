package notesdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
)

const createTableStatement = `
	CREATE TABLE profiles (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NULL,
		surname TEXT NULL,
		age INTEGER NULL
	)
`

var jsonRequest = []byte(`[
	{
		"id": 1,
		"name": "John",
		"surname": "Doe",
		"age": 30
	},
	{
		"id": 2,
		"name": null,
		"surname": "Smith",
		"age": 21
	},
	{
		"id": 3,
		"name": "Alice",
		"surname": null,
		"age": 25
	},
	{
		"id": 4,
		"name": "Bob",
		"surname": "Johnson",
		"age": null
	},
	{
		"id": 5,
		"name": null,
		"surname": null,
		"age": null
	}
]`)

func TestNull(t *testing.T) {
	// ProfileRequest моделирует входящий JSON от клиента.
	// Важно: для nullable-полей используем указатели (*string, *int64),
	// чтобы отличать два разных состояния:
	// 1) поле пришло как null (или не пришло) -> nil
	// 2) поле пришло с реальным значением -> не nil
	type ProfileRequest struct {
		ID int64 `json:"id"` // Не nullable поле: всегда ожидаем id.

		// Указатель на string: nil означает отсутствие значения (NULL).
		Name *string `json:"name,omitempty"`

		// Аналогично Name: nil <-> NULL.
		Surname *string `json:"surname,omitempty"`

		// Указатель на число: nil <-> NULL, не nil -> есть возраст.
		Age *int64 `json:"age,omitempty"`
	}

	// ProfileRow описывает строку, которую мы читаем/пишем в БД.
	// Здесь показываем сразу несколько подходов к nullable:
	// - sql.NullString
	// - *string
	// - sql.Null[int64] (generic-тип в Go 1.22+)
	type ProfileRow struct {
		ID int64 // Первичный ключ, не NULL.

		// NullString хранит и значение, и флаг Valid.
		// Если Valid=false, значит в БД был NULL.
		Name sql.NullString

		// Вариант через указатель: nil означает NULL.
		Surname *string

		// Generic nullable-тип для чисел (аналог NullInt64).
		// Valid=false -> NULL, Valid=true -> значение лежит в V.
		Age sql.Null[int64]
	}

	// Поднимаем тестовую БД (в памяти или во временном окружении).
	db := NewTestDB(t)

	// Создаем таблицу с nullable-колонками name/surname/age.
	_, err := db.Exec(createTableStatement)
	if err != nil {
		// Любая ошибка DDL должна сразу падать в тесте.
		t.Fatal(err)
	}
	t.Logf("Table created successfully")

	// Сюда распарсим JSON-массив профилей.
	var profilesRequest []ProfileRequest

	// Декодируем JSON в Go-структуры.
	// Для nullable-полей получим либо nil, либо указатель на значение.
	err = json.Unmarshal(jsonRequest, &profilesRequest)
	if err != nil {
		// Если JSON некорректный, дальнейшая проверка бессмысленна.
		t.Fatal(err)
	}

	// Готовим слайс строк для вставки в БД.
	// Емкость сразу = количеству входных профилей, чтобы избежать лишних realloc.
	insertRows := make([]ProfileRow, 0, len(profilesRequest))

	// Проходим по каждому входному профилю и конвертируем типы
	// из API-представления в DB-представление.
	for _, requestProfile := range profilesRequest {
		// Добавляем одну подготовленную строку в insertRows.
		insertRows = append(insertRows, ProfileRow{
			// ID копируем напрямую: это non-null поле.
			ID: requestProfile.ID,

			// Конвертируем *string -> sql.NullString.
			Name: func() sql.NullString {
				// Если указатель nil, значит значение отсутствует (NULL).
				if requestProfile.Name == nil {
					return sql.NullString{Valid: false}
				}

				// Иначе переносим строку и ставим Valid=true.
				return sql.NullString{String: *requestProfile.Name, Valid: true}
			}(),

			// Для Surname оставляем указатель как есть:
			// nil/не nil уже корректно кодирует NULL/значение.
			Surname: requestProfile.Surname,

			// Конвертируем *int64 -> sql.Null[int64].
			Age: func() sql.Null[int64] {
				// nil означает отсутствие значения -> Valid=false.
				if requestProfile.Age == nil {
					return sql.Null[int64]{Valid: false}
				}

				// Есть значение: кладем в V и помечаем Valid=true.
				return sql.Null[int64]{V: *requestProfile.Age, Valid: true}
			}(),
		})
	}

	// Начинаем собирать INSERT для batch-вставки нескольких строк.
	query := "INSERT INTO profiles (id, name, surname, age) VALUES "
	// Здесь будем хранить значения для плейсхолдеров (?, ?, ?, ?).
	args := []any{}

	// Формируем SQL и список аргументов динамически.
	for i, p := range insertRows {
		// Между кортежами значений добавляем запятую.
		if i > 0 {
			query += ", "
		}

		// Один профиль = один набор из 4 плейсхолдеров.
		query += "(?, ?, ?, ?)"

		// В args кладем значения в том же порядке, что и в SQL.
		args = append(args, p.ID, p.Name, p.Surname, p.Age)
	}

	t.Log("SQL for batch INSERT:", query)
	t.Log("Args for batch INSERT:", args)

	// Выполняем batch INSERT.
	result, err := db.Exec(query, args...)
	if err != nil {
		// Ошибка вставки означает, что конвертация nullable-данных сломана,
		// либо нарушена схема таблицы.
		t.Fatal(err)
	}
	// Проверяем, сколько строк реально вставилось.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Rows Affected: %d", rowsAffected)

	// Читаем обратно все строки, чтобы проверить, как NULL пришел из БД.
	rows, err := db.Query("SELECT id, name, surname, age FROM profiles")
	if err != nil {
		// Ошибка запроса делает тест невалидным.
		t.Fatal(err)
	}

	// Всегда закрываем rows, чтобы освободить ресурсы драйвера/соединения.
	defer rows.Close()

	// Сюда соберем результат чтения из БД.
	selectRows := make([]ProfileRow, 0)

	// Идем по каждой строке результата.
	for rows.Next() {
		// zero-value важен:
		// Name.Valid=false, Surname=nil, Age.Valid=false по умолчанию.
		var p ProfileRow

		// Scan заполняет поля структуры из текущей строки результата.
		// Для nullable колонок корректно выставляются Valid/nil.
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Age); err != nil {
			// Любая ошибка сканирования должна падать тестом.
			t.Fatal(err)
		}

		// Сохраняем прочитанную строку.
		selectRows = append(selectRows, p)
	}

	// Проверяем ошибки итерации после завершения цикла rows.Next().
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	// Для наглядности печатаем каждую строку в лог теста.
	// Это помогает увидеть, как NULL отображается в каждом подходе.
	for _, p := range selectRows {
		// Начинаем с обязательного поля ID.
		resultString := fmt.Sprintf("ID: %d", p.ID)

		// sql.NullString: читаем значение только когда Valid=true.
		if p.Name.Valid {
			resultString += ", Name: " + p.Name.String
		} else {
			resultString += ", Name: NULL"
		}

		// *string: nil означает NULL.
		if p.Surname != nil {
			resultString += ", Surname: " + *p.Surname
		} else {
			resultString += ", Surname: NULL"
		}

		// sql.Null[int64]: когда Valid=true, значение берем из поля V.
		if p.Age.Valid {
			resultString += fmt.Sprintf(", Age: %d", p.Age.V)
		} else {
			resultString += ", Age: NULL"
		}

		// Выводим итоговую строку в test log.
		t.Log(resultString)
	}
}
