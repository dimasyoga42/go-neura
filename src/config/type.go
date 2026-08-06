package config

type Xtal struct {
	Name string `json:"name"`
	Tipe string `json:"tipe"`
	Stat string `json:"stat"`
	Rute string `json:"rute"`
	Max  string `json:"max"`
}

type Regis struct {
	Name      string `json:"name"`
	Deskripsi string `json:"effect"`
	Lv        string `json:"max_lv"`
	Lv_stode  string `json:"levels_studied"`
}
type Trait struct {
	Name string `json:"name"`
	Deskripsi string `json:"stat_effect"`
}
type ApiRespon struct {
	Data any     `json:"data"`
	Time float64 `json:"response_time_ms"`
}
