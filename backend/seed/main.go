package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Ayah struct {
	Ayah   int    `json:"ayah"`
	Surah  int    `json:"surah"`
	Page   int    `json:"page"`
	Glyphs string `json:"glyphs"`
	Text   string `json:"text"`
}

func main() {
	bytes, err := os.ReadFile("../../qpc-glyphs/out/ayahs.json")
	if err != nil {
		panic(err)
	}

	Ayahs := map[string]Ayah{}
	err = json.Unmarshal(bytes, &Ayahs)
	if err != nil {
		panic(err)
	}

	AyahsList := []Ayah{}
	for _, ayah := range Ayahs {
		AyahsList = append(AyahsList, ayah)
	}

	sort.Slice(AyahsList, func(i, j int) bool {
		return AyahsList[i].Surah < AyahsList[j].Surah || (AyahsList[i].Surah == AyahsList[j].Surah && AyahsList[i].Ayah < AyahsList[j].Ayah)
	})

	sqlInsert := func(ayah Ayah) string {
		return fmt.Sprintf("INSERT INTO AYAHS (ayah, surah, page, glyphs, text) VALUES (%v, %v, %v, '%v', '%v');", ayah.Ayah, ayah.Surah, ayah.Page, ayah.Glyphs, ayah.Text)
	}

	inserts := ""
	for _, ayah := range AyahsList {
		// Save to database
		insertQuery := sqlInsert(ayah)
		inserts += insertQuery + "\n\n"
	}

	err = os.WriteFile("./sql/ayahs.sql", []byte(inserts), 0644)
}
