package DataLayer

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDBURI = "mongodb://localhost:27017"
)

var (
	MongoClient *mongo.Client
	MsgCenterDB *mongo.Database

	WaitSendMsgColl  *mongo.Collection
	UserFriendsColl  *mongo.Collection
	UserBlackColl    *mongo.Collection
	UserGroupChat    *mongo.Collection
	UserSubscription *mongo.Collection

	GroupChats    *mongo.Collection
	Subscriptions *mongo.Collection
)

func init() {
	var err error
	MongoClient, err = mongo.Connect(getTimeOutCtx(10), options.Client().ApplyURI(MongoDBURI))
	if nil != err {
		log.Fatal(err)
	}
	err = MongoClient.Ping(getTimeOutCtx(3), readpref.Primary())
	if nil != err {
		log.Fatal(err)
	}
	MsgCenterDB = MongoClient.Database("MsgCenter")
	WaitSendMsgColl = MsgCenterDB.Collection("WaitSendMsg")
	UserFriendsColl = MsgCenterDB.Collection("UserFriends")
	UserBlackColl = MsgCenterDB.Collection("UserBlackList")

	GroupChats = MsgCenterDB.Collection("GroupChats")
	Subscriptions = MsgCenterDB.Collection("Subscriptions")

}

func getTimeOutCtx(expire time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), expire*time.Second)
	return ctx
}

type TempWaitSendMsg struct {
	Id      int64    `bson:"_id"`
	Message [][]byte `bson:"message"`
}

//Save the message which sent failed because the target user is offline.
func MongoSaveWaitSendMessage(id int64, data []byte) error {
	_, err := WaitSendMsgColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"message": data}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: save WaitSendMessage fail for user(%d), error detail: %s", id, err.Error())
		return err
	} else {
		log.Printf("WaitSendMessage: save an message for user(%d)", id)
		return nil
	}
}

// Query the messages should be sent to current user.
func MongoQueryWaitSendMessage(id int64) ([][]byte, error) {
	temp := new(TempWaitSendMsg)
	err := WaitSendMsgColl.FindOneAndDelete(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query WaitSendMessage fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.Message, nil
}

type TempFriends struct {
	Id      int64   `bson:"_id"`
	Friends []int64 `bson:"friends"`
}

// Add a friend's id for the current user.
func MongoAddFriendId(srcId, dstId int64) error {
	_, err := UserFriendsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$addToSet": bson.M{"friends": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add friends fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

// Query the id of friends of the current user
func MongoQueryFriendsId(id int64) ([]int64, error) {
	temp := new(TempFriends)
	err := UserFriendsColl.FindOne(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query friends fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.Friends, nil
}

// Remove a friend's id for the current user.
func MongoDelFriendId(srcId, dstId int64) error {
	_, err := UserFriendsColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$pull": bson.M{"friends": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: remove a friend fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

type TempBlackList struct {
	Id        int64   `bson:"_id"`
	BlackList []int64 `bson:"blacklist"`
}

// Add a user'id into the blacklist of current user
func MongoBlackListAdd(srcId, dstId int64) error {
	_, err := UserBlackColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$addToSet": bson.M{"blacklist": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add friends fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

// Query the id of the friends who marked black by current user
func MongoQueryBlackList(id int64) ([]int64, error) {
	temp := new(TempBlackList)
	err := UserBlackColl.FindOne(getTimeOutCtx(3), bson.M{"_id": id}).Decode(temp)
	if nil != err {
		log.Printf("Error: query blacklist fail for user(%d), error detail: %s", id, err.Error())
		return nil, err
	}
	return temp.BlackList, nil
}

// Move a user's id out from the blacklist of current user
func MongoBlackListDel(srcId, dstId int64) error {
	_, err := UserBlackColl.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": srcId},
		bson.M{"$pull": bson.M{"blacklist": dstId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: move a friend out from blacklist fail for user(%d), error detail: %s", srcId, err.Error())
		return err
	}
	return nil
}

type TempGroupChat struct {
	Id      int64   `bson:"_id"`
	UsersId []int64 `bson:"users_id"`
}

// Add a group's id into the groups_id of current user
func MongoGroupChatAddUser(groupId, userId int64) error {
	_, err := GroupChats.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": groupId},
		bson.M{"$addToSet": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: add user fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return err
	}
	return nil
}

// Query the user's id of the group
func MongoQueryGroupChatUsers(groupId int64) ([]int64, error) {
	temp := new(TempGroupChat)
	err := GroupChats.FindOne(getTimeOutCtx(3), bson.M{"_id": groupId}).Decode(temp)
	if nil != err {
		log.Printf("Error: query users fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return nil, err
	}
	return temp.UsersId, nil
}

// Query the all groups information
func MongoQueryGroupChatAll() ([]TempGroupChat, error) {
	ctx := getTimeOutCtx(30)
	curs, err := GroupChats.Find(ctx, bson.D{})
	if nil != err {
		log.Printf("Error: query all group chat information fail")
		return nil, err
	}
	defer curs.Close(ctx)
	data := make([]TempGroupChat, 0)
	for curs.Next(ctx) {
		temp := new(TempGroupChat)
		err := curs.Decode(temp)
		if nil != err {
			log.Printf("Error: query all group chat error, detail: %s", err.Error())
			continue
		}
		data = append(data, *temp)
	}
	return data, nil
}

// Move a user's id out from a group chat
func MongoGroupChatDelUser(groupId, userId int64) error {
	_, err := GroupChats.UpdateOne(getTimeOutCtx(3),
		bson.M{"_id": groupId},
		bson.M{"$pull": bson.M{"users_id": userId}},
		options.Update().SetUpsert(true))
	if nil != err {
		log.Printf("Error: remove user fail for groupChat(%d), error detail: %s", groupId, err.Error())
		return err
	}
	return nil
}
