package theme

import (
	"image/color"
	_"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

//go:embed .\fonts\SimSun-01.ttf
var fontsbyte []byte

func (m MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return fyne.NewStaticResource("SimSun-01.ttf", fontsbyte)
}
func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
