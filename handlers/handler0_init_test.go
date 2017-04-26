package handlers

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/utils"
)

type testRecord struct {
	testNo         int
	roomId         string
	userId         string
	messageId      string
	platform       string
	in             string
	out            string
	httpStatusCode int
}

func TestMain(m *testing.M) {
	datastoreProvider := datastore.GetProvider()
	err := datastoreProvider.Connect()
	if err != nil {
		log.Println(err.Error())
	}
	datastoreProvider.Init()
	/*
		Mux = bone.New().Prefix("/" + utils.API_VERSION)
		SetUserMux()
		SetRoomMux()
		SetRoomUserMux()
		SetMessageMux()
		SetAssetMux()
		SetDeviceMux()
	*/
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	StartServer(ctx)
	utils.IsTesting = true
	testRC := m.Run()
	err = datastoreProvider.DropDatabase()
	if err != nil {
		log.Println(err.Error())
	}
	os.Exit(testRC)
}

func TestIndex(t *testing.T) {
	testRecord := &testRecord{
		testNo:         1,
		out:            `(?m)^"Swagchat API version v[0-9]"$`,
		httpStatusCode: 200,
	}
	ts := httptest.NewServer(Mux)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/")

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

func TestNotFound(t *testing.T) {
	testRecord := &testRecord{
		testNo:         1,
		out:            ``,
		httpStatusCode: 404,
	}
	ts := httptest.NewServer(Mux)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/not-found")

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

//func BenchmarkPostRoom(b *testing.B) {
//	datastoreProvider := datastore.GetProvider()
//	err := datastoreProvider.Connect()
//	if err != nil {
//		log.Println(err.Error())
//	}
//	datastoreProvider.Init()
//	Mux = bone.New().Prefix("/" + utils.API_VERSION)
//	SetRoomMux()
//	ts := httptest.NewServer(Mux)
//	defer ts.Close()
//
//	for i := 0; i < b.N; i++ {
//		in := `
//				{
//					"name": "dennis room"
//				}
//			`
//		reader := strings.NewReader(in)
//		http.Post(ts.URL+"/"+utils.API_VERSION+"/rooms", "application/json", reader)
//	}
//	err = datastoreProvider.DropDatabase()
//	if err != nil {
//		log.Println(err.Error())
//	}
//}
