package main

import (
	"GoMCL2/messages"
	"GoMCL2/mylog"
	"GoMCL2/theme"
	"GoMCL2/weblib"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Infos struct {
	Arch string

	Workpath   string
	Mcdir      string
	Assetsdir  string
	Indexdir   string
	Objdir     string
	ObjBackdir string
	Libdir     string
	Verdir     string

	SourcesResult gjson.Result

	VersionManifestResult gjson.Result

	VersionManifestV2Result gjson.Result

	ConfigResult gjson.Result

	Logfile *os.File

	AssetsRootURL          string
	VersionManifestV2URL   string
	VersionManifestURL     string
	LibRootURL             string
	IndexRootURL           string
	VersionRootURL         string
	AuthlibInjectorRootURL string
	ClientJarRootURL       string

	NeedReplaceURL bool

	SourceName string
}

func needReplaceAssetsRootURL() bool {
	return infos.NeedReplaceURL && infos.AssetsRootURL != ""
}
func needReplaceLibRootURL() bool {
	return infos.NeedReplaceURL && infos.LibRootURL != ""
}
func needReplaceIndexRootURL() bool {
	return infos.NeedReplaceURL && infos.IndexRootURL != ""
}
func needReplaceClienJarRootURL() bool {
	return infos.NeedReplaceURL && infos.ClientJarRootURL != ""
}

var infos = &Infos{Arch: GetOsBit()}
var TheError error = nil
var (
	settingsIcon fyne.Resource
	mainwinIcon  fyne.Resource
)
var (
	//go:embed mainIcon.png
	mainwinIconBytes []byte

	//go:embed settingsIcon.png
	settingsIconBytes []byte
)

var (
	//go:embed sources.json
	defaultSources []byte
)

func init() {
	workpath, err := os.Getwd()
	if err != nil {
		TheError = err
		return
	}
	infos.Workpath = workpath
	logf, err := os.OpenFile("mcllog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		TheError = err
		return
	}
	infos.Logfile = logf
	b, err := os.ReadFile("SOURCE")
	if err != nil {
		logf.WriteString(fmt.Sprintf("[%s init/ERROR]: "+err.Error(), time.Now().String()))
		TheError = err
		return
	}
	infos.SourceName = string(b)
	var mcdir string = filepath.Join(workpath, ".minecraft")
	infos.Mcdir = mcdir
	assetsdir := filepath.Join(mcdir, "assets")
	infos.Assetsdir = assetsdir
	indexdir := filepath.Join(assetsdir, "indexes")
	infos.Indexdir = indexdir
	objdir := filepath.Join(assetsdir, "objects")
	infos.Objdir = objdir
	objbackdir := filepath.Join(assetsdir, "virtual", "legacy")
	infos.ObjBackdir = objbackdir
	libdir := filepath.Join(mcdir, "libraries")
	infos.Libdir = libdir
	verdir := filepath.Join(mcdir, "versions")
	infos.Verdir = verdir
	err = MakeDirs(infos.Mcdir, infos.Assetsdir, infos.Indexdir, infos.Objdir, infos.ObjBackdir, infos.Libdir, infos.Verdir)
	if err != nil {
		mylog.LogError(infos.Logfile, err, "init")
		TheError = err
		return
	}
	sources, err := os.ReadFile("sources.json")
	var isUseDefaultSources bool = false
	if err != nil {
		if os.IsNotExist(err) {
			sources = defaultSources
			isUseDefaultSources = true
		} else {
			TheError = err
			mylog.LogError(infos.Logfile, err, "init")
			return
		}
	}
	sourcesResult := gjson.ParseBytes(sources)
	defaultSourcesResult := gjson.ParseBytes(defaultSources)
	if !isUseDefaultSources {
		defaultSourcesResult.ForEach(func(key, value gjson.Result) bool {
			if !value.IsObject() {
				return false
			}
			valueMap := value.Map()
			valueMaps := map[string]string{}
			for k, v := range valueMap {
				valueMaps[k] = v.String()
			}
			var err error
			sources, err = sjson.SetBytes(sources, ReplaceByMap(key.String(), map[string]string{
				".": "\\.",
			}), valueMaps)
			if err != nil {
				TheError = err
				return false
			}
			return true
		})
		if TheError != nil {
			mylog.LogError(infos.Logfile, TheError, "init")
			return
		}
	}
	infos.SourcesResult = sourcesResult
	sourcessResult := sourcesResult.Get(ReplaceByMap(infos.SourceName, map[string]string{
		".": "\\.",
	}))
	if sourcessResult.Exists() {
		infos.AssetsRootURL = sourcessResult.Get("assetsRoot").String()
		mylog.LogInfo(infos.Logfile, "set assetRootURL: "+infos.AssetsRootURL, "init")
		infos.AuthlibInjectorRootURL = sourcessResult.Get("authlib-injectorRoot").String()
		mylog.LogInfo(infos.Logfile, "set authlib-injectorRootURL: "+infos.AuthlibInjectorRootURL, "init")
		infos.LibRootURL = sourcessResult.Get("libRoot").String()
		mylog.LogInfo(infos.Logfile, "set libRootURL: "+infos.LibRootURL, "init")
		infos.ClientJarRootURL = sourcessResult.Get("verjarRoot").String()
		mylog.LogInfo(infos.Logfile, "set client jar root URL: "+infos.ClientJarRootURL, "init")
		infos.IndexRootURL = sourcessResult.Get("indexRoot").String()
		mylog.LogInfo(infos.Logfile, "set indexRootURL: "+infos.IndexRootURL, "init")
		infos.VersionManifestURL = sourcessResult.Get("version_manifest").String()
		mylog.LogInfo(infos.Logfile, "set version_manifestURL: "+infos.VersionManifestURL, "init")
		infos.VersionManifestV2URL = sourcessResult.Get("version_manifest_v2").String()
		mylog.LogInfo(infos.Logfile, "set version_manifest_v2URL: "+infos.VersionManifestV2URL, "init")
		infos.NeedReplaceURL = sourcessResult.Get("needDo").Bool()
		mylog.LogInfo(infos.Logfile, fmt.Sprintf("need replace: %t", infos.NeedReplaceURL), "init")
		infos.VersionRootURL = sourcessResult.Get("versionRoot").String()
		mylog.LogInfo(infos.Logfile, "set version root URL: "+infos.VersionRootURL, "init")
	} else {
		infos.NeedReplaceURL = false
	}
	var verinfoURL string = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"
	if infos.NeedReplaceURL && infos.VersionManifestV2URL != "" {
		verinfoURL = infos.VersionManifestV2URL
	}
	mylog.LogInfo(infos.Logfile, "getting version_manifest_v2", "init")
	rb, err := weblib.GetBytesFromString(verinfoURL)
	if err != nil {
		TheError = err
		mylog.LogError(infos.Logfile, err, "init")
		return
	}
	mylog.LogInfo(infos.Logfile, "got version_manifest_v2", "init")
	infos.VersionManifestV2Result = gjson.ParseBytes(rb)
	settingsIcon = fyne.NewStaticResource("settingsIcon.png", settingsIconBytes)
	mainwinIcon = fyne.NewStaticResource("mainIcon.png", mainwinIconBytes)
}

const (
	Newline = "\n"
)

var nativePath string
var a fyne.App = app.New()
var mainSize = fyne.NewSize(600, 600)

func newSource() {
	addnewsourcewin := a.NewWindow("新建源向导")
	addnewsourcewin.Resize(mainSize)
	addnewsourcewin.SetIcon(settingsIcon)
	librooturlEntry := widget.NewEntry()
	verrooturlEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	addnewsourcewin.SetContent(container.NewVBox(
		widget.NewLabel("源名字:"),
		nameEntry,
		widget.NewLabel("库的根URL:"),
		librooturlEntry,
		widget.NewLabel("版本文件根URL:"),
		verrooturlEntry,
		widget.NewButton("确定", func() {
			name := nameEntry.Text
			namepath := ReplaceByMap(name, map[string]string{
				".": "\\.",
			})
			s, err := sjson.Set(infos.SourcesResult.String(), namepath+".libRoot", librooturlEntry)
			if err != nil {
				mylog.LogError(infos.Logfile, err, "add source")
				return
			}
			err = os.WriteFile("sources.json", []byte(s), 0666)
			if err != nil {
				mylog.LogError(infos.Logfile, err, "add source")
				return
			}
		}),
	))
	addnewsourcewin.Show()
}
func setting() {
	w := a.NewWindow("设置")
	w.Resize(mainSize)
	w.SetIcon(settingsIcon)
	sources := []string{}
	infos.SourcesResult.ForEach(func(key, value gjson.Result) bool {
		sources = append(sources, key.String())
		return true
	})
	sourceEntry := widget.NewSelectEntry(sources)
	sourceEntry.SetText(infos.SourceName)
	var hassave = false
	w.SetContent(container.NewVBox(
		widget.NewButton("新建源", newSource),
		widget.NewLabel("下载源:"),
		sourceEntry,
		widget.NewButton("保存", func() {
			err := WriteBytes("SOURCE", []byte(sourceEntry.Text))
			if err != nil {
				mylog.LogError(infos.Logfile, err, "set source")
				return
			}
			hassave = true
			w.Close()
		}),
	))
	w.SetOnClosed(func() {
		if hassave {
			return
		}
		save, err := messages.AskYesNo("Golang MCL", "是否保存设置")
		if err != nil {
			mylog.LogError(infos.Logfile, err, "settings")
			return
		}
		if !save {
			return
		}
		err = WriteBytes("SOURCE", []byte(sourceEntry.Text))
		if err != nil {
			mylog.LogError(infos.Logfile, err, "set source")
			return
		}
	})
	w.Show()
}
func main() {
	if TheError != nil {
		os.Exit(1)
	}
	defer infos.Logfile.Close()
	verjsons := infos.VersionManifestV2Result.Get("versions.#.id|@ugly").String()
	verjsons = strings.ReplaceAll(verjsons, "\"", "")
	verjsons = strings.TrimLeft(verjsons, "[")
	verjsons = strings.TrimRight(verjsons, "]")
	allVersions := strings.Split(verjsons, ",")
	a.Settings().SetTheme(&theme.MyTheme{})
	w := a.NewWindow("Golang Minecraft Launcher 2")
	w.SetIcon(mainwinIcon)
	w.Resize(mainSize)
	verchooseEntry := widget.NewSelectEntry(allVersions)
	outLabel := widget.NewLabel("")
	downloadfunc := func(b *widget.ProgressBarInfinite) {
		b.Show()
		b.Start()
		needVer := verchooseEntry.Text
		var (
			verUrl  string = ""
			verSha1 string = ""
			verJson gjson.Result
		)
		verUrlReuslt := infos.VersionManifestV2Result.Get("versions.#(id==\"" + needVer + "\").url")
		if !verUrlReuslt.Exists() {
			messages.ShowError("Golang MCL", "版本文件:\n无URL")
			return
		}
		verUrl = verUrlReuslt.String()
		verSha1Reuslt := infos.VersionManifestV2Result.Get("versions.#(id==\"" + needVer + "\").sha1")
		if !verSha1Reuslt.Exists() {
			messages.ShowError("Golang MCL", "版本文件:\n无Sha1")
		}
		verSha1 = verSha1Reuslt.String()
		if verUrl == "" || verSha1 == "" {
			messages.ShowError("Golang MCL", "无URL或sha1")
			return
		}
		var verDir string = filepath.Join(infos.Verdir, needVer)
		err := os.MkdirAll(verDir, 0666)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", err.Error())
			return
		}
		if infos.NeedReplaceURL && infos.VersionRootURL != "" {
			verUrl = ReplaceByMap(verUrl, map[string]string{
				"https://piston-meta.mojang.com":  infos.VersionRootURL,
				"https://launchermeta.mojang.com": infos.VersionRootURL,
				"https://launcher.mojang.com":     infos.VersionRootURL,
			})
		}
		rb, err := weblib.GetBytesFromString(verUrl)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", err.Error())
			return
		}
		gotsha1 := Sha1Bytes(rb)
		if gotsha1 != verSha1 {
			err = NewHashNotSame(verSha1, gotsha1)
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", "Get versionJson:\n"+err.Error())
			return
		}
		verJson = gjson.ParseBytes(rb)
		err = WriteFmtJsonBytes(filepath.Join(verDir, needVer+".json"), rb)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", "Write verJson:\n"+err.Error())
			return
		}
		i := verJson.Get("assetIndex")
		indexesB, err := GetIndexBytes(i)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", "Get assets index:\n"+err.Error())
			return
		}
		indexname := i.Get("id").String()
		indexpath := filepath.Join(infos.Indexdir, indexname+".json")
		err = WriteFmtJsonBytes(indexpath, indexesB)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", "Write indexes\n"+err.Error())
			return
		}
		indexes := gjson.ParseBytes(indexesB)
		cp, err := GetLib(verJson, outLabel)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", err.Error())
			return
		}
		w.Resize(fyne.NewSize(600, 600))
		err = GetObj(outLabel, indexes)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", err.Error())
			return
		}
		err = Backup(verJson, cp, false)
		if err != nil {
			mylog.LogError(infos.Logfile, err, "downloadpkg")
			messages.ShowError("Golang MCL", err.Error())
			return
		}
		outLabel.SetText("downloaded: " + needVer)
		b.Stop()
		b.Hide()
	}
	downloadbar := widget.NewProgressBarInfinite()
	downbox := container.NewVBox(
		verchooseEntry,
		widget.NewButton("         下载         ", func() {
			downloadfunc(downloadbar)
		}),
	)
	outbox := container.NewVBox(
		outLabel,
		downloadbar,
	)
	downloadbar.Hide()
	downloadBox := container.NewHBox(
		downbox,
		outbox,
	)
	userEntry := widget.NewEntry()
	passEntry := widget.NewPasswordEntry()
	emailEntry := widget.NewEntry()
	logintypeEntry := widget.NewSelectEntry([]string{
		"littleskin",
		"Mircosoft",
	})

	w.SetContent(container.NewVBox(
		downloadBox,
		widget.NewButton("启动", func() {
			needVer := verchooseEntry.Text
			if !FilenameIsExist(filepath.Join(infos.Verdir, needVer, needVer+".json")) {
				return
			}
			rb, err := os.ReadFile("backup.json")
			if err != nil {
				panic(err)
			}
			result := gjson.GetBytes(rb, ReplaceByMap(needVer, map[string]string{
				".": "\\.",
			}))
			if !result.Exists() {
				messages.ShowError("Golang MCL", "该版本未下载，\n请下载后重试")
				return
			}
			cmds := result.Get("cmd").String()
			usernames, uuid, token, types, err := Login(logintypeEntry, userEntry, passEntry, emailEntry, int(result.Get("complianceLevel").Int()))
			if err != nil {
				panic(err)
			}
			verJsons, err := os.ReadFile(filepath.Join(infos.Verdir, needVer, needVer+".json"))
			if err != nil {
				return
			}
			indexPath := gjson.GetBytes(verJsons, "assetIndex.id").String() + ".json"
			indexJson, err := os.ReadFile(filepath.Join(infos.Indexdir, indexPath))
			if err != nil {
				return
			}
			err = GetObj(outLabel, gjson.ParseBytes(indexJson))
			if err != nil {
				return
			}
			_ = os.Remove("start.bat")
			f, err := os.Create("start.bat")
			if err != nil {
				mylog.LogError(infos.Logfile, err, "launchpkg")
				return
			}
			cmds = ReplaceByMap(cmds, map[string]string{
				"${auth_player_name}":  usernames,
				"${auth_uuid}":         uuid,
				"${auth_access_token}": token,
				"${user_type}":         types,
				"${resolution_width}":  "854",
				"${resolution_height}": "480",
			})
			if logintypeEntry.Text == "littleskin" {
				p, err := Getauthlib()
				if err != nil {
					mylog.LogError(infos.Logfile, err, "loginpkg")
					return
				}
				cmds = ReplaceByMap(cmds, map[string]string{
					"-XX:HeapDumpPath=MojangTricksIntelDriversForPerformance_javaw.exe_minecraft.exe.heapdump": "-javaagent:\"" + p + "\"=https://littleskin.cn/api/yggdrasil -XX:HeapDumpPath=MojangTricksIntelDriversForPerformance_javaw.exe_minecraft.exe.heapdump",
				})
			}
			_, err = f.WriteString(cmds)
			if err != nil {
				f.Close()
				mylog.LogError(infos.Logfile, err, "launchpkg")
				return
			}
			f.Close()
			err = RunCmdFileWithOutOutput("start.bat")
			if err != nil {
				mylog.LogError(infos.Logfile, err, "launchpkg")
				return
			}
		}),
		widget.NewLabel("username:"),
		userEntry,
		widget.NewLabel("password:"),
		passEntry,
		widget.NewLabel("e-mail:"),
		emailEntry,
		widget.NewLabel("type:"),
		logintypeEntry,
		widget.NewButton("设置", setting),
		widget.NewButton("清空日志", func() {
			infos.Logfile.Close()
			logf, err := os.OpenFile("mcllog.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				mylog.LogError(infos.Logfile, err, "launchpkg")
				return
			}
			infos.Logfile = logf
		}),
	))
	w.SetOnClosed(func() {
		a.Quit()
	})
	w.Show()
	a.Run()
}
