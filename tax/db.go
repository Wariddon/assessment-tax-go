package tax

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

type Test struct {
	ID   int    `json:"id"`
	Test string `json:"test"`
}

var db *sql.DB

func init() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	db = conn
}

// Test godoc
// @summary Health Check
// @description Health checking for the service
// @id Test
// @produce plain
// @router /test [get]
func GetTest(c echo.Context) error {
	row := db.QueryRow("SELECT * FROM test")
	tst := Test{}
	err := row.Scan(&tst.ID, &tst.Test)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	//fmt.Printf("tst % #v\n", tst)

	return c.JSON(http.StatusOK, tst)
}

func GetAllowances() (map[string]float64, error) {

	rows, _ := db.Query("SELECT * FROM allowances")
	allowances := make(map[string]float64)

	for rows.Next() {
		var id int
		var allowance string
		var amount float64
		err := rows.Scan(&id, &allowance, &amount)
		if err != nil {
			return nil, err
		}
		allowances[allowance] = amount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allowances, nil
}
