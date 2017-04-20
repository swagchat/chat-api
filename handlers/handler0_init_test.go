package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/utils"
	"github.com/go-zoo/bone"
)

type testRecord struct {
	testNo         int
	roomId         string
	userId         string
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
	Mux = bone.New().Prefix("/" + utils.API_VERSION)
	SetUserMux()
	SetRoomMux()
	SetRoomUserMux()
	SetMessageMux()
	SetAssetMux()
	SetDeviceMux()
	testRC := m.Run()
	err = datastoreProvider.DropDatabase()
	if err != nil {
		log.Println(err.Error())
	}
	os.Exit(testRC)
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
