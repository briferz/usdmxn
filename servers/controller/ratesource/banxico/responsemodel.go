package banxico

type responseModel struct {
	BMX Bmx `json:"bmx"`
}

type Bmx struct {
	Series []Serie `json:"series"`
}

type Serie struct {
	IdSerie string `json:"idSerie"`
	Titulo  string `json:"titulo"`
	Datos []Dato `json:"datos"`
}

type Dato struct {
	Fecha string `json:"fecha"`
	Dato  string `json:"dato"`
}
