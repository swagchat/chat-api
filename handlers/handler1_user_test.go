package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/fairway-corp/swagchat-api/utils"
)

var createUserIds []string

type userStruct struct {
	UserId string `json:"userId,omitempty"`
}

func TestPostUser(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			in: `
				{
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 2,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png"
				}
			`,
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 3,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png",
					"informationUrl": "http://localhost/dennis"
				}
			`,
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 4,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png",
					"informationUrl": "http://localhost/dennis",
					"metaData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 5,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png",
					"informationUrl": "http://localhost/dennis",
					"metaData": {"key": "value"},
					"isPublic": true
				}
			`,
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 6,
			in: `
				{
					"userId": "custom-user-id-1",
					"name": "dennis-1"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"dennis-1","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 7,
			in: `
				{
					"userId": "custom-user-id-2",
					"name": "dennis-2"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-2","name":"dennis-2","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 8,
			in: `
				{
					"userId": "custom-user-id-3",
					"name": "dennis-3"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-3","name":"dennis-3","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 9,
			in: `
				{
					"userId": "custom-user-id-for-delete",
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-for-delete","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"accessToken":"[a-zA-Z0-9-._~+/]+","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 10,
			in: `
				{
					"userId": "custom-user-id-1",
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"title":"An error occurred while creating user item.","status":500,"detail":".*","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
		{
			testNo: 11,
			in: `
				{
					"userId": "custom_id",
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create user item\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"userId","reason":"userId is invalid\. Available characters are alphabets, numbers and hyphens\."}\]}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 12,
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error\. \(Create user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 13,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error\. \(Create user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 14,
			in: `
				{
					"pictureUrl": "http://example.com/picture.png"
				}
			`,
			out:            `(?m)^{"title":"Request parameter error\. \(Create user item\)","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"name","reason":"name is required, but it's empty\."}\]}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		res, err := http.Post(ts.URL+"/"+utils.API_VERSION+"/users", "application/json", reader)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}

		if testRecord.httpStatusCode == 201 {
			user := &userStruct{}
			_ = json.Unmarshal(data, user)
			createUserIds = append(createUserIds, user.UserId)
		}
	}
}

func TestGetUsers(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			out:            `(?m)^{"users":[{"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-1","name":"dennis-1","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-2","name":"dennis-2","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-3","name":"dennis-3","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-for-delete","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}]}$`,
			httpStatusCode: 200,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users")
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestGetUser(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createUserIds) != 9 {
		t.Fatalf("createUserIds length error \n[expected]%d\n[result  ]%d", 9, len(createUserIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			userId:         createUserIds[0],
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			userId:         createUserIds[1],
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			userId:         createUserIds[2],
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         4,
			userId:         createUserIds[3],
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         5,
			userId:         createUserIds[4],
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},

		{
			testNo:         6,
			userId:         createUserIds[5],
			out:            `(?m)^{"userId":"custom-user-id-1","name":"dennis-1","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         7,
			userId:         createUserIds[6],
			out:            `(?m)^{"userId":"custom-user-id-2","name":"dennis-2","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         8,
			userId:         createUserIds[7],
			out:            `(?m)^{"userId":"custom-user-id-3","name":"dennis-3","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         9,
			userId:         createUserIds[8],
			out:            `(?m)^{"userId":"custom-user-id-for-delete","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         10,
			userId:         "not-exist-user-id",
			out:            ``,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.userId)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestPutUser(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			userId: "custom-user-id-1",
			in: `
				{
					"name": "Jeremy"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"Jeremy","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 2,
			userId: "custom-user-id-1",
			in: `
				{
					"pictureUrl": "http://localhost/images/jeremy.png"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 3,
			userId: "custom-user-id-1",
			in: `
				{
					"informationUrl": "http://localhost/jeremy"
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"metaData":{},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 4,
			userId: "custom-user-id-1",
			in: `
				{
					"metaData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"metaData":{"key":"value"},"isPublic":false,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 5,
			userId: "custom-user-id-1",
			in: `
				{
					"isPublic": true
				}
			`,
			out:            `(?m)^{"userId":"custom-user-id-1","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"isCanBlock":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}$`,
			httpStatusCode: 200,
		},
		{
			testNo: 6,
			userId: "custom-user-id-1",
			in: `
				{
					"unreadCount": -1
				}
			`,
			out:            `(?m)^{"title":"Json parse error\. \(Update user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 7,
			userId: "custom-user-id-1",
			in: `
			`,
			out:            `(?m)^{"title":"Json parse error\. \(Update user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			userId: "custom-user-id-1",
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json parse error\. \(Update user item\)","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 9,
			userId: "not-exist-user-id",
			in: `
				{
					"name": "Not Exist"
				}
			`,
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.userId, reader)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestDeleteUser(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			userId:         "custom-user-id-for-delete",
			out:            `(?m)^$`,
			httpStatusCode: 204,
		},
		{
			testNo:         2,
			userId:         "not-exist-user-id",
			out:            `(?m)^$`,
			httpStatusCode: 404,
		},
	}

	for _, testRecord := range testTable {
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.userId, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}
