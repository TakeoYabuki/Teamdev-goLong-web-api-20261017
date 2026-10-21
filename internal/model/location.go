package model

// Location は観測都市の情報です。
type Location struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// Locations は観測地点として使う固定都市データです。
// 座標は各都市の代表地点です。
var Locations = []Location{
	{ID: "sapporo", Name: "札幌市", Lat: 43.0618, Lon: 141.3545},
	{ID: "aomori", Name: "青森市", Lat: 40.8222, Lon: 140.7474},
	{ID: "sendai", Name: "仙台市", Lat: 38.2682, Lon: 140.8694},
	{ID: "tokyo", Name: "東京都心", Lat: 35.6812, Lon: 139.7671},
	{ID: "niigata", Name: "新潟市", Lat: 37.9162, Lon: 139.0364},
	{ID: "nagoya", Name: "名古屋市", Lat: 35.1815, Lon: 136.9066},
	{ID: "osaka", Name: "大阪市", Lat: 34.6937, Lon: 135.5023},
	{ID: "hiroshima", Name: "広島市", Lat: 34.3853, Lon: 132.4553},
	{ID: "kochi", Name: "高知市", Lat: 33.5597, Lon: 133.5311},
	{ID: "fukuoka", Name: "福岡市", Lat: 33.5902, Lon: 130.4017},
	{ID: "kagoshima", Name: "鹿児島市", Lat: 31.5966, Lon: 130.5571},
	{ID: "naha", Name: "那覇市", Lat: 26.2124, Lon: 127.6809},
}
