// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"

	"github.com/goki/gi"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/driver"
	"github.com/goki/gi/units"
	"github.com/goki/ki"
)

func main() {
	driver.Main(func(app oswin.App) {
		mainrun()
	})
}

var CurFilename = ""

func mainrun() {
	width := 1024
	height := 768

	// turn this on to see a trace of the rendering
	// gi.Update2DTrace = true
	// gi.Render2DTrace = true
	// gi.Layout2DTrace = true

	rec := ki.Node{}          // receiver for events
	rec.InitName(&rec, "rec") // this is essential for root objects not owned by other Ki tree nodes

	win := gi.NewWindow2D("GoGi SVG Test Window", width, height, true)
	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	vp.Fill = true

	vlay := vp.AddNewChild(gi.KiT_Frame, "vlay").(*gi.Frame)
	vlay.Lay = gi.LayoutCol

	brow := vlay.AddNewChild(gi.KiT_Layout, "brow").(*gi.Layout)
	brow.Lay = gi.LayoutRow
	brow.SetStretchMaxWidth()

	svgrow := vlay.AddNewChild(gi.KiT_Layout, "svgrow").(*gi.Layout)
	svgrow.Lay = gi.LayoutRow
	svgrow.SetProp("align-vert", gi.AlignMiddle)
	svgrow.SetProp("align-horiz", "center")
	svgrow.SetProp("margin", 2.0) // raw numbers = px = 96 dpi pixels
	svgrow.SetStretchMaxWidth()
	svgrow.SetStretchMaxHeight()

	svg := svgrow.AddNewChild(gi.KiT_SVG, "svg").(*gi.SVG)
	svg.Fill = true
	svg.SetStretchMaxWidth()
	svg.SetStretchMaxHeight()

	loads := brow.AddNewChild(gi.KiT_Button, "loadsvg").(*gi.Button)
	// loads.SetProp("vertical-align", gi.AlignMiddle)
	loads.SetText("Load SVG")

	fnm := brow.AddNewChild(gi.KiT_TextField, "cur-fname").(*gi.TextField)
	fnm.SetMinPrefWidth(units.NewValue(20, units.Em))
	// fnm.SetProp("vertical-align", AlignMiddle)

	loads.ButtonSig.Connect(win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			path, fn := filepath.Split(CurFilename)
			gi.FileViewDialog(vp, path, fn, "Load SVG", "", win, func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					dlg, _ := send.(*gi.Dialog)
					CurFilename := gi.FileViewDialogValue(dlg)
					fnm.SetText(CurFilename)
					svg.LoadXML(CurFilename)
					svg.SetNormXForm()
				}
			})
		}
	})

	fnm.TextFieldSig.Connect(win.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.TextFieldDone) {
			tf := send.(*gi.TextField)
			CurFilename = tf.Text()
			svg.LoadXML(CurFilename)
			svg.SetNormXForm()
		}
	})

	vp.UpdateEndNoSig(updt)

	win.StartEventLoop()
}
