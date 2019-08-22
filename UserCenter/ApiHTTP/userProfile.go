package ApiHTTP

import (
	"../DataLayer"
	"../utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"strconv"
)

const (
	AuthTokenSalt      = "this is a auth token salt"
	AuthTokenAliveTime = 3600 * 24 //unit:second
	AuthTokenIssuer    = "userCenter"
	JWTGetUserId         = "user_id"

	PhotoSaveFoldPath   = "/Users/wzy/GitProrgram/PrivateIM/UserCenter/static/photos/"
	PhotoSuffix         = ".png"
	DefaultAvatarUrl    = "/static/photos/defaultAvatar.jpg"
	PhotosUrlPrefix     = "/static/photos/" // if you use oss , should change this value
	MaxAvatarUploadSize = 100 * 2 << 10

	QRCodeBaseUrl = "http://127.0.0.1:8080/qrcode/"
)

// GetProfile HTTP API function
func GetProfile(c *gin.Context) {
	user := DataLayer.UserBasic{Id: c.MustGet(JWTGetUserId).(int64)}
	err := DataLayer.MySQLGetUserById(&user)
	if nil != err {
		c.JSON(404, gin.H{"error": "get user information fail"})
		return
	}
	c.JSON(200, user)
}

type TempProfile struct {
	Name   string `json:"name" binding:"nameValidator"`
	Mobile string `json:"mobile" binding:"mobileValidator"`
	Gender int    `json:"gender" binding:"genderValidator"`
}

// PutProfile HTTP API function
func PutProfile(c *gin.Context) {
	// Validate the params
	tempProfileP, err := parseTempProfile(c)
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	errs := binding.Validator.ValidateStruct(tempProfileP)
	if nil != errs {
		c.JSON(400, gin.H{"errors": errs.Error()})
		return
	}
	// Update user info
	userId := c.MustGet(JWTGetUserId)
	err = DataLayer.MySQLUpdateProfile(tempProfileP.Name, tempProfileP.Mobile, tempProfileP.Gender, userId.(int64))
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tempProfileP)
}

// Parse the JsonBodyParams to TempProfile
func parseTempProfile(c *gin.Context) (*TempProfile, error) {
	// Parse the JsonBodyParams to map
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	if n == 0 {
		return nil, errors.New("not have any JsonBodyParams")
	}

	tempDict := make(map[string]interface{})
	_ = json.Unmarshal(buf[0:n], &tempDict)
	// Check the integrity of parameters
	name, ok := tempDict["name"]
	if !ok {
		return nil, errors.New("`name` not exited in JsonBodyParams")
	}
	mobile, ok := tempDict["mobile"]
	if !ok {
		return nil, errors.New("`mobile` not exited in JsonBodyParams")
	}
	gender, ok := tempDict["gender"]
	if !ok {
		return nil, errors.New("`gender` not exited in JsonBodyParams")
	}

	tempProfileP := &TempProfile{
		Name:   name.(string),
		Mobile: mobile.(string),
		Gender: int(gender.(float64))}
	return tempProfileP, nil
}

// GetAvatar HTTP API function
func GetAvatar(c *gin.Context) {
	userId := c.MustGet(JWTGetUserId)
	avatar := new(string)
	err := DataLayer.MySQLGetUserAvatar(userId.(int64), avatar)
	if nil != err {
		c.JSON(500, gin.H{"error": "query avatar fail"})
		return
	}
	if *avatar == "" {
		c.JSON(200, gin.H{"avatar": DefaultAvatarUrl})
		return
	}
	c.JSON(200, gin.H{"avatar_url": PhotosUrlPrefix + *avatar + PhotoSuffix})
}

// PutAvatar HTTP API function
func PutAvatar(c *gin.Context) {
	// get file data and hash value as name
	file, err := c.FormFile("new_avatar")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size > MaxAvatarUploadSize || file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=100kb"})
		return
	}
	hashName, data, err := utils.GinFormFileHash(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// check if the hash name is existed more then one in the table
	// if true, it meanings that has the same file already upload.
	// not need to save the file again.
	count := DataLayer.MySQLAvatarHashNameCount(hashName)
	if count == 0 {
		// save the file data to local or static server
		if err := UploadAvatarLocal(data, hashName); nil != err {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// save the information into database
	userId := c.MustGet(JWTGetUserId)
	err = DataLayer.MySQLPutUserAvatar(userId.(int64), hashName)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"avatar_url": PhotosUrlPrefix + hashName + PhotoSuffix})
}

// GetQRCode HTT API function
func GetQrCode(c *gin.Context) {
	// try to get QrCode hash name from database. if existed, return.
	userId := c.MustGet(JWTGetUserId)
	hashNameP := new(string)
	err := DataLayer.MySQLGetUserQRCode(userId.(int64), hashNameP)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if *hashNameP != "" {
		c.JSON(200, gin.H{"qr_code": PhotosUrlPrefix + *hashNameP + PhotoSuffix})
		return
	}

	// if the qr code hash name is not existed, create an new and save
	content := QRCodeContent(userId.(int64))
	data, _ := utils.CreatQRCodeBytes(content)
	*hashNameP = utils.BytesDataHash(data)
	err = SaveQRCodeLocal(data, *hashNameP)
	if nil != err {
		c.JSON(500, gin.H{"error": "create QRCode fail"})
		return
	}
	err = DataLayer.MySQLPutUserQRCode(userId.(int64), *hashNameP)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"qr_code_url": PhotosUrlPrefix + *hashNameP + PhotoSuffix})

}

// ParseQrCode HTTP API function
func ParseQrCode(c *gin.Context) {
	// get file data
	file, err := c.FormFile("qr_code")
	if nil != err {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if file.Size == 0 {
		c.JSON(400, gin.H{"error": "the upload image size need gt=0kb and lte=2MB"})
		return
	}
	_, data, err := utils.GinFormFileHash(file)
	if nil != err {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	content, err := utils.ParseQRCodeBytes(data)
	if nil != err {
		c.JSON(400, gin.H{"error": "QRCode parse fail"})
		return
	}
	c.JSON(200, gin.H{"qr_content": content})

}

// save avatar file to local
func UploadAvatarLocal(data []byte, hashName string) error {
	prefix := PhotoSaveFoldPath
	suffix := PhotoSuffix
	path := prefix + hashName + suffix
	if err := utils.UploadFileToLocal(data, path); nil != err {
		return err
	}
	return nil
}

// todo make the content for create a QRCode
func QRCodeContent(userId int64) string {
	//QRCodeBaseUrl
	return "https://www.baidu.com " + strconv.FormatInt(userId,10) // temp value
}

// save QRCode file to local
func SaveQRCodeLocal(data []byte, hashName string) error {
	savePath := PhotoSaveFoldPath + hashName + PhotoSuffix
	err := ioutil.WriteFile(savePath, data, 0644)
	if nil != err {
		return err
	}
	return nil
}

// todo upload the data to cloud with a hashName
func UploadDataToCloud(data []byte, hashName string) error {
	return nil
}