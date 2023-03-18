package main

import (
	"GoMCL2/weblib"
	"GoMCL2/zipfile"
	"errors"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

func GetLib(verJson gjson.Result, outLabel *widget.Label) (cp string, err error) {
	if !verJson.Exists() {
		err = &NotFound{}
		return
	}
	clientDownload := verJson.Get("downloads")
	if !clientDownload.Exists() {
		err = &NotFound{}
		return
	}
	clients := clientDownload.Get("client")
	if !clients.Exists() {
		err = &NotFound{}
		return
	}
	clientUrls := clients.Get("url")
	if !clientUrls.Exists() {
		err = &NotFound{}
		return
	}
	clientUrl := clientUrls.String()
	if needReplaceClienJarRootURL() {
		clientUrl = ReplaceByMap(clientUrl, map[string]string{
			"https://piston-data.mojang.com":  infos.ClientJarRootURL,
			"https://piston-meta.mojang.com":  infos.ClientJarRootURL,
			"https://launchermeta.mojang.com": infos.ClientJarRootURL,
			"https://launcher.mojang.com":     infos.ClientJarRootURL,
		})
	}
	var clientB []byte
	clientB, err = weblib.GetBytesFromString(clientUrl)
	if err != nil {
		return
	}
	needVer := verJson.Get("id").String()
	clientPath := filepath.Join(infos.Verdir, needVer, needVer+".jar")
	err = WriteBytes(clientPath, clientB)
	if err != nil {
		return
	}
	nativeList := []int{}
	libarray := verJson.Get("libraries").Array()
	for k, v := range libarray {
		download := v.Get("downloads")
		if !download.Exists() {
			err = &NotFound{}
			return
		}
		if download.Get("classifiers").Exists() {
			nativeList = append(nativeList, k)
		}
		rules := v.Get("rules")
		if rules.Exists() {
			var notsame bool = false
			for _, vr := range rules.Array() {
				action := vr.Get("action")
				if action.Exists() {
					if action.String() != "allow" {
						continue
					}
				}
				oses := vr.Get("os")
				if oses.Exists() {
					oses.ForEach(func(key, value gjson.Result) bool {
						if key.String() == "name" && value.String() != "windows" {
							return false
						}
						if key.String() == "arch" && value.String() != infos.Arch {
							notsame = true
							return false
						}
						return true
					})
				}
			}
			if notsame {
				continue
			}
		}
		outLabel.Refresh()
		artifact := download.Get("artifact")
		if !artifact.Exists() {
			if !download.Get("classifiers").Exists() {
				err = &NotFound{}
				return
			}
		}
		path := artifact.Get("path").String()
		path = ReplaceByMap(path, map[string]string{
			"/": "\\",
		})
		path = filepath.Join(infos.Libdir, path)
		err = os.MkdirAll(filepath.Dir(path), 0666)
		if err != nil {
			return
		}
		urls := artifact.Get("url").String()
		if needReplaceLibRootURL() {
			urls = ReplaceByMap(urls, map[string]string{
				"https://libraries.minecraft.net": infos.LibRootURL,
			})
		}
		var rb []byte
		rb, err = weblib.GetBytesFromString(urls)
		if err != nil {
			return
		}
		err = WriteBytes(path, rb)
		if err != nil {
			return
		}
		outLabel.SetText(path)
		cp += path + ";"
		time.Sleep(800)
	}
	cp += clientPath
	nativePath = filepath.Join(infos.Verdir, needVer, "native")
	err = os.MkdirAll(nativePath, 0666)
	if len(nativeList) == 0 {
		return
	}
	err = os.RemoveAll(".\\nativestmp")
	if err != nil {
		return
	}
	err = os.MkdirAll(".\\nativestmp", 0666)
	if err != nil {
		return
	}
	for _, v := range nativeList {
		native := libarray[v]
		downloads := native.Get("downloads")
		classifiers := downloads.Get("classifiers")
		natives := native.Get("natives")
		winkey := natives.Get("windows")
		if !winkey.Exists() {
			continue
		}
		windownloads := classifiers.Get(winkey.String())
		if !windownloads.Exists() {
			continue
		}
		var rb weblib.Bytes
		if windownloads.Get("url").String() == "" {
			err = errors.New("empty URL:" + windownloads.Get("url").String())
			return
		}
		url := windownloads.Get("url").String()
		if needReplaceLibRootURL() {
			url = ReplaceByMap(url, map[string]string{
				"https://libraries.minecraft.net": infos.LibRootURL,
			})
		}
		rb, err = weblib.GetBytesFromString(url)
		if err != nil {
			return
		}
		err = WriteBytes(filepath.Join("nativestmp", "tmp.jar"), rb)
		if err != nil {
			return
		}
		err = zipfile.Unzip(filepath.Join("nativestmp", "tmp.jar"), nativePath)
		if err != nil {
			return
		}
		var allnatives []string
		allnatives, err = GetAllFilePath(nativePath, allnatives)
		if err != nil {
			return
		}
		for _, v := range allnatives {
			if filepath.Ext(v) != ".dll" {
				err = os.Remove(v)
				if err != nil {
					return
				}
				continue
			}
			var dllbit string
			dllbit, err = GetDllBit(v)
			if err != nil {
				return
			}
			if dllbit != infos.Arch {
				err = os.Remove(v)
				if err != nil {
					return
				}
				continue
			}
		}
		time.Sleep(800)
	}
	err = os.RemoveAll(filepath.Join(nativePath, "META-INF"))
	return
}
