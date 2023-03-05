package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/model"

	"github.com/config"
)

// AllEmployee = Select Employee API
func AllEmployee(w http.ResponseWriter, r *http.Request) {
	var employee model.Employee
	var response model.Response
	var arrEmployee []model.Employee

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, nama, kecamatan, kelurahan, user, tps, jumlah_suara, gambar FROM voting")

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
		)
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

	_, err = db.Exec("INSERT INTO voting(nama, kelurahan, kecamatan, user,tps, jumlah_suara, gambar) VALUES(?, ?, ?, ?, ?, ?, ?)", nama, kelurahan, kecamatan, user, tps, jumlah_suara, gambar)

	if err != nil {
		log.Print(err)
		return
	}
	response.Status = 200
	response.Message = "Insert data successfully"
	fmt.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

// func GetAllData(w http.ResponseWriter, r *http.Request) {
// 	// Prepare the SQL query
// 	query := "SELECT kecamatan, MAX(CASE WHEN row_number = 1 THEN nama END) AS nama_1, MAX(CASE WHEN row_number = 2 THEN nama END) AS nama_2, MAX(CASE WHEN row_number = 3 THEN nama END) AS nama_3, SUM(jumlah_suara) AS total_suara FROM ( SELECT nama, kecamatan, jumlah_suara, ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS row_number FROM voting ) subquery WHERE row_number <= 3 GROUP BY kecamatan"

