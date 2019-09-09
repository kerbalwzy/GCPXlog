package Config

const (
	MySQLDataRPCServerAddress = "0.0.0.0:23331"
	MongoDataRPCServerAddress = "0.0.0.0:23332"

	UserDbMySQLURI = "root:123456@tcp(127.0.0.1:3306)/IMUserCenter?charset=utf8" +
		"&parseTime=true"

	MsgDbMongoURI = "mongodb://localhost:27017"

	TimeDisplayFormat = "2006-01-02 15:04:05"

	DataLayerSrvCAPem       = "/Users/wzy/GitPrograms/PrivateIM/DataService/CATSL/ca.pem"
	DataLayerSrvCAServerPem = "/Users/wzy/GitPrograms/PrivateIM/DataService/CATSL/server/server.pem"
	DataLayerSrvCAServerKey = "/Users/wzy/GitPrograms/PrivateIM/DataService/CATSL/server/server.key"
	DataLayerSrvCAClientPem = "/Users/wzy/GitPrograms/PrivateIM/DataService/CATSL/client/client.pem"
	DataLayerSrvCAClientKey = "/Users/wzy/GitPrograms/PrivateIM/DataService/CATSL/client/client.key"
)
