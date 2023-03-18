package main

import (
	"GoMCL2/weblib"
	"math/rand"

	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

const nilStr = ""

type NameNotFoundError struct {
	Name string
}

func (err NameNotFoundError) Error() string {
	return "name not found: " + err.Name
}
func NewNameNotFoundError(name string) *NameNotFoundError {
	return &NameNotFoundError{Name: name}
}

type NameListNotFoundError struct{}

func (err NameListNotFoundError) Error() string {
	return "name list not found"
}
func NewNameListNotFoundError() *NameListNotFoundError {
	return &NameListNotFoundError{}
}
func RandString(lenth int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, lenth)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func Login(userType *widget.SelectEntry, username, password, userEmail *widget.Entry, complianceLevel int) (string, string, string, string, error) {
	switch userType.Text {
	case "littleskin":
		tokens := RandString(len("ssssdddadfsdfsfsffsxxdsfewfsdf"))
		value := map[string]interface{}{
			"agent": map[string]interface{}{
				"name":    "Minecraft",
				"version": complianceLevel,
			},
			"username":    userEmail.Text,
			"password":    password.Text,
			"clientToken": tokens,
			"requestUser": false,
		}
		getstr, err := weblib.PostMapGotStrInStrWithHeader("https://littleskin.cn/api/yggdrasil/authserver/authenticate", value, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			return nilStr, nilStr, nilStr, nilStr, err
		}
		availableProfiles := gjson.Get(getstr, "availableProfiles")
		if !availableProfiles.Exists() || !availableProfiles.IsArray() {
			return nilStr, nilStr, nilStr, nilStr, NewNameListNotFoundError()
		}
		var uuid string = ""
		var hasName bool = false
		for _, v := range availableProfiles.Array() {
			if v.Get("name").String() == username.Text {
				hasName = true
				uuid = v.Get("id").String()
			}
		}
		if !hasName {
			return nilStr, nilStr, nilStr, nilStr, NewNameNotFoundError(username.Text)
		}
		accessToken := gjson.Get(getstr, "accessToken").String()
		return username.Text, uuid, accessToken, "mojang", nil
	}
	return nilStr, nilStr, nilStr, nilStr, nil
}
