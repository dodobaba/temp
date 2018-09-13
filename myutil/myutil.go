package myutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"shoppingzone/conf"
	"shoppingzone/mylib/mylog"
	"strings"
	"time"
)

//ExistPath : check the path or file is exists
func ExistPath(path string) <-chan bool {
	out := make(chan bool, 3)
	go func() {
		_, err := os.Stat(path)
		if err == nil {
			out <- true
		}
		if os.IsNotExist(err) {
			out <- false
			mylog.Tf("[Error]", "MyUtil", "ExistPath", " %s", err.Error())
		}
		close(out)
	}()
	return out
}

//ListFiles : list all files in the path
func ListFiles(path string) <-chan string {
	out := make(chan string, 3)
	go func() {
		dirList, err := ioutil.ReadDir(path)
		if err != nil {
			mylog.Tf("[Error]", "MyUtil", "ListFiles", "List files in "+path+" was err. %s", err.Error())
		}
		for _, v := range dirList {
			out <- v.Name()
		}
		close(out)
	}()
	return out
}

//ReadFiles : read file in root path , out put file context
func ReadFiles(files <-chan string, rootpath string) <-chan ResoultReadFile {
	out := make(chan ResoultReadFile, 3)
	go func() {
		for file := range files {
			rs, err := ioutil.ReadFile(rootpath + "/" + file)
			if err != nil {
				mylog.Tf("[Error]", "MyUtil", "ReadFiles", "%s was read error. %s", file, err.Error())
			}
			out <- ResoultReadFile{FileName: file, Context: string(rs)}
		}
		close(out)
	}()
	return out
}

//String2JSON : transfer string to JSON
func String2JSON(str string) <-chan Message {
	out := make(chan Message)
	go func() {
		var mess Message
		err := json.Unmarshal([]byte(str), &mess)
		if err != nil {
			mylog.Tf("[Error]", "MyUtil", "String2JSON", "Fail to reansfer string to json. %s", err.Error())
		}
		out <- mess
		close(out)
	}()
	return out
}

//Int2Byte : transfer int to byte
func Int2Byte(i int) []byte {
	x := int32(i)
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, x)
	return buff.Bytes()
}

//Byte2Int : transfer byte to int
func Byte2Int(b []byte) int {
	buff := bytes.NewBuffer(b)
	var x int32
	binary.Read(buff, binary.BigEndian, &x)
	return int(x)
}

//GetRandomString : get random string
func GetRandomString(num int) <-chan string {
	output := make(chan string)
	go func() {
		chart := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		var t string
		for i := 0; i < num; i++ {
			j := GetRandomNumber(len(chart) - 1)
			t += chart[j]
		}
		output <- t
	}()
	return output
}

//GetRandomRangeNumber :
func GetRandomRangeNumber(digital int) string {
	var out int
	for i := 0; i < digital; i++ {
		n := GetRandomNumber(8) + 1
		out += n * int(math.Pow(10, float64(i)))
	}
	return fmt.Sprintf("%d", out)
}

//GetRandomNumber :
func GetRandomNumber(num int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(num)
}

//EnCrypto : crypto string
func EnCrypto(d string, s string) string {
	data := append([]byte(d), s...)
	has := sha512.Sum512_256(data)
	str := fmt.Sprintf("%x", has)
	return str
}

//TakeHashWord :
func TakeHashWord(s string) string {
	ts := time.Now().Format("2006-01-02 03:04:05.000")
	HashWord := EnCrypto(s+ts, conf.Config.PwdSrec)
	out := `{"UserKey":"` + s + `","ts":"` + ts + `","HashWord":"` + HashWord + `"}`
	return out
}

//EnCryptoToken :
//@s #string to encrypto
func EnCryptoToken(s string) <-chan string {
	out := make(chan string)
	go func() {
		key, _ := hex.DecodeString(conf.Config.CryptoKey)
		plaintext := []byte(s)
		block, err := aes.NewCipher(key)
		if err != nil {
			mylog.Tf("[Error]", "MyUtil", "EnCryptoToken", "Failed to Block key by Encrypto.  %s", err.Error())
		}
		ciphertext := make([]byte, aes.BlockSize+len(plaintext))
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(cr.Reader, iv); err != nil {
			mylog.Tf("[Error]", "MyUtil", "EnCryptoToken", "Failed to fill vi.  %s", err.Error())
		}
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
		out <- fmt.Sprintf("%x", ciphertext)
		close(out)
	}()
	return out
}

