package amocrm_v4

import (
	"gorm.io/gorm"
	"net/http"
)

type Amo struct {
	Contact Ct
	Lead    Ld
	Task    Tsk
	Catalog Ctg
}

type authSettings struct {
	client            http.Client
	integrationID     string
	integrationSecret string
	endpoint          string
	redirectUri       string
	accessToken       string
	storage           *AuthAmoStorageConfig
}

type InitAmoConfig struct {
	Domain       string
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
}

type AuthAmoStorageConfig struct {
	DB        *gorm.DB
	TableName string
	AppName   string
}

type AmoAuthorizationDataStorage struct {
	Storage    string `json:"storage"`
	ConnectURI string `json:"connect_uri"`
}

func NewClient(initConfig *InitAmoConfig, storageConfig *AuthAmoStorageConfig) *Amo {
	err := createConnection(initConfig, storageConfig)
	if err != nil {
		panic(err)
	}
	return &Amo{}
}
