package gui

import (
	"fyne.io/fyne/v2/app"
)

func AppRender() {
	a := app.New()
	w := a.NewWindow("Hello World")

	m := NewMap()
	w.SetContent(m)
	w.ShowAndRun()
}

//func MapRender() {
//	// action
//	w.SetContent(m)
//	// verify
//}
