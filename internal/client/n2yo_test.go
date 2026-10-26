package client

import (
	"encoding/json"
	"testing"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

func TestN2YOResponseDecode(t *testing.T) {
	// N2YOから返されるJSONを再現
	jsonData := `{
    "info": {
        "satid": 25544,
        "satname": "SPACE STATION",
        "transactionscount": 1,
        "passescount": 1
    },
    "passes": [
        {
            "startAz": 274.66,
            "startAzCompass": "W",
            "startEl": 12.02,
            "startUTC": 1792523360,
            "maxAz": 333.04,
            "maxAzCompass": "NNW",
            "maxEl": 13.72,
            "maxUTC": 1792523635,
            "endAz": 34.33,
            "endAzCompass": "NE",
            "endEl": 0.27,
            "endUTC": 1792523920,
            "mag": 0.2,
            "duration": 210,
            "startVisibility": 1792523710
        }
    ]
}`

	var data model.N2YOResponse

	// JSONを構造体に変換
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		t.Fatalf("JSONデコード失敗: %v", err)
	}

	// 通過件数を確認
	if len(data.Passes) != 1 {
		t.Fatalf("通過件数が異なります: got=%d, want=1", len(data.Passes))
	}

	p := data.Passes[0]

	// 必要な12項目を確認
	if p.StartUTC != 1792523360 {
		t.Errorf("StartUTC: got=%d", p.StartUTC)
	}
	if p.StartAz != 274.66 {
		t.Errorf("StartAz: got=%f", p.StartAz)
	}
	if p.StartAzCompass != "W" {
		t.Errorf("StartAzCompass: got=%s", p.StartAzCompass)
	}
	if p.MaxUTC != 1792523635 {
		t.Errorf("MaxUTC: got=%d", p.MaxUTC)
	}
	if p.MaxAz != 333.04 {
		t.Errorf("MaxAz: got=%f", p.MaxAz)
	}
	if p.MaxAzCompass != "NNW" {
		t.Errorf("MaxAzCompass: got=%s", p.MaxAzCompass)
	}
	if p.MaxEl != 13.72 {
		t.Errorf("MaxEl: got=%f", p.MaxEl)
	}
	if p.EndUTC != 1792523920 {
		t.Errorf("EndUTC: got=%d", p.EndUTC)
	}
	if p.EndAz != 34.33 {
		t.Errorf("EndAz: got=%f", p.EndAz)
	}
	if p.EndAzCompass != "NE" {
		t.Errorf("EndAzCompass: got=%s", p.EndAzCompass)
	}
	if p.Duration != 210 {
		t.Errorf("Duration: got=%d", p.Duration)
	}
	if p.Mag != 0.2 {
		t.Errorf("Mag: got=%f", p.Mag)
	}

	if p.StartEl != 12.02 {
		t.Errorf("StartEl: got %v, want 12.02", p.StartEl)
	}

	if p.EndEl != 0.27 {
		t.Errorf("EndEl: got %v, want 0.27", p.EndEl)
	}
}
