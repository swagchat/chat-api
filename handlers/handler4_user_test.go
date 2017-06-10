package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/fairway-corp/swagchat-api/utils"
)

func TestGetUsers2(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			out:            `(?m)^{"users":[{"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-1","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-2","name":"dennis-2","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"userId":"custom-user-id-3","name":"dennis-3","unreadCount":0,"metaData":{},"isPublic":false,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}]}$`,
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

func TestGetUser2(t *testing.T) {
	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createUserIds) != 9 {
		t.Fatalf("createUserIds length error \n[expected]%d\n[result  ]%d", 9, len(createUserIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			userId:         "custom-user-id-1",
			out:            `(?m)^{"userId":"[a-z0-9-]+","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"metaData":{"key":"value"},"isPublic":true,"created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","rooms":\[{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","metaData":{},"type":1,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","metaData":{},"type":2,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","metaData":{},"type":3,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","pictureUrl":"http://localhost/images/dennis_room.png","metaData":{},"type":1,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{},"type":1,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 1","pictureUrl":"http://localhost/images/dennis_room.png","informationUrl":"http://localhost/dennis_room","metaData":{"key":"value"},"type":1,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"},{"roomId":"[a-z0-9-]+","userId":"custom-user-id-1","name":"room name 2","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","metaData":{"key":"value"},"type":2,"lastMessage":"","lastMessageUpdated":"","created":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","modified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruUnreadCount":0,"ruMetaData":{},"ruCreated":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z","ruModified":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})Z"}\],"devices":\[{"userId":"custom-user-id-1","platform":2,"token":"jkl","notificationDeviceId":"jkl"}\]}$`,
			httpStatusCode: 200,
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
