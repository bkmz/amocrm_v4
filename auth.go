package amocrm_v4

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type amoAuthRequestType string

const (
	amoAuthorizationAuthCode     amoAuthRequestType = "authorization_code"
	amoAuthorizationRefreshToken amoAuthRequestType = "refresh_token"
)

type AuthorizationData struct {
	gorm.Model
	AppName      string    `gorm:"column:app_name"`
	RefreshToken string    `gorm:"column:refresh_token"`
	ExpiresIn    time.Time `gorm:"column:expires_in"`
}

type authResp struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authRequest struct {
	ClientId     string             `json:"client_id"`
	ClientSecret string             `json:"client_secret"`
	GrantType    amoAuthRequestType `json:"grant_type"`
	Code         string             `json:"code,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty"`
	RedirectUri  string             `json:"redirect_uri"`
}

var client authSettings

func (a *authSettings) getUrl(path string) string {
	return fmt.Sprintf("%s%s", a.endpoint, path)
}

func createConnection(init *InitAmoConfig, storage *AuthAmoStorageConfig) error {
	// Проверяем наличие записи в таблице с данными об авторизации в АМО
	client.endpoint = fmt.Sprintf("https://%s.amocrm.com", init.Domain)
	client.integrationID = init.ClientID
	client.integrationSecret = init.ClientSecret
	client.client = http.Client{}
	client.storage = storage

	err := client.open(init.Code)
	if err != nil {
		return err
	}
	// запускаем фоновую задачу для обновления access_token
	go client.refresher()

	return nil
}

func (a *authSettings) open(authCode string) error {

	amoAuthorizationData := AuthorizationData{}
	result := a.storage.DB.Table(a.storage.TableName).Where("app_name = ?", a.storage.AppName).First(&amoAuthorizationData)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Errorf("Ошибка при получении данных об авторизации в АМО: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		ret := authResp{}
		opts := requestOpts{
			Method: http.MethodPost,
			Path:   "/oauth2/access_token",
			DataParameters: authRequest{
				ClientId:     a.integrationID,
				ClientSecret: a.integrationSecret,
				GrantType:    amoAuthorizationAuthCode,
				Code:         authCode,
				RedirectUri:  a.redirectUri,
			},
			Ret: ret,
		}

		err := httpRequest(opts)
		if err != nil {
			return err
		}

		client.accessToken = ret.AccessToken
		amoAuthorizationData = AuthorizationData{
			AppName:      a.storage.AppName,
			RefreshToken: ret.RefreshToken,
			ExpiresIn:    time.Now().Add(time.Duration(ret.ExpiresIn)*time.Second - 1*time.Minute),
		}

		err = a.storage.DB.Table(a.storage.TableName).Create(&amoAuthorizationData).Error
		if err != nil {
			return fmt.Errorf("ошибка при создании записи в БД об авторизации в АМО: %v", err)
		}
	} else {
		var ret = authResp{}
		opts := requestOpts{
			Method: http.MethodPost,
			Path:   "/oauth2/access_token",
			DataParameters: &authRequest{
				ClientId:     a.integrationID,
				ClientSecret: a.integrationSecret,
				GrantType:    amoAuthorizationRefreshToken,
				RefreshToken: amoAuthorizationData.RefreshToken,
				RedirectUri:  a.redirectUri,
			},
			Ret: ret,
		}

		err := httpRequest(opts)
		if err != nil {
			return fmt.Errorf("ошибка получения нового access_token: %v", err)
		}

		client.accessToken = ret.AccessToken
		// сохраняем новый refresh token в БД
		err = updateAuthDataInDB(
			a.storage.DB, a.storage.TableName, a.storage.AppName, ret.RefreshToken,
			time.Now().Add(time.Duration(ret.ExpiresIn)*time.Second-1*time.Minute),
		)
		if err != nil {
			return fmt.Errorf("ошибка при сохранении нового refresh_token в БД: %v", err)
		}

	}

	return nil
}

func (a *authSettings) refresher() {
	ticker := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-ticker.C:
			auth, err := getAuthDataFromDB(a.storage.DB, a.storage.TableName, a.storage.AppName)
			if err != nil {
				log.Errorf("Ошибка при получении данных об авторизации в АМО: %v", err)
				continue
			}

			if auth.ExpiresIn.Before(time.Now().Add(-5 * time.Minute)) {
				var ret = authResp{}
				opts := requestOpts{
					Method: http.MethodPost,
					Path:   "/oauth2/access_token",
					DataParameters: &authRequest{
						ClientId:     a.integrationID,
						ClientSecret: a.integrationSecret,
						GrantType:    amoAuthorizationRefreshToken,
						RefreshToken: auth.RefreshToken,
						RedirectUri:  a.redirectUri,
					},
					Ret: ret,
				}
				err = httpRequest(opts)
				if err != nil {
					log.Errorf("Ошибка при обновлении авторизационного токена: %v", err)
					continue
				}

				client.accessToken = ret.AccessToken
				// сохраняем новый refresh token в БД
				err = updateAuthDataInDB(
					a.storage.DB, a.storage.TableName, a.storage.AppName, ret.RefreshToken,
					time.Now().Add(time.Duration(ret.ExpiresIn)*time.Second-1*time.Minute),
				)
				if err != nil {
					log.Errorf("Ошибка при сохранении нового refresh token в БД: %v", err)
				}
			} else {
				log.Infof("Авторизационный токен истекает %v, обновление не требуется", auth.ExpiresIn)
			}
		}
	}
}

func getAuthDataFromDB(db *gorm.DB, table string, appName string) (*AuthorizationData, error) {
	amoAuthorizationData := &AuthorizationData{}

	result := db.Table(table).Where("app_name = ?", appName).First(&amoAuthorizationData)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("не найдена запись об авторизации в АМО")
	}

	return amoAuthorizationData, nil
}

func updateAuthDataInDB(db *gorm.DB, table string, appName string, refreshToken string, expiresIn time.Time) error {
	return db.Table(table).Where("app_name = ?", appName).
		Updates(map[string]interface{}{
			"refresh_token": refreshToken,
			"expires_in":    expiresIn,
		}).Error
}
