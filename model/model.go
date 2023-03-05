package model

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json"total_page"`
}

type Employee struct {
	Id          string `form:"id" json:"id"`
	Nama        string `form:"nama" json:"nama"`
	Kelurahan   string `form:"kelurahan" json:"kelurahan"`
	Kecamatan   string `form:"kecamatan" json":"kecamatan"`
	User        string `form:"user" json":"user"`
	TPS         string `form:"tps" json:"tps"`
	Gambar      string `form:"gambar" json:"gambar"`
	JumlahSuara int    `form:"jumlah_suara" json":"jumlah_suara"`
	CreatedAt   string `form:"created_at" json:"created_at"`
}
type Kandidat struct {
	Id     string `form:"id" json:"id"`
	Nama   string `form:"nama" json:"nama"`
	UserID int    `form:"user_id" json:"user_id"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type ResponseKandidat struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Kandidat
}

type Kecamatan struct {
	Id        int    `json:"id"`
	Kecamatan string `json:"kecamatan"`
}
type Kelurahan struct {
	Id        int    `json:"id"`
	Kelurahan string `json:"kelurahan"`
	Kecamatan string `json:"kecamatan"`
}

type KecamatanResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Kecamatan
}

type KelurahanResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Kelurahan
}
