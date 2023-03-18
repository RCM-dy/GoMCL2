package main

import "github.com/tidwall/gjson"

type NotFound struct{}

func (n *NotFound) Error() string {
	return "not found"
}

func GetIndexBytes(indexes gjson.Result) ([]byte, error) {
	if !indexes.Exists() {
		return nil, &NotFound{}
	}
	sha1 := indexes.Get("sha1").String()
	url := indexes.Get("url").String()
	if needReplaceIndexRootURL() {
		url = ReplaceByMap(url, map[string]string{
			"https://piston-meta.mojang.com":  infos.IndexRootURL,
			"https://launchermeta.mojang.com": infos.IndexRootURL,
			"https://launcher.mojang.com":     infos.IndexRootURL,
		})
	}
	rb, err := GetBytesFromNets(url)
	if err != nil {
		return nil, err
	}
	rsha1 := Sha1Bytes(rb)
	if rsha1 != sha1 {
		return nil, NewHashNotSame(sha1, rsha1)
	}
	return rb, nil
}
