package service

import (
	"bytes"
	"image"

	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"shoppingzone/conf"
	constants "shoppingzone/mylib/myconst"
	"shoppingzone/mylib/mylog"
	"shoppingzone/mylib/mymgo"
	"shoppingzone/mylib/myredis"
	"shoppingzone/myutil"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
CheckContentType :
@ 'c' #the http context
@ 'obj' #struct for http body
*/
func CheckContentType(c *gin.Context, obj interface{}) (string, string, error) {
	ip := c.ClientIP()
	contentType := c.Request.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		if err := c.ShouldBindJSON(obj); err != nil {
			return contentType, ip, err
		}
		break
	default:
		if err := c.ShouldBind(obj); err != nil {
			return contentType, ip, err
		}
		break
	}
	return contentType, ip, nil
}

/*
Cache :
@ 'k' #cache key
@ 'v' #cache value
@ 'exp' #expire time secends
*/
func Cache(k string, v string, exp int) <-chan bool {
	out := make(chan bool)
	go func() {
		err := myredis.Set(k, v, exp)
		if err != nil {
			mylog.Tf("[Error]", "HTTP Service", "Failed to save cache . %s", err.Error())
			out <- false
		}
		out <- true
		close(out)
	}()
	return out
}

/*
Authorization ：
@any roles of Group,Role,Label,Ruouter, can muilt use "," to slipt
*/
func Authorization(arg ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var merchantKey, userKey []byte
		ip := c.ClientIP()
		mylog.SetIP(ip)
		if c.GetHeader("Authorization") == "" && c.Query("token") == "" {
			c.Abort()
			mylog.Tf("[Error]", "App", "Authorization", "Fail to authorization,must auth.")
			c.JSON(200, gin.H{"Message": "Fail to authorization,must auth", "Status": "Fail"})
			return
		}
		h := c.GetHeader("Authorization")
		if h == "" && c.Query("token") != "" && len(c.Query("token")) >= 128 {
			th := []byte(c.Query("token"))
			merchantKey = th[:64]
			userKey = th[64:]
			h = string(userKey)
		}
		userToken, e := myutil.DeCryptoToken(h)
		if e != nil {
			c.Abort()
			mylog.Tf("[Error]", "App", "Authorization", "Fail to authorization,must auth. %s", e.Error())
			c.JSON(200, gin.H{"Message": "Fail to authorization,must auth", "Status": "Fail", "err": e.Error()})
			return
		}
		UserKey := <-userToken
		type cuser struct {
			ID              bson.ObjectId `bson:"_id" form:"_id" json:"_id"`
			Active          bool          `bson:"Active" form:"Active" json:"Active"`
			UserKey         string        `bson:"UserKey" form:"UserKey" json:"UserKey"`
			SecureLabel     []string      `bson:"SecureLabel" form:"SecureLabel" json:"SecureLabel"`
			SecureGroup     []string      `bson:"SecureGroup" form:"SecureGroup" json:"SecureGroup"`
			SecureRouter    []string      `bson:"SecureRouter" form:"SecureRouter" json:"SecureRouter"`
			SecureRole      []string      `bson:"SecureRole" form:"SecureRole" json:"SecureRole"`
			AdminMerchant   []string      `bson:"AdminMerchant" form:"AdminMerchant" json:"AdminMerchant"`
			ManagerMerchant []string      `bson:"ManagerMerchant" form:"ManagerMerchant" json:"ManagerMerchant"`
			UserMerchant    []string      `bson:"UserMerchant" form:"UserMerchant" json:"UserMerchant"`
		}
		var u cuser
		query := func(c *mgo.Collection) error {
			b := bson.M{"UserKey": UserKey, "Active": true}
			return c.Find(b).One(&u)
		}
		err := mymgo.Do("User", query)
		if err != nil {
			c.Abort()
			mylog.Tf("[Error]", "App", "Authorization", "Fail to find this user or this user havn't authorization. %s | %s", err.Error(), UserKey)
			c.JSON(200, gin.H{"err": "Fail to find this user or this user havn't authorization", "error": err.Error()})
			return
		} else {
			for _, uav := range arg {
				switch uav {
				case constants.MERCHANTADMIN:
					if !myutil.InArray(u.AdminMerchant, string(merchantKey)) {
						c.Abort()
						mylog.Tf("[Error]", "App", "Authorization", "Fail to this user access to this merchant. %s ", merchantKey)
						c.JSON(200, gin.H{"Status": "Fail", "Message": "Fail to this user access to this merchant."})
						return
					} else {
						mylog.Tf("[Info]", "App", "Authorization", "Success to this user access to this merchant with admin. %s", userKey)
					}

					break
				case constants.MERCHANTMANAGER:
					if !myutil.InArray(u.ManagerMerchant, string(merchantKey)) {
						c.Abort()
						mylog.Tf("[Error]", "App", "Authorization", "Fail to this user access to this merchant. %s ", merchantKey)
						c.JSON(200, gin.H{"Status": "Fail", "Message": "Fail to this user access to this merchant."})
						return
					} else {
						mylog.Tf("[Info]", "App", "Authorization", "Success to this user access to this merchant with manager. %s", userKey)
					}

					break
				case constants.MERCHANTUSER:
					if !myutil.InArray(u.UserMerchant, string(merchantKey)) {
						c.Abort()
						mylog.Tf("[Error]", "App", "Authorization", "Fail to this user access to this merchant. %s ", merchantKey)
						c.JSON(200, gin.H{"Status": "Fail", "Message": "Fail to this user access to this merchant."})
						return
					} else {
						mylog.Tf("[Info]", "App", "Authorization", "Success to this user access to this merchant with user. %s", userKey)
					}

					break
				case constants.ADMIN:
					if !myutil.InArray(u.SecureGroup, constants.ADMIN) && !myutil.InArray(u.SecureLabel, constants.ADMIN) && !myutil.InArray(u.SecureRole, constants.ADMIN) && !myutil.InArray(u.SecureRouter, constants.ADMIN) {
						c.Abort()
						mylog.Tf("[Error]", "App", "Authorization", "Fail to this user access to function. %s", userKey)
						c.JSON(200, gin.H{"Status": "Fail", "Message": "Fail to this user access to this function."})
						return
					} else {
						mylog.Tf("[Info]", "App", "Authorization", "Success to this user access to this function. %s", userKey)
					}
					break
				default:
					c.Abort()
					mylog.Tf("[Error]", "App", "Authorization", "Must have currect auth.")
					c.JSON(200, gin.H{"Status": "Fail", "Message": "Must have currect auth."})
					return
				}
			}
		}
		c.Next()
	}
}

