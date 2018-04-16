package idp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type keycloakProvider struct {
	baseEndpoint string
}

type kcUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	CreatedTimestamp string `json:"createdTimestamp"`
	Enabled          bool   `json:"enabled"`
}

type kcToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefleshExpiresIn int    `json:"refresh_expires_in"`
	RefleshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
}

func (kp *keycloakProvider) Init() error {
	return nil
}

func (kp *keycloakProvider) Post(ctx context.Context) (*models.User, error) {
	kcUser := &kcUser{
		Username: utils.GenerateUUID(), // Finally overwritten with userID of keycloak
		Enabled:  true,
	}
	kcUserByte, err := json.Marshal(kcUser)

	realm := ctx.Value(utils.CtxRealm)
	endpoint := fmt.Sprintf("%s/auth/admin/realms/%s/users", kp.baseEndpoint, realm)
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(kcUserByte),
	)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	req.Header.Set("Content-Type", "application/json")

	token, err := kp.clientToken(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("[keycloak]create uesr failure. HTTP Status code=%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	hl := resp.Header.Get("Location")
	userID := strings.Replace(hl, fmt.Sprintf("%s/", endpoint), "", 1)

	user := &models.User{
		UserID: userID,
		Name:   userID,
	}

	user.BeforeInsertGuest()

	user, err = datastore.Provider(ctx).InsertUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	userToken, err := kp.userToken(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.AccessToken = userToken

	return user, nil
}

func (kp *keycloakProvider) Get(ctx context.Context, userID string) (*models.User, error) {
	userToken, err := kp.userToken(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	user, err := datastore.Provider(ctx).SelectUser(userID, true, true, true)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.AccessToken = userToken

	return user, nil
}

func (kp *keycloakProvider) clientToken(ctx context.Context) (string, error) {
	realm := ctx.Value(utils.CtxRealm)
	endpoint := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", kp.baseEndpoint, realm)

	values := url.Values{}
	values.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(
		"POST",
		endpoint,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	appClient, err := datastore.Provider(ctx).SelectLatestAppClientByName("manage-guests")
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	basicAuth := utils.BasicAuth(appClient.ClientID, appClient.ClientSecret)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	kcToken := &kcToken{}
	err = json.Unmarshal(body, &kcToken)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	return kcToken.AccessToken, nil
}

func (kp *keycloakProvider) userToken(ctx context.Context, userID string) (string, error) {
	realm := ctx.Value(utils.CtxRealm)
	endpoint := fmt.Sprintf("%s/auth/admin/realms/%s/users/%s/impersonation", kp.baseEndpoint, realm, userID)

	req, err := http.NewRequest(
		"POST",
		endpoint,
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	token, err := kp.clientToken(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[keycloak]user not found")
	}
	defer resp.Body.Close()

	var userToken string
	cookies := resp.Cookies()
	for _, v := range cookies {
		if v.Name == "KEYCLOAK_IDENTITY" && v.Value != "" {
			userToken = v.Value
			break
		}
	}

	return userToken, nil
}
