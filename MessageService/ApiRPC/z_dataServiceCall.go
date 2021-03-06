package ApiRPC

import (
	"../RpcClientPbs/mongoPb"
	"../RpcClientPbs/mysqlPb"
)

func SaveDelayMessage(userId int64, message []byte) error {
	// todo test code used in separate development, need remove later
	return nil

	// code to actually use
	param := &mongoPb.IdAndMessage{Id: userId, Message: message}
	_, err := GetMongoDataClient().PutSaveDelayMessage(getTimeOutCtx(3), param)
	return err
}

func GetDelayMessages(userId int64) ([][]byte, error) {
	// todo test code used in separate development, need remove later
	testDelayMessages := [][]byte{
		[]byte("<test delay message 1>"),
		[]byte("<test delay message 2>"),
		[]byte("<test delay message 3>"),
		[]byte("<test delay message 4>"),
	}
	return testDelayMessages, nil

	// code to actually use
	param := &mongoPb.Id{Value: userId}
	data, err := GetMongoDataClient().GetDelayMessage(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.MessageList, nil
}

func GetUserFriends(userId int64) ([]int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{userId - 1, userId + 1, userId + 2}, nil

	// code to actually use
	param := &mysqlPb.Id{Value: userId}
	data, err := GetMySQLDataClient().GetEffectiveFriendsIdListByIdPlus(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.Data, nil
}

func GetUserBlacklist(userId int64) ([]int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{userId + 5}, nil

	// code to actually use
	param := &mysqlPb.Id{Value: userId}
	data, err := GetMySQLDataClient().GetBlacklistFriendsIdListByIdPlus(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.Data, nil
}

func GetGroupChatUsers(groupId int64) ([]int64, error) {
	// todo test code used in separate development, need remove later
	return []int64{0, 1, 2, 3}, nil

	// code to actually use
	param := &mysqlPb.IdAndIsDelete{Id: groupId}
	data, err := GetMySQLDataClient().GetUserIdListOfGroupChat(getTimeOutCtx(3), param)
	if nil != err {
		return nil, err
	}
	return data.Data, nil

}

func SaveUserChatHistory(joinId string, message []byte) error {
	// todo test code used in separate development, need remove later
	return nil

	// code to actually use
	param := &mongoPb.JoinIdAndMessage{JoinId: joinId, Message: message}
	_, err := GetMongoDataClient().PutSaveUserChatHistory(getTimeOutCtx(3), param)
	return err
}

func SaveGroupChatHistory(groupId int64, message []byte) error {
	// todo test code used in separate development, need remove later
	return nil

	param := &mongoPb.IdAndMessage{Id: groupId, Message: message}
	_, err := GetMongoDataClient().PutSaveGroupChatHistory(getTimeOutCtx(3), param)
	return err
}

func GetSubscriptionInfo(subsId int64) (managerId int64, fans []int64, err error) {
	// todo test code used in separate development, need remove later
	return 1, []int64{1, 2, 3, 4, 5}, nil

	// code to actually use
	param := &mysqlPb.IdAndIsDelete{Id: subsId}
	subscription, err := GetMySQLDataClient().GetOneSubscriptionById(getTimeOutCtx(3), param)
	if nil != err {
		return 0, nil, err
	}
	data, err := GetMySQLDataClient().GetUserIdListOfSubscription(getTimeOutCtx(3), param)

	return subscription.ManagerId, data.Data, nil

}

func SaveSubscriptionMessageHistory(subsId int64, message []byte) error {
	// todo test code used in separate development, need remove later
	return nil

	// code to actually use
	param := &mongoPb.IdAndMessage{Id: subsId, Message: message}
	_, err := GetMongoDataClient().PutSaveSubscriptionHistory(getTimeOutCtx(3), param)
	return err
}
