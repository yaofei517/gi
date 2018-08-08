// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"github.com/goki/gi/units"
	"github.com/goki/ki"
	"github.com/goki/ki/kit"
)

////////////////////////////////////////////////////////////////////////////////////////
// MenuBar

// MenuBar is a Layout (typically LayoutHoriz) that renders a gradient
// background and has convenience methods for adding menus
type MenuBar struct {
	Layout
	MainMenu bool `desc:"is this the main menu bar for a window?  controls whether displayed on macOS"`
}

var KiT_MenuBar = kit.Types.AddType(&MenuBar{}, MenuBarProps)

var MenuBarProps = ki.Props{
	"padding": units.NewValue(2, units.Px),
	"margin":  units.NewValue(0, units.Px),
	// "spacing":          units.NewValue(2, units.Px),
	"color":            &Prefs.FontColor,
	"background-color": "linear-gradient(pref(ControlColor), highlight-10)",
}

// MenuBarStdRender does the standard rendering of the bar
func (mb *MenuBar) MenuBarStdRender() {
	st := &mb.Sty
	rs := &mb.Viewport.Render
	pc := &rs.Paint

	pos := mb.LayData.AllocPos
	sz := mb.LayData.AllocSize
	pc.FillBox(rs, pos, sz, &st.Font.BgColor)
}

func (mb *MenuBar) Render2D() {
	if len(mb.Kids) == 0 { // todo: check for mac menu and don't render -- also need checks higher up
		return
	}
	if mb.FullReRenderIfNeeded() {
		return
	}
	if mb.PushBounds() {
		mb.MenuBarStdRender()
		mb.LayoutEvents()
		mb.RenderScrolls()
		mb.Render2DChildren()
		mb.PopBounds()
	} else {
		mb.DisconnectAllEvents(AllPris) // uses both Low and Hi
	}
}

// ConfigMenus configures Action items as children of MenuBar with the given
// names, which function as the main menu panels for the menu bar (File, Edit,
// etc).  Access the resulting menus as .KnownChildByName("name").(*Action),
// and
func (mb *MenuBar) ConfigMenus(menus []string) {
	sz := len(menus)
	tnl := make(kit.TypeAndNameList, sz+1)
	typ := KiT_Action
	for i, m := range menus {
		tnl[i].Type = typ
		tnl[i].Name = m
	}
	tnl[sz].Type = KiT_Stretch
	tnl[sz].Name = "menstr"
	_, updt := mb.ConfigChildren(tnl, false)
	for i, m := range menus {
		ma := mb.Kids[i].(*Action)
		ma.SetText(m)
		ma.SetAsMenu()
	}
	mb.UpdateEnd(updt)
}

////////////////////////////////////////////////////////////////////////////////////////
// ToolBar

// ToolBar is a Layout (typically LayoutHoriz) that renders a gradient
// background and is useful for holding Actions that do things
type ToolBar struct {
	Layout
}

var KiT_ToolBar = kit.Types.AddType(&ToolBar{}, ToolBarProps)

var ToolBarProps = ki.Props{
	"padding":          units.NewValue(2, units.Px),
	"margin":           units.NewValue(0, units.Px),
	"color":            &Prefs.FontColor,
	"background-color": "linear-gradient(pref(ControlColor), highlight-10)",
}

// ToolBarStdRender does the standard rendering of the bar
func (mb *ToolBar) ToolBarStdRender() {
	st := &mb.Sty
	rs := &mb.Viewport.Render
	pc := &rs.Paint

	pos := mb.LayData.AllocPos
	sz := mb.LayData.AllocSize
	pc.FillBox(rs, pos, sz, &st.Font.BgColor)
}

func (mb *ToolBar) Render2D() {
	if len(mb.Kids) == 0 { // todo: check for mac menu and don't render -- also need checks higher up
		return
	}
	if mb.FullReRenderIfNeeded() {
		return
	}
	if mb.PushBounds() {
		mb.ToolBarStdRender()
		mb.LayoutEvents()
		mb.RenderScrolls()
		mb.Render2DChildren()
		mb.PopBounds()
	} else {
		mb.DisconnectAllEvents(AllPris) // uses both Low and Hi
	}
}