type loadFileProcess struct {
	hashname   string
	ip         string
	contentype string
}

func (l *loadFileProcess) querydb(loadfile UpLoadFileRs, contenttyps []string) <-chan bool {
	out := make(chan bool, 10)
	go func() {
		b := bson.M{"HashName": l.hashname, "ContentType": bson.M{"$in": contenttyps}}
		query := func(c *mgo.Collection) error {
			return c.Find(b).One(&loadfile)
		}
		err := mymgo.Do("Files", query)
		if err != nil {
			mylog.Tf("[Error]", "APP", "LoadImage", "Fail to query this file in db. %s | %s", l.hashname, err.Error())
			out <- false
		} else {
			l.contentype = loadfile.ContentType
			out <- true
		}
		close(out)
	}()
	return out
}

func (l *loadFileProcess) openfile(s <-chan bool) <-chan *os.File {
	out := make(chan *os.File, 3)
	go func() {
		if <-s {
			imgfile, err := os.Open(conf.Uploadconfig.Uploadpath + `/` + l.hashname)
			if err != nil {
				mylog.Tf("[Error]", "APP", "LoadImage", "Fail to read file. %s | %s", l.hashname, err.Error())
			}
			out <- imgfile
		}
		close(out)
	}()
	return out
}

func (l *loadFileProcess) processImage(f <-chan *os.File, width string, height string) []byte {
	out := make(chan []byte, 10)
	go func() {
		var img, outimg image.Image
		var err error
		rs := <-f
		if rs != nil {
			imgbuff := new(bytes.Buffer)
			nw, _ := strconv.Atoi(width)
			nh, _ := strconv.Atoi(height)
			switch l.contentype {
			case `image/jpeg`:
				img, err = jpeg.Decode(rs)
				break
			case `image/png`:
				img, err = png.Decode(rs)
				break
			case `image/gif`:
				img, err = gif.Decode(rs)
				break
			default:
				img, err = jpeg.Decode(rs)
				break
			}
			defer rs.Close()
			if err != nil {
				mylog.Tf("[Error]", "APP", "LoadImage", "Fail to decode file. %s | %s", l.hashname, err.Error())
			}
			bonds := img.Bounds()
			dx := bonds.Dx()
			dy := bonds.Dy()
			if height == "" && width == "" {
				outimg = img
				mylog.Tf("[Info]", "APP", "LoadImage", "Success to load image original and sent it. %s", l.hashname)
			} else if height == "" && width != "" {
				outimg = resize.Resize(uint(nw), uint(nw*dy/dx), img, resize.Lanczos3)
				mylog.Tf("[Info]", "APP", "LoadImage", "Success to load image with width %d and sent it. %s", nw, l.hashname)
			} else if height != "" && width == "" {
				outimg = resize.Resize(uint(nh*dx/dy), uint(nh), img, resize.Lanczos3)
				mylog.Tf("[Info]", "APP", "LoadImage", "Success to load image with height %d and sent it. %s", nh, l.hashname)
			} else if height != "" && width != "" {
				outimg = resize.Resize(uint(nw), uint(nh), img, resize.Lanczos3)
				mylog.Tf("[Info]", "APP", "LoadImage", "Success to load image with height %d and width %d and sent it. %s", nw, nh, l.hashname)
			}
			switch l.contentype {
			case `image/jpeg`:
				err = jpeg.Encode(imgbuff, outimg, nil)
				break
			case `image/png`:
				err = png.Encode(imgbuff, outimg)
				break
			case `image/gif`:
				err = gif.Encode(imgbuff, outimg, nil)
				break
			default:
				err = jpeg.Encode(imgbuff, outimg, nil)
				break
			}
			if err != nil {
				mylog.Tf("[Error]", "APP", "LoadImage", "Fail to encode file. %s | %s", l.hashname, err.Error())
			}
			out <- imgbuff.Bytes()
		} else {
			out <- []byte(`Can't load this image!`)
		}
		close(out)
	}()
	return <-out
}

