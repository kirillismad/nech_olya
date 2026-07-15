package notesdb

import (
	"database/sql"
	"fmt"
	"testing"
)

/*

create table profiles (
	id integer PRIMARY KEY NOT NULL, -- .... -1, 0, 1, ...
	name text NULL,
	surname text NULL,
	age integer NULL
)

*/

func TestNull(t *testing.T) {
	// Воображаемый реквест, который мы получили от клиента. В реальном приложении это будет приходить в виде JSON.
	// Поля могут быть null, поэтому мы используем указатели, чтобы отличать "поле указано" от "поле не указано".
	// Иногда zero value может быть признаком того что поле не было указано. Это актуально для тех случаев когда значение не может быть физически равно zero value
	// Например, можно было бы сделать Name string, и просто воспринимать пустую строку как "не указано". В реальной жизни, в таких случаях, разработчики обычно просто договариваются что использовать zero value или указатель.
	type Request struct {
		ID      int64   `json:"id"`
		Name    *string `json:"name,omitempty"`
		Surname *string `json:"surname,omitempty"`
		Age     *int64  `json:"age,omitempty"`
	}
	// Структура, которая будет использоваться для сканирования данных из базы данных. Здесь мы используем sql.NullString, *string (указатель) и sql.Null[int64] для обработки возможных NULL значений.
	type ProfileRow struct {
		ID      int64
		Name    sql.NullString
		Surname *string
		Age     sql.Null[int64] // go 1.22, джинерики (generic types) = обобщения типов
	}

	db := newTestDB(t)

	_, err := db.Exec("CREATE TABLE profiles (id INTEGER PRIMARY KEY NOT NULL, name TEXT NULL, surname TEXT NULL, age INTEGER NULL)")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO profiles (id, name, surname, age) VALUES (1, 'John', 'Doe', 30), (2, NULL, 'Smith', NULL), (3, 'Alice', NULL, 25)")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, surname, age FROM profiles")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	profiles := make([]ProfileRow, 0)

	for rows.Next() {
		var p ProfileRow // zero value of ProfileRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Age); err != nil {
			t.Fatal(err)
		}
		profiles = append(profiles, p)
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	for _, p := range profiles {
		resultString := fmt.Sprintf("ID: %d", p.ID)

		if p.Name.Valid {
			resultString += ", Name: " + p.Name.String
		} else {
			resultString += ", Name: NULL"
		}

		if p.Surname != nil {
			resultString += ", Surname: " + *p.Surname
		} else {
			resultString += ", Surname: NULL"
		}

		if p.Age.Valid {
			resultString = fmt.Sprintf("%s, Age: %d", resultString, p.Age.V)
		} else {
			resultString += ", Age: NULL"
		}

		t.Log(resultString)
	}

}
