package model

type N2YOResponse struct {
	Info struct {
		SatID       int `json:"satid"`
		PassesCount int `json:"passescount"`
	} `json:"info"`

	Passes []N2YOPass `json:"passes"`
}

type N2YOPass struct {
	StartUTC       int64    `json:"startUTC"`
	StartAz        float64  `json:"startAz"`
	StartAzCompass string   `json:"startAzCompass"`
	StartEl        float64  `json:"startEl"`
	MaxUTC         int64    `json:"maxUTC"`
	MaxAz          float64  `json:"maxAz"`
	MaxAzCompass   string   `json:"maxAzCompass"`
	MaxEl          float64  `json:"maxEl"`
	EndUTC         int64    `json:"endUTC"`
	EndAz          float64  `json:"endAz"`
	EndAzCompass   string   `json:"endAzCompass"`
	EndEl          float64  `json:"endEl"`
	Duration       int      `json:"duration"`
	Mag            *float64 `json:"mag"`
}
