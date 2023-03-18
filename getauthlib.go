package main

import (
	"GoMCL2/weblib"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
)

func Getauthlib() (string, error) {
	rootUrl := "https://authlib-injector.yushi.moe/"
	if infos.NeedReplaceURL && infos.AuthlibInjectorRootURL != "" {
		rootUrl = infos.AuthlibInjectorRootURL
	}
	list, err := GetBytesFromNets(rootUrl + "artifact/latest.json")
	if err != nil {
		return "", err
	}
	listResult := gjson.ParseBytes(list)
	urls := listResult.Get("download_url").String()
	sha256s := listResult.Get("checksums.sha256").String()
	rb, err := weblib.GetBytesFromString(urls)
	if err != nil {
		return "", err
	}
	s := sha256.New()
	s.Write(rb)
	gots := hex.EncodeToString(s.Sum(nil))
	if gots != sha256s {
		return "", NewHashNotSame(sha256s, gots)
	}
	name := GetNameFromUrl(urls)
	name, err = filepath.Abs(name)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(name, rb, 0666)
	return name, err
}