//DeCryptoToken :
func DeCryptoToken(s string) (<-chan string, error) {
	out := make(chan string)
	var e error
	go func(e error) {
		key, _ := hex.DecodeString(conf.Config.CryptoKey)
		ciphertext, _ := hex.DecodeString(s)
		block, err := aes.NewCipher(key)
		if err != nil {
			mylog.Tf("[Error]", "MyUtil", "DeCryptoToken", "Failed to Block key by Decrypto. %s", err.Error())
			e = err
			out <- s
			return
		}
		if len(ciphertext) < aes.BlockSize {
			mylog.Tf("[Error]", "MyUtil", "DeCryptoToken", "The ciphertext too short.")
			e = errors.New("The ciphertext too short")
			out <- s
			return
		}
		iv := ciphertext[:aes.BlockSize]
		ciphertext = ciphertext[aes.BlockSize:]
		stream := cipher.NewCFBDecrypter(block, iv)
		stream.XORKeyStream(ciphertext, ciphertext)
		mylog.Tf("[Info]", "MyUtil", "DeCryptoToken", "Success to decryto token. %s | %s", s, fmt.Sprintf("%s", ciphertext))
		e = nil
		out <- fmt.Sprintf("%s", ciphertext)
		close(out)
	}(e)
	return out, e
}

//InArray :
func InArray(array []string, value ...interface{}) bool {
	if len(value) == 1 && reflect.TypeOf(value[0]).String() == "string" {
		for _, t := range array {
			if value[0] == t {
				return true
			}
		}
		return false
	}
	for _, t := range array {
		for _, y := range value {
			if t == y {
				return true
			}
		}
	}
	return false
}

/*
CompareArrary :
@Param {arrary} a
@Param {arrary} b
@Retrun {arrary} va
@Return {arrary} vb
@Return {arrary} uab
*/
func CompareArrary(a []string, b []string) ([]string, []string, []string) {
	var va, vb, uab []string
	va = a
	vb = b
	if len(a) > 0 && len(b) > 0 {
		for ia, ta := range a {
			for ib, tb := range b {
				if ta == tb {
					uab = append(uab, ta)
					va = DropArrayElement(va, ia)
					vb = DropArrayElement(vb, ib)
				}
			}
		}
	}
	return va, vb, uab
}

//DropArrayElement :
func DropArrayElement(a []string, idx int) []string {
	return append(a[:idx], a[idx+1:]...)
}

//CheckOSArgs :
func CheckOSArgs() OutOSArgs {
	var out OutOSArgs
	var tp []string
	if len(os.Args) > 1 {
		for idx, v := range os.Args {
			if v == "-log" {
				out.LogPath = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-log=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-log=" {
				out.LogPath = strings.Split(v, "=")[1]
				continue
			}
			if v == "-P" {
				out.HTTPPort = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-P=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-P=" {
				out.HTTPPort = strings.Split(v, "=")[1]
				continue
			}
			if v == "-p" {
				out.NetProt = os.Args[idx+1]
				continue
			}
			tp = regexp.MustCompile(`^-p=`).FindAllString(v, -1)
			if len(tp) > 0 && tp[0] == "-p=" {
				out.NetProt = strings.Split(v, "=")[1]
				continue
			}
		}
	}
	if out.LogPath == "" {
		out.LogPath = conf.Config.LogFile
	}
	if out.HTTPPort == "" {
		out.HTTPPort = conf.Config.HTTPPort
	}
	if out.NetProt == "" {
		out.NetProt = conf.Config.NetPort
	}
	return out
}

//TARfile :
func TARfile(p string, f string) {
	tarCmd := exec.Command("tar", "-zcf", p+"/tar/"+f+".tar.gz", p+"/"+f, "--remove-files")
	tarCmd.Run()
	rmCmd := exec.Command("rm", "-fv", p+"/"+f)
	rmCmd.Run()
}
