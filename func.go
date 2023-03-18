package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
)

func WriteBytes(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func GetNameFromUrl(url string) string {
	urlArray := strings.Split(url, "/")
	return urlArray[len(urlArray)-1]
}
func ReplaceByMap(s string, c map[string]string) string {
	ss := s
	for k, v := range c {
		ss = strings.ReplaceAll(ss, k, v)
	}
	return ss
}
func GetBytesFromNets(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
func Sha1Bytes(datas []byte) string {
	s := sha1.New()
	s.Write(datas)
	return hex.EncodeToString(s.Sum(nil))
}
func MakeDirs(dirs ...string) error {
	for _, v := range dirs {
		err := os.MkdirAll(v, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
func runCMD(commands string) (output string, err error) {
	cmdpath := filepath.Join(os.Getenv("windir"), "System32", "cmd.exe")
	cmd := exec.Command(cmdpath, "/c", commands)
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Run()
	var o []byte
	o, err = io.ReadAll(out)
	if err != nil {
		return
	}
	output = string(o)
	return
}
func GetJavas() (m map[string]string, err error) {
	m = make(map[string]string)
	var c string
	c, err = runCMD("where java")
	if err != nil {
		return
	}
	var o []byte
	for _, v := range strings.Split(c, "\n") {
		if v == "" {
			continue
		}
		v = strings.ReplaceAll(v, "\r", "")
		cmd := exec.Command(v, "-version")
		out := bytes.NewBuffer(nil)
		cmd.Stdout = out
		cmd.Stderr = out
		cmd.Run()
		o, err = io.ReadAll(out)
		if err != nil {
			return
		}
		output := string(o)
		outs := []string{}
		for _, v1 := range strings.Split(output, "\n") {
			if v1 == "" {
				continue
			}
			outs = append(outs, strings.ReplaceAll(v1, "\r", ""))
		}
		ver := strings.Split(outs[0], " ")[2]
		ver = strings.ReplaceAll(ver, "\"", "")
		vers := strings.Split(ver, ".")[0]
		two := strings.Split(ver, ".")[1]
		if vers == "1" && two == "8" {
			vers = two
		}
		vs, ok := m[vers]
		if ok {
			if len(vs) < len(v) {
				continue
			}
		}
		m[vers] = v
	}
	return
}
func FmtJsonBytes(data []byte) ([]byte, error) {
	var s bytes.Buffer
	err := json.Indent(&s, data, "", "    ")
	if err != nil {
		return []byte(""), err
	}
	return s.Bytes(), nil
}
func WriteFmtJsonBytes(filename string, jsondata []byte) error {
	jsonb, err := FmtJsonBytes(jsondata)
	if err != nil {
		return err
	}
	return WriteBytes(filename, jsonb)
}
func IsRuleSameFrom_gjson_Result(rules gjson.Result, isdomouser bool) bool {
	if !rules.Exists() {
		return true
	}
	var notsame bool = false
	for _, vr := range rules.Array() {
		action := vr.Get("action")
		if action.Exists() {
			if action.String() != "allow" {
				continue
			}
		}
		features := vr.Get("features")
		if features.Exists() {
			features.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "is_demo_user" && value.Bool() == isdomouser {
					notsame = true
					return false
				}
				if key.String() == "has_custom_resolution" && value.Bool() {
					notsame = true
					return false
				}
				return true
			})
			if notsame {
				break
			}
		}
		oses := vr.Get("os")
		if oses.Exists() {
			oses.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "name" && value.String() == "windows" {
					notsame = true
					return false
				}
				if key.String() == "arch" && value.String() == infos.Arch {
					notsame = true
					return false
				}
				return true
			})
		}
	}
	return !notsame
}
func RunCmdFile(file string) (string, error) {
	cmd := exec.Command(filepath.Join(os.Getenv("windir"), "System32", "cmd.exe"), "/C", file)
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cmd.Wait()
	return string(out), nil
}
func RunCmdFileWithOutOutput(file string) error {
	cmd := exec.Command(filepath.Join(os.Getenv("windir"), "System32", "cmd.exe"), "/C", file)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

type AssertExpention struct{}

func (a *AssertExpention) Error() string {
	return "assertExpention"
}
func AssertsTrue(got bool) {
	if !got {
		panic(&AssertExpention{})
	}
}
func FilenameIsExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return os.IsNotExist(err)
}
func GetAllFilePath(p string, s []string) ([]string, error) {
	rd, err := os.ReadDir(p)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := filepath.Join(p, fi.Name())
			fullDir, err := filepath.Abs(fullDir)
			if err != nil {
				return s, err
			}
			s, err = GetAllFilePath(fullDir, s)
			if err != nil {
				return s, err
			}
		} else {
			fullName := filepath.Join(p, fi.Name())
			fullName, err := filepath.Abs(fullName)
			if err != nil {
				return s, err
			}
			s = append(s, fullName)
		}
	}
	return s, nil
}
func GetDllBit(p string) (string, error) {
	cmd := exec.Command("dumpbin.exe", "/headers", p)
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	o, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}
	return strings.Split(strings.Split(string(o), " machine (")[1], ")")[0], nil
}
func GetOsBit() string {
	return "x" + fmt.Sprintf("%d", 32<<(^uint(0)>>63))
}
