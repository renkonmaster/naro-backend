package main

import (
	"database/sql"
	"maps"
	"naro-backend/handler"
	"testing"
)

func Test_sumPopulationByCountryCode(t *testing.T) {
	// ここにテストケースを書いていく
	cases := []struct {
		name   string           // テストケースの名前
		cities []handler.City           // テストケースの入力
		want   map[string]int64 // 期待される結果
	}{
		{
			name:   "empty input",
			cities: []handler.City{},
			want:   map[string]int64{},
		},
		{
			name:	"only one input",
			cities: []handler.City{
				{
					ID:          1,
					Name:        sql.NullString{String: "Tokyo", Valid: true},
					CountryCode: sql.NullString{String: "JP", Valid: true},
					District:    sql.NullString{String: "Kanto", Valid: true},
					Population:  sql.NullInt64{Int64: 13929286, Valid: true},
				},
			},
			want: 	map[string]int64{
				"JP": 13929286, // Tokyoの人口
			},
		},
		{
			name:	"multiple inputs with same country code",
			cities: []handler.City{
				{
					ID:          1,
					Name:        sql.NullString{String: "Tokyo", Valid: true},
					CountryCode: sql.NullString{String: "JP", Valid: true},
					District:    sql.NullString{String: "Kanto", Valid: true},
					Population:  sql.NullInt64{Int64: 13929286, Valid: true},
				},
				{
					ID:          2,
					Name:        sql.NullString{String: "Osaka", Valid: true},
					CountryCode: sql.NullString{String: "JP", Valid: true},
					District:    sql.NullString{String: "Kansai", Valid: true},
					Population:  sql.NullInt64{Int64: 8839469, Valid: true},
				},
				{
					ID:          3,
					Name:        sql.NullString{String: "Kyoto", Valid: true},
					CountryCode: sql.NullString{String: "JP", Valid: true},
					District:    sql.NullString{String: "Kansai", Valid: true},
					Population:  sql.NullInt64{Int64: 1474570, Valid: true},
				},
			},
			want: map[string]int64{
				"JP":  13929286 + 8839469 + 1474570, // Tokyo + Osaka + Kyoto
			},
		},
		{
			name:   "input with city.CountryCode.Valid = false",
			cities: []handler.City{
				{
					ID:          1,
					Name:        sql.NullString{String: "Tokyo", Valid: true},
					CountryCode: sql.NullString{String: "", Valid: false}, // 無効
					District:    sql.NullString{String: "Kanto", Valid: true},
					Population:  sql.NullInt64{Int64: 13929286, Valid: true},
				},
				{
					ID:          2,
					Name:        sql.NullString{String: "Osaka", Valid: true},
					CountryCode: sql.NullString{String: "JP", Valid: true},
					District:    sql.NullString{String: "Kansai", Valid: true},
					Population:  sql.NullInt64{Int64: 8839469, Valid: true},
				},
			},
			want: map[string]int64{
				"JP": 8839469, // Osakaのみ
			},
		},
	}
	for _, tt := range cases {
		// サブテストの実行
		t.Run(tt.name, func(t *testing.T) {
			got := sumPopulationByCountryCode(tt.cities)
			if !maps.Equal(got, tt.want) {
				t.Errorf("sumPopulationByCountryCode(%v) = %v, want %v", tt.cities, got, tt.want)
			}
		})
	}
}