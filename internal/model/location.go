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
	{ID: "morioka", Name: "盛岡市", Lat: 39.7021, Lon: 141.1545},
	{ID: "sendai", Name: "仙台市", Lat: 38.2682, Lon: 140.8694},
	{ID: "akita", Name: "秋田市", Lat: 39.7200, Lon: 140.1026},
	{ID: "yamagata", Name: "山形市", Lat: 38.2554, Lon: 140.3396},
	{ID: "fukushima", Name: "福島市", Lat: 37.7608, Lon: 140.4747},

	{ID: "mito", Name: "水戸市", Lat: 36.3659, Lon: 140.4712},
	{ID: "utsunomiya", Name: "宇都宮市", Lat: 36.5551, Lon: 139.8828},
	{ID: "maebashi", Name: "前橋市", Lat: 36.3895, Lon: 139.0634},
	{ID: "saitama", Name: "さいたま市", Lat: 35.8617, Lon: 139.6455},
	{ID: "chiba", Name: "千葉市", Lat: 35.6074, Lon: 140.1065},
	{ID: "tokyo", Name: "東京都心", Lat: 35.6812, Lon: 139.7671},
	{ID: "yokohama", Name: "横浜市", Lat: 35.4437, Lon: 139.6380},

	{ID: "niigata", Name: "新潟市", Lat: 37.9162, Lon: 139.0364},
	{ID: "toyama", Name: "富山市", Lat: 36.6953, Lon: 137.2113},
	{ID: "kanazawa", Name: "金沢市", Lat: 36.5613, Lon: 136.6562},
	{ID: "fukui", Name: "福井市", Lat: 36.0641, Lon: 136.2196},
	{ID: "kofu", Name: "甲府市", Lat: 35.6620, Lon: 138.5683},
	{ID: "nagano", Name: "長野市", Lat: 36.6486, Lon: 138.1948},
	{ID: "gifu", Name: "岐阜市", Lat: 35.4233, Lon: 136.7607},
	{ID: "shizuoka", Name: "静岡市", Lat: 34.9756, Lon: 138.3828},
	{ID: "nagoya", Name: "名古屋市", Lat: 35.1815, Lon: 136.9066},

	{ID: "tsu", Name: "津市", Lat: 34.7186, Lon: 136.5050},
	{ID: "otsu", Name: "大津市", Lat: 35.0179, Lon: 135.8546},
	{ID: "kyoto", Name: "京都市", Lat: 35.0116, Lon: 135.7681},
	{ID: "osaka", Name: "大阪市", Lat: 34.6937, Lon: 135.5023},
	{ID: "kobe", Name: "神戸市", Lat: 34.6901, Lon: 135.1956},
	{ID: "nara", Name: "奈良市", Lat: 34.6851, Lon: 135.8048},
	{ID: "wakayama", Name: "和歌山市", Lat: 34.2305, Lon: 135.1708},

	{ID: "tottori", Name: "鳥取市", Lat: 35.5011, Lon: 134.2351},
	{ID: "matsue", Name: "松江市", Lat: 35.4681, Lon: 133.0484},
	{ID: "okayama", Name: "岡山市", Lat: 34.6551, Lon: 133.9195},
	{ID: "hiroshima", Name: "広島市", Lat: 34.3853, Lon: 132.4553},
	{ID: "yamaguchi", Name: "山口市", Lat: 34.1785, Lon: 131.4737},

	{ID: "tokushima", Name: "徳島市", Lat: 34.0703, Lon: 134.5548},
	{ID: "takamatsu", Name: "高松市", Lat: 34.3428, Lon: 134.0466},
	{ID: "matsuyama", Name: "松山市", Lat: 33.8392, Lon: 132.7657},
	{ID: "kochi", Name: "高知市", Lat: 33.5597, Lon: 133.5311},

	{ID: "fukuoka", Name: "福岡市", Lat: 33.5902, Lon: 130.4017},
	{ID: "saga", Name: "佐賀市", Lat: 33.2635, Lon: 130.3009},
	{ID: "nagasaki", Name: "長崎市", Lat: 32.7503, Lon: 129.8777},
	{ID: "kumamoto", Name: "熊本市", Lat: 32.8031, Lon: 130.7079},
	{ID: "oita", Name: "大分市", Lat: 33.2396, Lon: 131.6093},
	{ID: "miyazaki", Name: "宮崎市", Lat: 31.9111, Lon: 131.4239},
	{ID: "kagoshima", Name: "鹿児島市", Lat: 31.5966, Lon: 130.5571},
	{ID: "naha", Name: "那覇市", Lat: 26.2124, Lon: 127.6809},
}
