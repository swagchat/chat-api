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

	"github.com/betchi/go-gimei"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
)

type keycloakProvider struct {
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

func (kp *keycloakProvider) Init() error {
	return nil
}

func (kp *keycloakProvider) Post(ctx context.Context) (*models.User, error) {
	workspace := ctx.Value(utils.CtxWorkspace)
	gimei := gimei.NewName()
	name := fmt.Sprintf("%s(%s)(仮)", gimei.Kanji(), gimei.Katakana())

	setting, err := datastore.Provider(ctx).SelectLatestSetting()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if setting == nil {
		return nil, errors.Wrap(err, "setting values is nil")
	}
	var settingValues models.SettingValues
	json.Unmarshal(setting.Values, &settingValues)

	// Create keycloak user
	kcUser := &kcUser{
		Username: utils.GenerateUUID(), // Finally overwritten with userID of keycloak
		Enabled:  true,
	}
	kcUserByte, err := json.Marshal(kcUser)

	endpoint := fmt.Sprintf("%s/auth/admin/realms/%s/users", kp.baseEndpoint, workspace)
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(kcUserByte),
	)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	req.Header.Set("Content-Type", "application/json")

	manageUserToken, err := kp.clientToken(ctx, settingValues.Keycloak.ManageUserClient.ClientID, settingValues.Keycloak.ManageUserClient.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", manageUserToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Wrap(fmt.Errorf("[keycloak]create uesr failure. HTTP Endpoint=%s, HTTP Status code=%d", endpoint, resp.StatusCode), "")
	}
	defer resp.Body.Close()

	hl := resp.Header.Get("Location")
	userID := strings.Replace(hl, fmt.Sprintf("%s/", endpoint), "", 1)

	// Set guest role to keycloak user
	kcRoleMappings := []*kcRoleMapping{
		&kcRoleMapping{
			ID:   settingValues.Keycloak.GuestRoleID,
			Name: "guest",
		},
	}
	kcRoleMappingsByte, _ := json.Marshal(kcRoleMappings)

	endpoint = fmt.Sprintf("%s/auth/admin/realms/%s/users/%s/role-mappings/realm", kp.baseEndpoint, workspace, userID)
	req, err = http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(kcRoleMappingsByte),
	)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", manageUserToken))

	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusNoContent {
		return nil, errors.Wrap(fmt.Errorf("[keycloak]set role failure. HTTP Endpoint=%s, HTTP Status code=%d", endpoint, resp.StatusCode), "")
	}

	// Create user
	user := &models.User{
		UserID: userID,
		Name:   name,
	}
	user.BeforeInsertGuest()

	general := &protobuf.UserRole{
		UserID: user.UserID,
		RoleID: utils.RoleGeneral,
	}
	guest := &protobuf.UserRole{
		UserID: user.UserID,
		RoleID: utils.RoleGuest,
	}
	roles := []*protobuf.UserRole{general, guest}

	user, err = datastore.Provider(ctx).InsertUser(user, roles)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	guestUserToken, err := kp.clientToken(ctx, settingValues.Keycloak.GuestUserClient.ClientID, settingValues.Keycloak.GuestUserClient.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.AccessToken = guestUserToken

	return user, nil
}

func (kp *keycloakProvider) Get(ctx context.Context, userID string) (*models.User, error) {
	setting, err := datastore.Provider(ctx).SelectLatestSetting()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if setting == nil {
		return nil, errors.Wrap(err, "setting values is nil")
	}
	var settingValues models.SettingValues
	json.Unmarshal(setting.Values, &settingValues)

	user, err := datastore.Provider(ctx).SelectUser(userID, datastore.WithBlocks(true), datastore.WithDevices(true), datastore.WithRooms(true))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	guestUserToken, err := kp.clientToken(ctx, settingValues.Keycloak.GuestUserClient.ClientID, settingValues.Keycloak.GuestUserClient.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.AccessToken = guestUserToken

	return user, nil
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

	basicAuth := utils.BasicAuth(clientID, clientSecret)
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
