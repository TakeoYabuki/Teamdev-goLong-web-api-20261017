package handler

import (
	"net/url"
	"testing"
)

func TestCheckQuery(t *testing.T) {
	allowed := []string{"q", "page", "limit"}

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"パラメータなしは通す", "", false},
		{"知っている名前だけなら通す", "q=tokyo&page=2", false},
		{"知らない名前は断る", "city=tokyo", true},
		{"同じ名前が2回なら断る", "q=a&q=b", true},
		{"値が空なら断る", "page=", true},
		{"イコールなしも空として断る", "page", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatalf("クエリを読めません: %v", err)
			}

			err = checkQuery(query, allowed...)

			if tt.wantErr && err == nil {
				t.Error("断るはずが、通してしまいました")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("通すはずが、断りました: %v", err)
			}
		})
	}
}
