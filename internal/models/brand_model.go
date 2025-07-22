package models

type BrandModel struct {
	Id        string `json:"id"`
	BrandName string `json:"brand_name"`
	ModelName string `json:"model_name"`
	TypeName  string `json:"type_name"`
}
