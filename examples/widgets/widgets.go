// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"

	"github.com/goki/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/units"
	"github.com/goki/ki"
	"github.com/goki/ki/kit"
)

func main() {
	gimain.Main(func() {
		mainrun()
	})
}

func mainrun() {
	width := 1024
	height := 768

	// turn these on to see a traces of various stages of processing..
	// gi.Update2DTrace = true
	// gi.Render2DTrace = true
	// gi.Layout2DTrace = true
	// ki.SignalTrace = true

	rec := ki.Node{}          // receiver for events
	rec.InitName(&rec, "rec") // this is essential for root objects not owned by other Ki tree nodes

	win := gi.NewWindow2D("gogi-widgets-demo", "GoGi Widgets Demo", width, height, true) // true = pixel sizes

	icnm := "widget-wedge-down"

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	// style sheet
	var css = ki.Props{
		"button": ki.Props{
			"background-color": gi.Color{255, 240, 240, 255},
		},
		"#combo": ki.Props{
			"background-color": gi.Color{240, 255, 240, 255},
		},
		".hslides": ki.Props{
			"background-color": gi.Color{240, 225, 255, 255},
		},
		"kbd": ki.Props{
			"color": "blue",
		},
	}
	vp.CSS = css

	mfr := win.SetMainFrame()
	// mfr.SetProp("background-color", "linear-gradient(to top, red, lighter-80)")
	// mfr.SetProp("background-color", "linear-gradient(to right, red, orange, yellow, green, blue, indigo, violet)")
	// mfr.SetProp("background-color", "linear-gradient(to right, rgba(255,0,0,0), rgba(255,0,0,1))")
	// mfr.SetProp("background-color", "radial-gradient(red, lighter-80)")

	trow := mfr.AddNewChild(gi.KiT_Layout, "trow").(*gi.Layout)
	trow.Lay = gi.LayoutHoriz
	trow.SetStretchMaxWidth()

	trow.AddNewChild(gi.KiT_Stretch, "str1")
	title := trow.AddNewChild(gi.KiT_Label, "title").(*gi.Label)
	title.Text = `This is a <b>demonstration</b> of the
<span style="color:red">various</span> <i>GoGi</i> Widgets<br>
<large>Shortcuts: <kbd>Ctrl+Alt+P</kbd> = Preferences,
<kbd>Ctrl+Alt+E</kbd> = Editor, <kbd>Ctrl/Cmd +/-</kbd> = zoom</large>`
	title.SetProp("text-align", gi.AlignCenter)
	title.SetProp("vertical-align", gi.AlignTop)
	title.SetProp("font-family", "Times New Roman, serif")
	title.SetProp("font-size", "x-large")
	// title.SetProp("letter-spacing", 2)
	title.SetProp("line-height", 1.5)
	trow.AddNewChild(gi.KiT_Stretch, "str2")

	//////////////////////////////////////////
	//      Buttons

	mfr.AddNewChild(gi.KiT_Space, "blspc")
	blrow := mfr.AddNewChild(gi.KiT_Layout, "blrow").(*gi.Layout)
	blab := blrow.AddNewChild(gi.KiT_Label, "blab").(*gi.Label)
	blab.SetProp("font-family", "Arial Unicode")
	blab.Text = "Buttons: ⌘ "
	blab.Selectable = true

	brow := mfr.AddNewChild(gi.KiT_Layout, "brow").(*gi.Layout)
	brow.Lay = gi.LayoutHoriz
	brow.SetProp("horizontal-align", gi.AlignLeft)
	// brow.SetProp("horizontal-align", gi.AlignJustify)
	brow.SetStretchMaxWidth()

	button1 := brow.AddNewChild(gi.KiT_Button, "button1").(*gi.Button)
	button1.SetProp("#icon", ki.Props{ // note: must come before SetIcon
		"width":  units.NewValue(1.5, units.Em),
		"height": units.NewValue(1.5, units.Em),
	})
	button1.Tooltip = "press this <i>button</i> to pop up a dialog box"

	button1.SetIcon(icnm)
	button1.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
		if sig == int64(gi.ButtonClicked) { // note: 3 diff ButtonSig sig's possible -- important to check
			// vp.Win.Quit()
			gi.StringPromptDialog(vp, "Enter value here..", "Button1 Dialog", "This is a string prompt dialog!  Various specific types of dialogs are available.", nil, rec.This,
				func(recv, send ki.Ki, sig int64, data interface{}) {
					dlg := send.(*gi.Dialog)
					if sig == int64(gi.DialogAccepted) {
						val := gi.StringPromptDialogValue(dlg)
						fmt.Printf("got string value: %v\n", val)
					}
				})
		}
	})

	button2 := brow.AddNewChild(gi.KiT_Button, "button2").(*gi.Button)
	button2.SetText("Open GoGiEditor")
	// button2.SetProp("background-color", "#EDF")
	button2.Tooltip = "This button will open the GoGi GUI editor where you can edit this very GUI and see it update dynamically as you change things"
	button2.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
		if sig == int64(gi.ButtonClicked) {
			giv.GoGiEditor(vp)
		}
	})

	checkbox := brow.AddNewChild(gi.KiT_CheckBox, "checkbox").(*gi.CheckBox)
	checkbox.Text = "Toggle"

	// note: receiver for menut items with shortcuts must be a Node2D or Window
	mb1 := brow.AddNewChild(gi.KiT_MenuButton, "menubutton1").(*gi.MenuButton)
	mb1.SetText("Menu Button")
	mb1.Menu.AddMenuText("Menu Item 1", "Shift+Control+1", win.This, 1, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received menu action data: %v from menu action: %v\n", data, send.Name())
	})

	mi2 := mb1.Menu.AddMenuText("Menu Item 2", "", nil, 2, nil)

	mi2.Menu.AddMenuText("Sub Menu Item 2", "", win.This, 2.1, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received menu action data: %v from menu action: %v\n", data, send.Name())
	})

	mb1.Menu.AddSeparator("sep1")

	mb1.Menu.AddMenuText("Menu Item 3", "Control+3", win.This, 3, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received menu action data: %v from menu action: %v\n", data, send.Name())
	})

	brow.SetPropChildren("margin", units.NewValue(2, units.Ex))

	//////////////////////////////////////////
	//      Sliders

	mfr.AddNewChild(gi.KiT_Space, "slspc")
	slrow := mfr.AddNewChild(gi.KiT_Layout, "slrow").(*gi.Layout)
	slab := slrow.AddNewChild(gi.KiT_Label, "slab").(*gi.Label)
	slab.Text = "Sliders:"

	srow := mfr.AddNewChild(gi.KiT_Layout, "srow").(*gi.Layout)
	srow.Lay = gi.LayoutHoriz
	srow.SetProp("horizontal-align", "left")
	srow.SetStretchMaxWidth()

	slider1 := srow.AddNewChild(gi.KiT_Slider, "slider1").(*gi.Slider)
	slider1.Dim = gi.X
	slider1.Class = "hslides"
	slider1.Defaults()
	slider1.SetMinPrefWidth(units.NewValue(20, units.Em))
	slider1.SetMinPrefHeight(units.NewValue(2, units.Em))
	slider1.SetValue(0.5)
	slider1.Snap = true
	slider1.Tracking = true
	slider1.Icon = gi.IconName("widget-circlebutton-on")

	slider2 := srow.AddNewChild(gi.KiT_Slider, "slider2").(*gi.Slider)
	slider2.Dim = gi.Y
	slider2.Defaults()
	slider2.SetMinPrefHeight(units.NewValue(10, units.Em))
	slider2.SetMinPrefWidth(units.NewValue(1, units.Em))
	slider2.SetStretchMaxHeight()
	slider2.SetValue(0.5)

	slider1.SliderSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received slider signal: %v from slider: %v with data: %v\n", gi.SliderSignals(sig), send.Name(), data)
	})

	slider2.SliderSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received slider signal: %v from slider: %v with data: %v\n", gi.SliderSignals(sig), send.Name(), data)
	})

	scrollbar1 := srow.AddNewChild(gi.KiT_ScrollBar, "scrollbar1").(*gi.ScrollBar)
	scrollbar1.Dim = gi.X
	scrollbar1.Class = "hslides"
	scrollbar1.Defaults()
	scrollbar1.SetMinPrefWidth(units.NewValue(20, units.Em))
	scrollbar1.SetMinPrefHeight(units.NewValue(1, units.Em))
	scrollbar1.SetThumbValue(0.25)
	scrollbar1.SetValue(0.25)
	// scrollbar1.Snap = true
	scrollbar1.Tracking = true
	scrollbar1.SliderSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received scrollbar signal: %v from scrollbar: %v with data: %v\n", gi.SliderSignals(sig), send.Name(), data)
	})

	scrollbar2 := srow.AddNewChild(gi.KiT_ScrollBar, "scrollbar2").(*gi.ScrollBar)
	scrollbar2.Dim = gi.Y
	scrollbar2.Defaults()
	scrollbar2.SetMinPrefHeight(units.NewValue(10, units.Em))
	scrollbar2.SetMinPrefWidth(units.NewValue(1, units.Em))
	scrollbar2.SetStretchMaxHeight()
	scrollbar2.SetThumbValue(0.1)
	scrollbar2.SetValue(0.5)
	scrollbar2.SliderSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received scrollbar signal: %v from scrollbar: %v with data: %v\n", gi.SliderSignals(sig), send.Name(), data)
	})

	//////////////////////////////////////////
	//      Text Widgets

	mfr.AddNewChild(gi.KiT_Space, "tlspc")
	txlrow := mfr.AddNewChild(gi.KiT_Layout, "txlrow").(*gi.Layout)
	txlab := txlrow.AddNewChild(gi.KiT_Label, "txlab").(*gi.Label)
	txlab.Text = "Text Widgets:"
	txrow := mfr.AddNewChild(gi.KiT_Layout, "txrow").(*gi.Layout)
	txrow.Lay = gi.LayoutHoriz
	// txrow.SetProp("horizontal-align", gi.AlignJustify)
	txrow.SetStretchMaxWidth()

	edit1 := txrow.AddNewChild(gi.KiT_TextField, "edit1").(*gi.TextField)
	edit1.SetText("Edit this text")
	edit1.SetProp("min-width", "20em")
	edit1.TextFieldSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("Received line edit signal: %v from edit: %v with data: %v\n", gi.TextFieldSignals(sig), send.Name(), data)
	})
	// edit1.SetProp("inactive", true)
	edit1.UseCmpltr = true

	sb := txrow.AddNewChild(gi.KiT_SpinBox, "spin").(*gi.SpinBox)
	sb.Defaults()
	sb.HasMin = true
	sb.Min = 0.0
	sb.SpinBoxSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("SpinBox %v value changed: %v\n", send.Name(), data)
	})

	cb := txrow.AddNewChild(gi.KiT_ComboBox, "combo").(*gi.ComboBox)
	cb.ItemsFromTypes(kit.Types.AllImplementersOf(reflect.TypeOf((*gi.Node2D)(nil)).Elem(), false), true, true, 50)
	cb.ComboSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		fmt.Printf("ComboBox %v selected index: %v data: %v\n", send.Name(), sig, data)
	})

	// txrow.SetPropChildren("margin", units.NewValue(2, units.Ex))

	// main menu

	mmen := win.MainMenu
	mmen.ConfigMenus([]string{"File", "Edit", "Window"})

	emen := win.MainMenu.KnownChildByName("Edit", 1).(*gi.Action)
	emen.Menu = make(gi.Menu, 0, 10)
	emen.Menu.AddCopyCutPaste(win, true)

	vp.UpdateEndNoSig(updt)

	win.StartEventLoop()

	// note: never gets here..
	fmt.Printf("ending\n")
}
