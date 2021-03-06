package Config

const (
	// gRPC config
	MySQLDataRPCServerAddress   = "0.0.0.0:23331"
	MongoDataRPCServerAddress   = "0.0.0.0:23332"
	UserAuthRPCServerAddress    = "0.0.0.0:11111"
	MsgTransferRPCServerAddress = "0.0.0.0:12222"

	// http config
	UserCenterHttpServerAddress = "0.0.0.0:8080"

	PrivateIMRootCAPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/ca.pem"
	PrivateIMServerPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.pem"
	PrivateIMServerKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/server/server.key"
	PrivateIMClientPem = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.pem"
	PrivateIMClientKey = "/Users/wzy/GitPrograms/PrivateIM/CATLSFiles/client/client.key"

	// static file config, about operating the avatar and qrCode
	StaticFoldPath         = "/Users/wzy/GitPrograms/PrivateIM/UserService/static/"
	PhotoSaveFoldPath      = "/Users/wzy/GitPrograms/PrivateIM/UserService/static/photos/"
	PhotoSuffix            = ".png"
	PhotosUrlPrefix        = "/static/photos/" // if you use oss , should change this value
	DefaultAvatarPicName   = "defaultAvatar.jpg"
	AvatarPicUploadMaxSize = 100 * 2 << 10

	QRCodeBaseUrl = "http://127.0.0.1:8080/qrcontent/?"

	// password and auth token config
	PasswordHashSalt   = "this is a password hash salt"
	AuthTokenSalt      = "this is a auth token salt"
	AuthTokenAliveTime = 3600 * 24 //unit:second
	AuthTokenIssuer    = "userCenter"

	// reset password by authentication email config
	ResetPasswordTokenAliveTIme       = 60 * 5 // unit:second
	RestPasswordEmailSentTagAliveTime = 60 * 4 // unit:second
	RestPasswordPageBaseLink          = "http://127.0.0.1:8080/resetPassword.html?token="
	EmailServerHost                   = "smtp.163.com"
	EmailServerPort                   = 25
	EmailAuthUserName                 = "lhzqwlyyn@163.com"
	EmailAuthPassword                 = "wzy123456"
	RestPasswordEmailSubject          = "PrivateIM Reset Password"

	// redis db config
	RedisAddr = "127.0.0.1:6379"

	ElasticsearchServerAddress  = "0.0.0.0:9200"
)
