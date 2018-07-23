package idp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/pkg/errors"
)

type keycloakProvider struct {
	ctx          context.Context
	baseEndpoint string
}

type kcUser struct {
	ID               string   `json:"id"`
	Username         string   `json:"username"`
	CreatedTimestamp string   `json:"createdTimestamp"`
	RealmRoles       []string `json:"realmRoles"` // Read only
	Enabled          bool     `json:"enabled"`
}

type kcRoleMapping struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type Idp struct {
	Keycloak *Keycloak `json:"keycloak"`
}

type Keycloak struct {
	ManageUserClient *ManageUserClient `json:"manageUserClient"`
	GuestUserClient  *GuestUserClient  `json:"guestUserClient"`
	GuestRoleID      string            `json:"guestRoleId"`
}

type ManageUserClient struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type GuestUserClient struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (kp *keycloakProvider) Init() error {
	return nil
}

func (kp *keycloakProvider) Create() (string, string, error) {
	workspace := kp.ctx.Value(utils.CtxWorkspace)

	setting, err := datastore.Provider(kp.ctx).SelectOperatorSetting("1")
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	if setting == nil {
		err = errors.New("Operator setting is nil")
		logger.Error(err.Error())
		return "", "", err
	}

	var idp Idp
	json.Unmarshal(setting.Idp, &idp)

	// Create keycloak user
	kcUser := &kcUser{
		Username: utils.GenerateUUID(), // Finally overwritten with userID of keycloak
		Enabled:  true,
	}
	kcUserByte, err := json.Marshal(kcUser)

	endpoint := fmt.Sprintf("%s/auth/admin/realms/%s/users", kp.baseEndpoint, workspace)
	kcReq, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(kcUserByte),
	)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	kcReq.Header.Set("Content-Type", "application/json")

	manageUserToken, err := kp.clientToken(kp.ctx, idp.Keycloak.ManageUserClient.ClientID, idp.Keycloak.ManageUserClient.ClientSecret)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}

	kcReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", manageUserToken))
	client := &http.Client{}
	resp, err := client.Do(kcReq)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[keycloak]create uesr failure. HTTP Endpoint=%s, HTTP Status code=%d", endpoint, resp.StatusCode)
		logger.Error(err.Error())
		return "", "", err
	}
	defer resp.Body.Close()

	hl := resp.Header.Get("Location")
	userID := strings.Replace(hl, fmt.Sprintf("%s/", endpoint), "", 1)

	// Set guest role to keycloak user
	kcRoleMappings := []*kcRoleMapping{
		&kcRoleMapping{
			ID:   idp.Keycloak.GuestRoleID,
			Name: "guest",
		},
	}
	kcRoleMappingsByte, _ := json.Marshal(kcRoleMappings)

	endpoint = fmt.Sprintf("%s/auth/admin/realms/%s/users/%s/role-mappings/realm", kp.baseEndpoint, workspace, userID)
	kcReq, err = http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(kcRoleMappingsByte),
	)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	kcReq.Header.Set("Content-Type", "application/json")
	kcReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", manageUserToken))

	resp, err = client.Do(kcReq)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("[keycloak]set role failure. HTTP Endpoint=%s, HTTP Status code=%d", endpoint, resp.StatusCode)
		logger.Error(err.Error())
		return "", "", err
	}

	token, err := kp.clientToken(kp.ctx, idp.Keycloak.GuestUserClient.ClientID, idp.Keycloak.GuestUserClient.ClientSecret)
	if err != nil {
		logger.Error(err.Error())
		return "", "", err
	}

	return userID, token, nil
}

func (kp *keycloakProvider) GetToken() (string, error) {
	setting, err := datastore.Provider(kp.ctx).SelectOperatorSetting("1")
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	if setting == nil {
		err = errors.New("Operator setting is nil")
		logger.Error(err.Error())
		return "", err
	}
	var idp Idp
	json.Unmarshal(setting.Idp, &idp)

	token, err := kp.clientToken(kp.ctx, idp.Keycloak.GuestUserClient.ClientID, idp.Keycloak.GuestUserClient.ClientSecret)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return token, nil
}

func (kp *keycloakProvider) clientToken(ctx context.Context, clientID, clientSecret string) (string, error) {
	workspace := ctx.Value(utils.CtxWorkspace)
	if kp.baseEndpoint == "" || workspace == "" || clientID == "" || clientSecret == "" {
		return "", errors.Wrap(fmt.Errorf("[keycloak]Invalid params for create client accessToken. baseEndpoint=%s, realm=%s, ClientID=%s, ClientSecret=%s",
			kp.baseEndpoint, workspace, clientID, clientSecret), "")
	}

	endpoint := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", kp.baseEndpoint, workspace)

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

	basicAuth := basicAuth(clientID, clientSecret)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(fmt.Errorf("[keycloak]create client accessToken failure. HTTP Endpoint=%s, HTTP Status code=%d, ClientID=%s, ClientSecret=%s, Basic=%s", endpoint, resp.StatusCode, clientID, clientSecret, basicAuth), "")
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

func basicAuth(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
