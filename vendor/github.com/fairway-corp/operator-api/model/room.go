package model

import (
	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type Room struct {
	scpb.Room
	MetaData utils.JSONText
	Users    []*UserForRoom
}

type UserForRoom struct {
	scpb.UserForRoom
}
