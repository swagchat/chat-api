package models

type Asset struct {
	AssetId   string `json:"assetId"`
	SourceUrl string `json:"sourceUrl"`
	Mime      string `json:"mime"`
}
