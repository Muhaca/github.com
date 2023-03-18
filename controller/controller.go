package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/model"

	"github.com/config"
)

// select all data with pagination
// func AllEmployee(w http.ResponseWriter, r *http.Request) {
// 	var employee model.Employee
// 	var response model.Response
// 	var arrEmployee []model.Employee

// 	db := config.Connect()
// 	defer db.Close()

// 	// Get page and perPage query parameters
// 	page, err := strconv.Atoi(r.URL.Query().Get("page"))
// 	if err != nil {
// 		page = 1
// 	}

// 	perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))
// 	if err != nil {
// 		perPage = 10
// 	}

// 	// Calculate offset based on page and perPage
// 	offset := (page - 1) * perPage

// 	// Construct SQL query with LIMIT and OFFSET clauses
// 	query := fmt.Sprintf("SELECT id, nama, kecamatan, kelurahan, user, tps, jumlah_suara, gambar FROM voting LIMIT %d OFFSET %d", perPage, offset)

// 	rows, err := db.Query(query)

// 	if err != nil {
// 		log.Print(err)
// 	}

// 	for rows.Next() {
// 		err = rows.Scan(
// 			&employee.Id,
// 			&employee.Nama,
// 			&employee.Kecamatan,
// 			&employee.Kelurahan,
// 			&employee.User,
// 			&employee.TPS,
// 			&employee.JumlahSuara,
// 			&employee.Gambar,
// 		)
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		} else {
// 			arrEmployee = append(arrEmployee, employee)
// 		}
// 	}

// 	response.Status = 200
// 	response.Message = "Success"
// 	response.Data = arrEmployee

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	json.NewEncoder(w).Encode(response)
// }

func AllEmployee(w http.ResponseWriter, r *http.Request) {
	var employee model.Employee
	var response model.Response
	var arrEmployee []model.Employee

	db := config.Connect()
	defer db.Close()

	// Get the page and perPage query parameters from the request
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))
	if err != nil {
		perPage = 10
	}

	// Calculate the offset and limit based on the page and perPage values
	offset := (page - 1) * perPage
	limit := perPage

	// Query the database to get the total number of rows
	var totalRows int
	err = db.QueryRow("SELECT COUNT(*), created_at FROM voting").Scan(&totalRows)
	if err != nil {
		log.Print(err)
	}

	// Query the database to get the employees for the current page
	rows, err := db.Query(`
		SELECT id, nama, kecamatan, kelurahan, user, tps, jumlah_suara, gambar, created_at
		FROM voting
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(
			&employee.Id,
			&employee.Nama,
			&employee.Kecamatan,
			&employee.Kelurahan,
			&employee.User,
			&employee.TPS,
			&employee.JumlahSuara,
			&employee.Gambar,
			&employee.CreatedAt,
		)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			arrEmployee = append(arrEmployee, employee)
		}
	}

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(perPage)))

	// Create the meta object for the response
	meta := model.Meta{
		Page:       page,
		PerPage:    perPage,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	response.Status = 200
	response.Message = "Success"
	response.Meta = meta
	response.Data = arrEmployee

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	json.NewEncoder(w).Encode(response)

}

// AllEmployee = Select Employee API
func GetKecamatan(w http.ResponseWriter, r *http.Request) {
	var kecamatan model.Kecamatan
	var response []model.Kecamatan
	var resp model.KecamatanResponse

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, kecamatan FROM kecamatan")

	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(&kecamatan.Id, &kecamatan.Kecamatan)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			response = append(response, kecamatan)
		}
	}

	resp.Status = 200
	resp.Message = "Success"
	resp.Data = response

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	json.NewEncoder(w).Encode(resp)

}

func GetKelurahan(w http.ResponseWriter, r *http.Request) {
	kecamatan := r.URL.Query().Get("kecamatan")
	if kecamatan == "" {
		http.Error(w, "kecamatan parameter is required", http.StatusBadRequest)
		return
	}

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, kelurahan FROM kelurahan WHERE kecamatan = ?", kecamatan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var kelurahanList []model.Kelurahan
	for rows.Next() {
		var kelurahan model.Kelurahan
		err := rows.Scan(&kelurahan.Id, &kelurahan.Kelurahan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kelurahanList = append(kelurahanList, kelurahan)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	json.NewEncoder(w).Encode(kelurahanList)
}

func GetKandidat(w http.ResponseWriter, r *http.Request) {
	var employee model.Kandidat
	var response model.ResponseKandidat
	var arrEmployee []model.Kandidat

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, nama, user_id FROM voting_kandidat")

	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err = rows.Scan(&employee.Id, &employee.Nama, &employee.UserID)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			arrEmployee = append(arrEmployee, employee)
		}
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrEmployee

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	json.NewEncoder(w).Encode(response)

}

func GetUserByID(db *sql.DB, id int) (*model.Employee, error) {
	var employee model.Employee
	err := db.QueryRow("select * from voting where id=?").Scan(&employee.Id, &employee.Nama, &employee.Kecamatan, &employee.Kelurahan, &employee.User, &employee.TPS, &employee.JumlahSuara)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id tidak ditemukan", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		http.Error(w, "database error connection", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	emp, err := GetUserByID(db, id)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(emp)

}

// InsertEmployee = Insert Employee API
func InsertEmployee(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	// id := r.FormValue("id")
	nama := r.FormValue("nama")
	kelurahan := r.FormValue("kelurahan")
	kecamatan := r.FormValue("kecamatan")
	user := r.FormValue("user")
	tps := r.FormValue("tps")
	jumlah_suara := r.FormValue("jumlah_suara")
	gambar := r.FormValue("gambar")

	_, err = db.Exec("INSERT INTO voting(nama, kelurahan, kecamatan, user,tps, jumlah_suara, gambar, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, now())", nama, kelurahan, kecamatan, user, tps, jumlah_suara, gambar)

	if err != nil {
		log.Print(err)
		return
	}
	response.Status = 200
	response.Message = "Insert data successfully"
	fmt.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	json.NewEncoder(w).Encode(response)
}

// InsertEmployee = Insert Employee API
func InsertKandidat(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	// id := r.FormValue("id")
	nama := r.FormValue("nama")
	userid := r.FormValue("user_id")

	_, err = db.Exec("INSERT INTO voting_kandidat(nama, user_id) VALUES(?, ?)", nama, userid)

	if err != nil {
		log.Print(err)
		return
	}
	response.Status = 200
	response.Message = "Insert data successfully"
	fmt.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	json.NewEncoder(w).Encode(response)
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	kecamatan := r.URL.Query().Get("kecamatan")
	db := config.Connect()
	defer db.Close()

	var rows *sql.Rows
	var err error
	if kecamatan == "" {
		rows, err = db.Query("SELECT sum(jumlah_suara) as total_suara, nama, kecamatan FROM voting GROUP BY nama, kecamatan")
	} else {
		rows, err = db.Query("SELECT sum(jumlah_suara) as total_suara, nama FROM voting WHERE kecamatan = ? GROUP BY nama", kecamatan)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var totalSuara int64
		var nama string
		var kec string
		if kecamatan == "" {
			err := rows.Scan(&totalSuara, &nama, &kec)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := rows.Scan(&totalSuara, &nama)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		result := map[string]interface{}{
			"total_suara": totalSuara,
			"nama":        nama,
		}
		if kecamatan == "" {
			result["kecamatan"] = kec
		}
		results = append(results, result)
	}

	jsonResult, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

func GetDataKelurahan(w http.ResponseWriter, r *http.Request) {
	kelurahan := r.URL.Query().Get("kelurahan")
	db := config.Connect()
	defer db.Close()

	var rows *sql.Rows
	var err error
	if kelurahan == "" {
		rows, err = db.Query("SELECT sum(jumlah_suara) as total_suara, nama, kelurahan FROM voting GROUP BY nama, kelurahan")
	} else {
		rows, err = db.Query("SELECT sum(jumlah_suara) as total_suara, nama FROM voting WHERE kelurahan = ? GROUP BY nama", kelurahan)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var totalSuara int64
		var nama string
		var kel string
		if kelurahan == "" {
			err := rows.Scan(&totalSuara, &nama, &kel)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := rows.Scan(&totalSuara, &nama)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		result := map[string]interface{}{
			"total_suara": totalSuara,
			"nama":        nama,
		}
		if kelurahan == "" {
			result["kelurahan"] = kel
		}
		results = append(results, result)
	}

	jsonResult, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResult)
}

type Result struct {
	Kecamatan  string `json:"kecamatan"`
	Nama1      string `json:"nama_1"`
	Nama2      string `json:"nama_2"`
	Nama3      string `json:"nama_3"`
	TotalSuara int    `json:"total_suara"`
}

func GetAllData(w http.ResponseWriter, r *http.Request) {
	// Open database connection
	db := config.Connect()
	defer db.Close()
	var rows *sql.Rows
	var err error

	// Execute the query
	kecamatan := r.URL.Query().Get("kecamatan")

	if kecamatan == "" {
		rows, err = db.Query(`
        SELECT kecamatan,													
		MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
		MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
		MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
		MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
		MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
		MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
		COALESCE(MAX(CASE WHEN rn = 4 THEN nama END), '') AS nama_4, 
		COALESCE(MAX(CASE WHEN rn = 4 THEN jumlah_suara END), 0) AS total_suara_4,
		COALESCE(MAX(CASE WHEN rn = 5 THEN nama END), '') AS nama_5, 
		COALESCE(MAX(CASE WHEN rn = 5 THEN jumlah_suara END), 0) AS total_suara_5,
		COALESCE(MAX(CASE WHEN rn = 6 THEN nama END), '') AS nama_6, 
		COALESCE(MAX(CASE WHEN rn = 6 THEN jumlah_suara END), 0) AS total_suara_6,
		COALESCE(MAX(CASE WHEN rn = 7 THEN nama END), '') AS nama_7, 
		COALESCE(MAX(CASE WHEN rn = 7 THEN jumlah_suara END), 0) AS total_suara_7,
		COALESCE(MAX(CASE WHEN rn = 8 THEN nama END), '') AS nama_8, 
		COALESCE(MAX(CASE WHEN rn = 8 THEN jumlah_suara END), 0) AS total_suara_8,
		COALESCE(MAX(CASE WHEN rn = 9 THEN nama END), '') AS nama_9, 
		COALESCE(MAX(CASE WHEN rn = 9 THEN jumlah_suara END), 0) AS total_suara_9,
		COALESCE(MAX(CASE WHEN rn = 10 THEN nama END), '') AS nama_10, 
		COALESCE(MAX(CASE WHEN rn = 10 THEN jumlah_suara END), 0) AS total_suara_10
        FROM (
            SELECT kecamatan, 
                   nama, 
                   jumlah_suara, 
                   ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS rn
            FROM voting
        ) t
        WHERE rn <= 10
        GROUP BY kecamatan
    `)
	} else {
		rows, err = db.Query(`
        SELECT 	kecamatan,													
		MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
		MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
		MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
		MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
		MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
		MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
		COALESCE(MAX(CASE WHEN rn = 4 THEN nama END), '') AS nama_4, 
		COALESCE(MAX(CASE WHEN rn = 4 THEN jumlah_suara END), 0) AS total_suara_4,
		COALESCE(MAX(CASE WHEN rn = 5 THEN nama END), '') AS nama_5, 
		COALESCE(MAX(CASE WHEN rn = 5 THEN jumlah_suara END), 0) AS total_suara_5,
		COALESCE(MAX(CASE WHEN rn = 6 THEN nama END), '') AS nama_6, 
		COALESCE(MAX(CASE WHEN rn = 6 THEN jumlah_suara END), 0) AS total_suara_6,
		COALESCE(MAX(CASE WHEN rn = 7 THEN nama END), '') AS nama_7, 
		COALESCE(MAX(CASE WHEN rn = 7 THEN jumlah_suara END), 0) AS total_suara_7,
		COALESCE(MAX(CASE WHEN rn = 8 THEN nama END), '') AS nama_8, 
		COALESCE(MAX(CASE WHEN rn = 8 THEN jumlah_suara END), 0) AS total_suara_8,
		COALESCE(MAX(CASE WHEN rn = 9 THEN nama END), '') AS nama_9, 
		COALESCE(MAX(CASE WHEN rn = 9 THEN jumlah_suara END), 0) AS total_suara_9,
		COALESCE(MAX(CASE WHEN rn = 10 THEN nama END), '') AS nama_10, 
		COALESCE(MAX(CASE WHEN rn = 10 THEN jumlah_suara END), 0) AS total_suara_10
        FROM (
            SELECT kecamatan, 
                   nama, 
                   jumlah_suara, 
                   ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS rn
            FROM voting

			where kecamatan = ?
        ) t
        WHERE rn <= 10
        GROUP BY kecamatan
    `, kecamatan)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create the response data
	var result []map[string]interface{}
	for rows.Next() {
		var kecamatan string
		var suara1, suara2, suara3, suara4, suara5, suara6, suara7, suara8, suara9, suara10 int
		var nama1, nama2, nama3, nama4, nama5, nama6, nama7, nama8, nama9, nama10 string
		err = rows.Scan(&kecamatan, &nama1, &suara1, &nama2, &suara2, &nama3, &suara3, &nama4, &suara4, &nama5, &suara5, &nama6, &suara6, &nama7, &suara7, &nama8, &suara8, &nama9, &suara9, &nama10, &suara10)
		if err != nil {
			log.Fatal(err)
		}
		suara := []map[string]interface{}{
			{"nama_1": nama1, "total_suara_1": suara1},
			{"nama_2": nama2, "total_suara_2": suara2},
			{"nama_3": nama3, "total_suara_3": suara3},
			{"nama_4": nama4, "total_suara_4": suara4},
			{"nama_5": nama5, "total_suara_5": suara5},
			{"nama_6": nama6, "total_suara_6": suara6},
			{"nama_7": nama7, "total_suara_7": suara7},
			{"nama_8": nama8, "total_suara_8": suara8},
			{"nama_9": nama9, "total_suara_9": suara9},
			{"nama_10": nama10, "total_suara_10": suara10},
		}
		result = append(result, map[string]interface{}{"kecamatan": kecamatan, "suara": suara})
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// Encode and return the response data
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllDataKel(w http.ResponseWriter, r *http.Request) {
	// Open database connection

	kecamatan := r.URL.Query().Get("kecamatan")
	var rows *sql.Rows
	var err error
	db := config.Connect()
	defer db.Close()

	// Execute the query
	if kecamatan == "" {
		rows, err = db.Query(`
        SELECT kelurahan, 
		MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
		MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
		MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
		MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
		MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
		MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
		COALESCE(MAX(CASE WHEN rn = 4 THEN nama END), '') AS nama_4, 
		COALESCE(MAX(CASE WHEN rn = 4 THEN jumlah_suara END), 0) AS total_suara_4,
		COALESCE(MAX(CASE WHEN rn = 5 THEN nama END), '') AS nama_5, 
		COALESCE(MAX(CASE WHEN rn = 5 THEN jumlah_suara END), 0) AS total_suara_5,
		COALESCE(MAX(CASE WHEN rn = 6 THEN nama END), '') AS nama_6, 
		COALESCE(MAX(CASE WHEN rn = 6 THEN jumlah_suara END), 0) AS total_suara_6,
		COALESCE(MAX(CASE WHEN rn = 7 THEN nama END), '') AS nama_7, 
		COALESCE(MAX(CASE WHEN rn = 7 THEN jumlah_suara END), 0) AS total_suara_7,
		COALESCE(MAX(CASE WHEN rn = 8 THEN nama END), '') AS nama_8, 
		COALESCE(MAX(CASE WHEN rn = 8 THEN jumlah_suara END), 0) AS total_suara_8,
		COALESCE(MAX(CASE WHEN rn = 9 THEN nama END), '') AS nama_9, 
		COALESCE(MAX(CASE WHEN rn = 9 THEN jumlah_suara END), 0) AS total_suara_9,
		COALESCE(MAX(CASE WHEN rn = 10 THEN nama END), '') AS nama_10, 
		COALESCE(MAX(CASE WHEN rn = 10 THEN jumlah_suara END), 0) AS total_suara_10
        FROM (
            SELECT kelurahan, 
                   nama, 
                   jumlah_suara, 
                   ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS rn
            FROM voting
        ) t
        WHERE rn <= 10
        GROUP BY kelurahan
    `)
	} else {
		rows, err = db.Query(`
        SELECT kelurahan, 
		MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
		MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
		MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
		MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
		MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
		MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
		COALESCE(MAX(CASE WHEN rn = 4 THEN nama END), '') AS nama_4, 
		COALESCE(MAX(CASE WHEN rn = 4 THEN jumlah_suara END), 0) AS total_suara_4,
		COALESCE(MAX(CASE WHEN rn = 5 THEN nama END), '') AS nama_5, 
		COALESCE(MAX(CASE WHEN rn = 5 THEN jumlah_suara END), 0) AS total_suara_5,
		COALESCE(MAX(CASE WHEN rn = 6 THEN nama END), '') AS nama_6, 
		COALESCE(MAX(CASE WHEN rn = 6 THEN jumlah_suara END), 0) AS total_suara_6,
		COALESCE(MAX(CASE WHEN rn = 7 THEN nama END), '') AS nama_7, 
		COALESCE(MAX(CASE WHEN rn = 7 THEN jumlah_suara END), 0) AS total_suara_7,
		COALESCE(MAX(CASE WHEN rn = 8 THEN nama END), '') AS nama_8, 
		COALESCE(MAX(CASE WHEN rn = 8 THEN jumlah_suara END), 0) AS total_suara_8,
		COALESCE(MAX(CASE WHEN rn = 9 THEN nama END), '') AS nama_9, 
		COALESCE(MAX(CASE WHEN rn = 9 THEN jumlah_suara END), 0) AS total_suara_9,
		COALESCE(MAX(CASE WHEN rn = 10 THEN nama END), '') AS nama_10, 
		COALESCE(MAX(CASE WHEN rn = 10 THEN jumlah_suara END), 0) AS total_suara_10
        FROM (
            SELECT kelurahan, 
                   nama, 
                   jumlah_suara, 
                   ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS rn
            FROM voting

			where kecamatan = ?
        ) t
        WHERE rn <= 10
        GROUP BY kelurahan
    `, kecamatan)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create the response data
	var result []map[string]interface{}
	for rows.Next() {
		var kelurahan string
		var suara1, suara2, suara3, suara4, suara5, suara6, suara7, suara8, suara9, suara10 sql.NullInt64
		var nama1, nama2, nama3, nama4, nama5, nama6, nama7, nama8, nama9, nama10 sql.NullString
		err = rows.Scan(&kelurahan, &nama1, &suara1, &nama2, &suara2, &nama3, &suara3, &nama4, &suara4, &nama5, &suara5, &nama6, &suara6, &nama7, &suara7, &nama8, &suara8, &nama9, &suara9, &nama10, &suara10)
		if err != nil {
			log.Fatal(err)
		}
		suara := []map[string]interface{}{
			{"nama_1": replaceNullString(nama1), "total_suara_1": replaceNullInt(suara1)},
			{"nama_2": replaceNullString(nama2), "total_suara_2": replaceNullInt(suara2)},
			{"nama_3": replaceNullString(nama3), "total_suara_3": replaceNullInt(suara3)},
			{"nama_4": replaceNullString(nama4), "total_suara_4": replaceNullInt(suara4)},
			{"nama_5": replaceNullString(nama5), "total_suara_5": replaceNullInt(suara5)},
			{"nama_6": replaceNullString(nama6), "total_suara_6": replaceNullInt(suara6)},
			{"nama_7": replaceNullString(nama7), "total_suara_7": replaceNullInt(suara7)},
			{"nama_8": replaceNullString(nama8), "total_suara_8": replaceNullInt(suara8)},
			{"nama_9": replaceNullString(nama9), "total_suara_9": replaceNullInt(suara9)},
			{"nama_10": replaceNullString(nama10), "total_suara_10": replaceNullInt(suara10)},
		}
		result = append(result, map[string]interface{}{"kelurahan": kelurahan, "suara": suara})
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// Encode and return the response data
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatal(err)
	}
}

// Helper function to replace null string with empty string
func replaceNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// Helper function to replace null integer with zero
func replaceNullInt(i sql.NullInt64) int {
	if i.Valid {
		return int(i.Int64)
	}
	return 0
}