/*
LoadImage :
*/
func LoadImage(c *gin.Context) {

	var loadfile UpLoadFileRs
	contenttyps := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/bmp",
	}
	hashname := c.Param("hashkey")
	width := c.Query("w")
	height := c.Query("h")
	ip := c.ClientIP()
	mylog.SetIP(ip)
	peipline := loadFileProcess{
		hashname: hashname,
		ip:       ip,
	}
	outimage := peipline.processImage(peipline.openfile(peipline.querydb(loadfile, contenttyps)), width, height)
	c.Data(200, loadfile.ContentType, outimage)
}

/*
Uploadfiles :
@files #multipat form feild name must be 'files' to set file type
*/
func Uploadfiles(c *gin.Context) {
	ip := c.ClientIP()
	mylog.SetIP(ip)
	var hashArrary, fileArrary, failArray []string
	isExists := <-myutil.ExistPath(conf.Uploadconfig.Uploadpath)
	if !isExists {
		if err := os.Mkdir(conf.Uploadconfig.Uploadpath, 0777); err != nil {
			mylog.Tf("[Error]", "App", "Uploadfiles", "can't make path for upload files %s", err.Error())
			c.JSON(200, gin.H{"Status": "Fail", "Message": "error upload path", "Error": err.Error()})
			return
		}
	}
	form, err := c.MultipartForm()
	if err != nil {
		mylog.Tf("[Error]", "App", "Uploadfiles", " %s", err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	files := form.File["files"]
	checkRs, failRs := uploadfileCheckLimit(files)
	saveFiles, saveError := uploadfilesProcess(checkRs, c)
	inDBFiles := saveFilesToDB(saveFiles)
	for f := range inDBFiles {
		mylog.Tf("[Info]", "App", "Uploadfiles", "Success to upload file %s | %s | %s | %d", f.OriginalName, f.HashName, f.ContentType, f.Size)
		fileArrary = append(fileArrary, f.OriginalName)
		hashArrary = append(hashArrary, f.HashName)
	}
	for frs := range failRs {
		failArray = append(failArray, frs)
	}
	for serr := range saveError {
		failArray = append(failArray, serr)
	}
	c.JSON(200, gin.H{"Status": "Sucess", "Message": "uploaded!.", "Files": fileArrary, "Hash": hashArrary, "FailFiles": failArray})
}

func uploadfileCheckLimit(files []*multipart.FileHeader) (<-chan *multipart.FileHeader, <-chan string) {
	out := make(chan *multipart.FileHeader, conf.Uploadconfig.Maxfiles)
	out2 := make(chan string, conf.Uploadconfig.Maxfiles)
	go func() {
		for _, file := range files {
			tf := strings.Split(file.Filename, `.`)
			xf := tf[len(tf)-1]
			if !myutil.InArray(conf.Uploadconfig.Allowtype, xf) {
				mylog.Tf("[Error]", "App", "Uploadfiles", "can't upload this file %s ", file.Filename)
				out2 <- `Not support this file: '` + file.Filename + `'`
				continue
			}
			if !myutil.InArray(conf.Uploadconfig.Allowheaders, file.Header.Get("Content-Type")) {
				mylog.Tf("[Error]", "App", "Uploadfiles", "can't upload this meta data file %s type %s", file.Filename, file.Header.Get("Content-Type"))
				out2 <- `Not support this file: '` + file.Filename + `'`
				continue
			}
			if file.Size > conf.Uploadconfig.Maxsize {
				mylog.Tf("[Error]", "App", "Uploadfiles", "can't upload this file %s too large size %d", file.Filename, file.Size)
				out2 <- `File's size too large: '` + file.Filename + `'`
				continue
			}
			out <- file
		}
		close(out)
		close(out2)
	}()
	return out, out2
}

func uploadfilesProcess(infile <-chan *multipart.FileHeader, c *gin.Context) (<-chan UpLoadFileRs, <-chan string) {
	out := make(chan UpLoadFileRs, conf.Uploadconfig.Maxfiles)
	out2 := make(chan string, conf.Uploadconfig.Maxfiles)
	go func() {
		for file := range infile {
			tf := strings.Split(file.Filename, `.`)
			f := myutil.EnCrypto(tf[0], time.Now().Format(`20060102`))
			if err := c.SaveUploadedFile(file, conf.Uploadconfig.Uploadpath+`/`+f); err != nil {
				mylog.Tf("[Error]", "App", "Uploadfiles", "upload this file %s failed. %s", file.Filename, err.Error())
				out2 <- `'` + file.Filename + `' upload error:` + err.Error()
			} else {
				out <- UpLoadFileRs{
					HashName:     f,
					OriginalName: file.Filename,
					ContentType:  file.Header.Get("Content-Type"),
					Size:         file.Size,
				}
			}
		}
		close(out)
		close(out2)
	}()
	return out, out2
}

func saveFilesToDB(savefile <-chan UpLoadFileRs) <-chan UpLoadFileRs {
	out := make(chan UpLoadFileRs, conf.Uploadconfig.Maxfiles)
	go func() {
		for f := range savefile {
			b := bson.M{"HashName": f.HashName, "OriginalName": f.OriginalName, "ContentType": f.ContentType, "Size": f.Size, "Active": true, "CreateTime": time.Now()}
			query := func(c *mgo.Collection) error {
				return c.Insert(b)
			}
			err := mymgo.Do("Files", query)
			if err != nil {
				mylog.Tf("[Error]", "APP", "Uploadfiles", "Fail to save file into db. %s | %s", f.OriginalName, err.Error())
			} else {
				mylog.Tf("[Info]", "APP", "Uploadfiles", "Success to save file into db. %s ", f.OriginalName)
				out <- f
			}
		}
		close(out)
	}()
	return out
}
