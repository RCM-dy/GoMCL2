package main

import (
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

type NotSameType struct {
	need gjson.Type
	got  gjson.Type
}

func NewNotSameType(need, got gjson.Type) *NotSameType {
	return &NotSameType{need: need, got: got}
}
func (n *NotSameType) Error() string {
	return "not same type,\nneed: " + n.need.String() + "got: " + n.got.String()
}
func GetObj(outLabel *widget.Label, index gjson.Result) error {
	if !index.Exists() {
		return &NotFound{}
	}
	needBackup := false
	if index.Get("map_to_resources").Exists() {
		if index.Get("map_to_resources").IsBool() {
			if index.Get("map_to_resources").Bool() {
				needBackup = true
			}
		}
	}
	objs := index.Get("objects")
	if !objs.IsObject() {
		return NewNotSameType(gjson.JSON, objs.Type)
	}
	err := os.RemoveAll(infos.Objdir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(infos.ObjBackdir, 0666)
	if err != nil {
		return err
	}
	var theErr error = nil
	objs.ForEach(func(key, value gjson.Result) bool {
		var rootUrl string = "https://resources.download.minecraft.net/"
		if needReplaceAssetsRootURL() {
			rootUrl = infos.AssetsRootURL
		}
		hashcode := value.Get("hash").String()
		twoHash := hashcode[:2]
		url := rootUrl + twoHash + "/" + hashcode
		objpath := filepath.Join(infos.Objdir, twoHash, hashcode)
		objbackuppath := filepath.Join(infos.ObjBackdir, ReplaceByMap(key.String(), map[string]string{
			"/": "\\",
		}))
		err := os.MkdirAll(filepath.Dir(objpath), 0666)
		if err != nil {
			theErr = err
			return false
		}
		rb, err := GetBytesFromNets(url)
		if err != nil {
			theErr = err
			return false
		}
		err = WriteBytes(objpath, rb)
		if err != nil {
			theErr = err
			return false
		}
		if needBackup {
			err := os.MkdirAll(filepath.Dir(objbackuppath), 0666)
			if err != nil {
				theErr = err
				return false
			}
			err = WriteBytes(objbackuppath, rb)
			if err != nil {
				theErr = err
				return false
			}
		}
		outLabel.SetText(hashcode)
		time.Sleep(800)
		return true
	})
	return theErr
}
