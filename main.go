package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/controller"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getemployee", controller.AllEmployee).Methods("GET")
	router.HandleFunc("/getdatacalon", controller.GetUserHandler).Methods("GET")
	router.HandleFunc("/getkandidat", controller.GetKandidat).Methods("GET")
	router.HandleFunc("/gettotaldata", controller.GetAllData).Methods("GET")
	router.HandleFunc("/gettotaldatakel", controller.GetAllDataKel).Methods("GET")
	router.HandleFunc("/gettotaldatakecamatan", controller.GetDataHandler).Methods("GET")
	router.HandleFunc("/getkecamatan", controller.GetKecamatan).Methods("GET")
	router.HandleFunc("/gettotaldatakelurahan", controller.GetDataKelurahan).Methods("GET")
	router.HandleFunc("/getkelurahan", controller.GetKelurahan).Methods("GET")
	router.HandleFunc("/insertemployee", controller.InsertEmployee).Methods("POST")
	router.HandleFunc("/insertkandidat", controller.InsertKandidat).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Connected to port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))
}

// package main

// import (
// 	"database/sql"

// 	"github.com/handlers"

// 	"github.com/labstack/echo"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	e := echo.New()
// 	db := initDB("storage.db")
// 	migrate(db)

// 	// daftar api
// 	e.GET("/tasks", handlers.GetTasks(db))
// 	e.POST("/tasks", handlers.PutTask(db))
// 	e.PUT("/tasks", handlers.EditTask(db))
// 	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

// 	e.Logger.Fatal(e.Start(":8000"))
// }

// func initDB(filepath string) *sql.DB {
// 	db, err := sql.Open("sqlite3", filepath)

// 	if err != nil {
// 		panic(err)
// 	}

// 	if db == nil {
// 		panic("db nil")
// 	}

// 	return db
// }

// func migrate(db *sql.DB) {
// 	sql := `
//     CREATE TABLE IF NOT EXISTS tasks(
//         id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
// 		name VARCHAR NOT NULL,
// 		status INTEGER
//     );
//     `

// 	_, err := db.Exec(sql)
// 	// Exit if something goes wrong with our SQL statement above
// 	if err != nil {
// 		panic(err)
// 	}
// }
