package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	qrcode "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/goSeeFuture/gblog/configs"
	"github.com/gofiber/fiber/v2"
)

const authorField = "author"

func randDraftPreviewValue() (string, error) {
	var key = make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}

func getAuthorPreviewValue() (string, error) {
	_, err := os.Stat(authorField)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	var value string
	if os.IsNotExist(err) {
		value, err = randDraftPreviewValue()
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(authorField, []byte(value), os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	b, err := ioutil.ReadFile(authorField)
	if err != nil {
		return "", err
	}

	value = string(b)
	return value, nil
}

func draftQRCode() (*qrcode.QRCodeString, string) {
	draftPreviewValue, err := getAuthorPreviewValue()
	if err != nil {
		return nil, ""
	}

	url := configs.Setting.PublicDomain + "/" + authorField + "/" + draftPreviewValue
	return qrcode.New().Get(url), url
}

func isAuthor(c *fiber.Ctx) bool {
	sess, err := store.Get(c)
	if err != nil {
		panic(err) // middleware catch panic
	}

	v := sess.Get(authorField)
	fmt.Println(authorField, v)
	if v == nil {
		return false
	}

	return v.(bool)
}

func PrintAuthorSecret() {
	qrcode, url := draftQRCode()
	fmt.Println("预览草稿二维码：")
	qrcode.Print()
	fmt.Println("预览草稿链接：", url)
}