// 	db := config.Connect()
// 	defer db.Close()
// 	// Execute the query and get a rows object
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// Iterate over each row and scan it into a Result object
// 	var results []Result
// 	for rows.Next() {
// 		var result Result
// 		if err := rows.Scan(&result.Kecamatan, &result.Nama1, &result.Nama2, &result.Nama3, &result.TotalSuara); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		results = append(results, result)
// 	}
// 	if err := rows.Err(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Encode the results as JSON and send them as the response body
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	if err := json.NewEncoder(w).Encode(results); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func GetAllData(w http.ResponseWriter, r *http.Request) {
// 	// open a connection to your MySQL database
// 	db := config.Connect()
// 	defer db.Close()

// 	// execute your SQL query to retrieve the sum of jumlah_suara and nama grouped by nama
// 	rows, err := db.Query("SELECT SUM(jumlah_suara) as total_suara, nama, kecamatan FROM voting GROUP BY kecamatan")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// create a slice to store the results
// 	var results []map[string]interface{}

// 	// loop through the result set and append each row to the results slice
// 	for rows.Next() {
// 		var totalSuara int
// 		var nama string
// 		var kecamatan string
// 		if err := rows.Scan(&totalSuara, &nama, &kecamatan); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		result := make(map[string]interface{})
// 		result["total_suara"] = totalSuara
// 		result["nama"] = nama
// 		result["kecamatan"] = kecamatan
// 		results = append(results, result)
// 	}

// 	// encode the results slice as JSON and write it to the response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	if err := json.NewEncoder(w).Encode(results); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func GetAllData(w http.ResponseWriter, r *http.Request) {
// 	db := config.Connect()
// 	defer db.Close()

// 	var result []map[string]interface{}
// 	rows, err := db.Query(`SELECT kecamatan,
// 	MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1,
// 	MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
// 	MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2,
// 	MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
// 	MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3,
// 	MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3
//   FROM (
// 	SELECT kecamatan,
// 	  nama,
// 	  jumlah_suara,
// 	  ROW_NUMBER() OVER (PARTITION BY kecamatan ORDER BY jumlah_suara DESC) AS rn
// 	FROM voting
//   ) t
//   WHERE rn <= 3
//   GROUP BY kecamatan`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var kecamatan string
// 		var nama1 string
// 		var total_suara1 int
// 		var nama2 sql.NullString
// 		var total_suara2 sql.NullInt64
// 		var nama3 sql.NullString
// 		var total_suara3 sql.NullInt64

// 		err := rows.Scan(&kecamatan, &nama1, &total_suara1, &nama2, &total_suara2, &nama3, &total_suara3)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		resultRow := make(map[string]interface{})
// 		resultRow["kecamatan"] = kecamatan
// 		suara := []map[string]interface{}{}
// 		if nama1 != "" {
// 			suara = append(suara, map[string]interface{}{"nama_1": nama1, "total_suara_1": total_suara1})
// 		}
// 		if nama2.Valid {
// 			suara = append(suara, map[string]interface{}{"nama_2": nama2.String, "total_suara_2": int(total_suara2.Int64)})
// 		}
// 		if nama3.Valid {
// 			suara = append(suara, map[string]interface{}{"nama_3": nama3.String, "total_suara_3": int(total_suara3.Int64)})
// 		}
// 		resultRow["suara"] = suara

// 		result = append(result, resultRow)
// 	}

// 	response, err := json.Marshal(result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }

func GetAllData(w http.ResponseWriter, r *http.Request) {
	// Open database connection
	db := config.Connect()
	defer db.Close()

	// Execute the query
	rows, err := db.Query(`
        SELECT kecamatan, 
               MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
               MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
               MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
               MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
               MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
               MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
			   MAX(CASE WHEN rn = 4 THEN nama END) AS nama_4, 
               MAX(CASE WHEN rn = 4 THEN jumlah_suara END) AS total_suara_4,
			   MAX(CASE WHEN rn = 5 THEN nama END) AS nama_5, 
               MAX(CASE WHEN rn = 5 THEN jumlah_suara END) AS total_suara_5,
			   MAX(CASE WHEN rn = 6 THEN nama END) AS nama_6, 
               MAX(CASE WHEN rn = 6 THEN jumlah_suara END) AS total_suara_6,
			   MAX(CASE WHEN rn = 7 THEN nama END) AS nama_7, 
               MAX(CASE WHEN rn = 7 THEN jumlah_suara END) AS total_suara_7,
			   MAX(CASE WHEN rn = 8 THEN nama END) AS nama_8, 
               MAX(CASE WHEN rn = 8 THEN jumlah_suara END) AS total_suara_8,
			   MAX(CASE WHEN rn = 9 THEN nama END) AS nama_9, 
               MAX(CASE WHEN rn = 9 THEN jumlah_suara END) AS total_suara_9,
			   MAX(CASE WHEN rn = 10 THEN nama END) AS nama_10, 
               MAX(CASE WHEN rn = 10 THEN jumlah_suara END) AS total_suara_10
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
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create the response data
	var result []map[string]interface{}
	for rows.Next() {
		var kecamatan string
		var suara1, suara2, suara3, suara4, suara5, suara6, suara7, suara8, suara9, suara10 sql.NullInt64
		var nama1, nama2, nama3, nama4, nama5, nama6, nama7, nama8, nama9, nama10 sql.NullString
		err = rows.Scan(&kecamatan, &nama1, &suara1, &nama2, &suara2, &nama3, &suara3, &nama4, &suara4, &nama5, &suara5, &nama6, &suara6, &nama7, &suara7, &nama8, &suara8, &nama9, &suara9, &nama10, &suara10)
		if err != nil {
			log.Fatal(err)
		}
		suara := []map[string]interface{}{
			{"nama_1": nama1, "total_suara": suara1},
			{"nama_2": nama2, "total_suara": suara2},
			{"nama_3": nama3, "total_suara": suara3},
			{"nama_4": nama4, "total_suara": suara4},
			{"nama_5": nama5, "total_suara": suara5},
			{"nama_6": nama6, "total_suara": suara6},
			{"nama_7": nama7, "total_suara": suara7},
			{"nama_8": nama8, "total_suara": suara8},
			{"nama_9": nama9, "total_suara": suara9},
			{"nama_10": nama10, "total_suara": suara10},
		}
		result = append(result, map[string]interface{}{"kecamatan": kecamatan, "suara": suara})
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Encode and return the response data
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatal(err)
	}
}

// func GetAllDataKel(w http.ResponseWriter, r *http.Request) {
// 	db := config.Connect()
// 	defer db.Close()

// 	var result []map[string]interface{}
// 	rows, err := db.Query("SELECT kelurahan, " +
// 		"MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, " +
// 		"MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1, " +
// 		"MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, " +
// 		"MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2, " +
// 		"MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, " +
// 		"MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3 " +
// 		"FROM ( " +
// 		"  SELECT kelurahan, " +
// 		"         nama, " +
// 		"         jumlah_suara, " +
// 		"         ROW_NUMBER() OVER (PARTITION BY kelurahan ORDER BY jumlah_suara DESC) AS rn " +
// 		"  FROM voting " +
// 		") t " +
// 		"WHERE rn <= 3 " +
// 		"GROUP BY kelurahan")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var kelurahan string
// 		var nama1 string
// 		var total_suara1 int
// 		var nama2 sql.NullString
// 		var total_suara2 sql.NullInt64
// 		var nama3 sql.NullString
// 		var total_suara3 sql.NullInt64

// 		err := rows.Scan(&kelurahan, &nama1, &total_suara1, &nama2, &total_suara2, &nama3, &total_suara3)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		resultRow := make(map[string]interface{})
// 		resultRow["kelurahan"] = kelurahan
// 		suara := []map[string]interface{}{}
// 		if nama1 != "" {
// 			suara = append(suara, map[string]interface{}{"nama_1": nama1, "total_suara_1": total_suara1})
// 		}
// 		if nama2.Valid {
// 			suara = append(suara, map[string]interface{}{"nama_2": nama2.String, "total_suara_2": int(total_suara2.Int64)})
// 		}
// 		if nama3.Valid {
// 			suara = append(suara, map[string]interface{}{"nama_3": nama3.String, "total_suara_3": int(total_suara3.Int64)})
// 		}
// 		resultRow["suara"] = suara

// 		result = append(result, resultRow)
// 	}

// 	response, err := json.Marshal(result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }

func GetAllDataKel(w http.ResponseWriter, r *http.Request) {
	// Open database connection
	db := config.Connect()
	defer db.Close()

	// Execute the query
	rows, err := db.Query(`
        SELECT kelurahan, 
               MAX(CASE WHEN rn = 1 THEN nama END) AS nama_1, 
               MAX(CASE WHEN rn = 1 THEN jumlah_suara END) AS total_suara_1,
               MAX(CASE WHEN rn = 2 THEN nama END) AS nama_2, 
               MAX(CASE WHEN rn = 2 THEN jumlah_suara END) AS total_suara_2,
               MAX(CASE WHEN rn = 3 THEN nama END) AS nama_3, 
               MAX(CASE WHEN rn = 3 THEN jumlah_suara END) AS total_suara_3,
			   MAX(CASE WHEN rn = 4 THEN nama END) AS nama_4, 
               MAX(CASE WHEN rn = 4 THEN jumlah_suara END) AS total_suara_4,
			   MAX(CASE WHEN rn = 5 THEN nama END) AS nama_5, 
               MAX(CASE WHEN rn = 5 THEN jumlah_suara END) AS total_suara_5,
			   MAX(CASE WHEN rn = 6 THEN nama END) AS nama_6, 
               MAX(CASE WHEN rn = 6 THEN jumlah_suara END) AS total_suara_6,
			   MAX(CASE WHEN rn = 7 THEN nama END) AS nama_7, 
               MAX(CASE WHEN rn = 7 THEN jumlah_suara END) AS total_suara_7,
			   MAX(CASE WHEN rn = 8 THEN nama END) AS nama_8, 
               MAX(CASE WHEN rn = 8 THEN jumlah_suara END) AS total_suara_8,
			   MAX(CASE WHEN rn = 9 THEN nama END) AS nama_9, 
               MAX(CASE WHEN rn = 9 THEN jumlah_suara END) AS total_suara_9,
			   MAX(CASE WHEN rn = 10 THEN nama END) AS nama_10, 
               MAX(CASE WHEN rn = 10 THEN jumlah_suara END) AS total_suara_10
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
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create the response data
	var result []map[string]interface{}
	for rows.Next() {
		var kelurahan string
		var suara1, suara2, suara3, suara4, suara5, suara6, suara7, suara8, suara9, suara10 int
		var nama1, nama2, nama3, nama4, nama5, nama6, nama7, nama8, nama9, nama10 string
		err = rows.Scan(&kelurahan, &nama1, &suara1, &nama2, &suara2, &nama3, &suara3, &nama4, &suara4, &nama5, &suara5, &nama6, &suara6, &nama7, &suara7, &nama8, &suara8, &nama9, &suara9, &nama10, &suara10)
		if err != nil {
			log.Fatal(err)
		}
		suara := []map[string]interface{}{
			{"nama_1": nama1, "total_suara": suara1},
			{"nama_2": nama2, "total_suara": suara2},
			{"nama_3": nama3, "total_suara": suara3},
			{"nama_4": nama4, "total_suara": suara4},
			{"nama_5": nama5, "total_suara": suara5},
			{"nama_6": nama6, "total_suara": suara6},
			{"nama_7": nama7, "total_suara": suara7},
			{"nama_8": nama8, "total_suara": suara8},
			{"nama_9": nama9, "total_suara": suara9},
			{"nama_10": nama10, "total_suara": suara10},
		}
		result = append(result, map[string]interface{}{"kelurahan": kelurahan, "suara": suara})
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Encode and return the response data
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatal(err)
	}
}
