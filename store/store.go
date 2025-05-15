package store

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/vouv/srun/model"
)

const accountFileName = "account.json"

var RootPath string

// SetAccount 设置账号信息
func SetAccount(username, password string) (err error) {
	return WriteAccount(&model.Account{
		Username: username,
		Password: password,
	})
}

// ReadAccount 读取账号信息
func ReadAccount() (account *model.Account, err error) {
	file, err := OpenAccountFile(os.O_RDONLY)
	if err != nil {
		log.Debugf("打开账号文件错误, %s,", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(base64.NewDecoder(base64.RawStdEncoding, file)).Decode(&account)
	return
}

// OpenAccountFile 打开账号文件
func OpenAccountFile(flag int) (file *os.File, err error) {
	accountFilename, err := getAccountFilename()
	if err != nil {
		return
	}
	return os.OpenFile(accountFilename, flag, 0600)
}

// WriteAccount 写入账号信息
func WriteAccount(account *model.Account) (err error) {
	file, err := OpenAccountFile(os.O_CREATE | os.O_TRUNC | os.O_WRONLY)
	if err != nil {
		log.Debugf("打开账号文件错误, %s", err)
		return
	}

	defer file.Close()

	enc := base64.NewEncoder(base64.RawStdEncoding, file)
	err = json.NewEncoder(enc).Encode(account)
	if err != nil {
		return err
	}
	return enc.Close()
}

// getAccountFilename 获取账号文件路径
func getAccountFilename() (fileSrc string, err error) {
	storageDir := filepath.Join(RootPath, ".srun")
	if _, sErr := os.Stat(storageDir); sErr != nil {
		if mErr := os.MkdirAll(storageDir, 0755); mErr != nil {
			log.Debugf("mkdir `%s` error, %s", storageDir, mErr)
			return
		}
	}
	fileSrc = filepath.Join(storageDir, accountFileName)
	return
}

func init() {
	curUser, gErr := user.Current()
	if gErr != nil {
		log.Fatalln("无法读取账户信息, 请重新设置账户信息")
	} else {
		RootPath = curUser.HomeDir
	}
}
