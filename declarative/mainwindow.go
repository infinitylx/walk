// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package declarative

import "github.com/lxn/walk"

type MainWindow struct {
	AssignTo         **walk.MainWindow
	Name             string
	Enabled          Property
	Visible          Property
	Font             Font
	MinSize          Size
	MaxSize          Size
	ContextMenuItems []MenuItem
	OnKeyDown        walk.KeyEventHandler
	OnKeyPress       walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnDropFiles      walk.DropFilesEventHandler
	OnSizeChanged    walk.EventHandler
	Icon             Property
	Title            Property
	Size             Size
	DataBinder       DataBinder
	Layout           Layout
	Children         []Widget
	MenuItems        []MenuItem
	StatusBarItems   []StatusBarItem
	ToolBarItems     []MenuItem // Deprecated: use ToolBar instead
	ToolBar          ToolBar
	Expressions      func() map[string]walk.Expression
	Functions        map[string]func(args ...interface{}) (interface{}, error)
}

func (mw MainWindow) Create() error {
	w, err := walk.NewMainWindow()
	if err != nil {
		return err
	}

	tlwi := topLevelWindowInfo{
		Name:             mw.Name,
		Enabled:          mw.Enabled,
		Visible:          mw.Visible,
		Font:             mw.Font,
		ToolTipText:      "",
		MinSize:          mw.MinSize,
		MaxSize:          mw.MaxSize,
		ContextMenuItems: mw.ContextMenuItems,
		OnKeyDown:        mw.OnKeyDown,
		OnKeyPress:       mw.OnKeyPress,
		OnKeyUp:          mw.OnKeyUp,
		OnMouseDown:      mw.OnMouseDown,
		OnMouseMove:      mw.OnMouseMove,
		OnMouseUp:        mw.OnMouseUp,
		OnSizeChanged:    mw.OnSizeChanged,
		DataBinder:       mw.DataBinder,
		Layout:           mw.Layout,
		Children:         mw.Children,
		Icon:             mw.Icon,
		Title:            mw.Title,
	}

	builder := NewBuilder(nil)

	w.SetSuspended(true)
	builder.Defer(func() error {
		w.SetSuspended(false)
		return nil
	})

	builder.deferBuildMenuActions(w.Menu(), mw.MenuItems)

	return builder.InitWidget(tlwi, w, func() error {
		if len(mw.ToolBar.Items) > 0 {
			var tb *walk.ToolBar
			if mw.ToolBar.AssignTo == nil {
				mw.ToolBar.AssignTo = &tb
			}

			if err := mw.ToolBar.Create(builder); err != nil {
				return err
			}

			old := w.ToolBar()
			w.SetToolBar(*mw.ToolBar.AssignTo)
			old.Dispose()
		} else {
			builder.deferBuildActions(w.ToolBar().Actions(), mw.ToolBarItems)
		}

		for _, sbi := range mw.StatusBarItems {
			s := walk.NewStatusBarItem()
			if sbi.AssignTo != nil {
				*sbi.AssignTo = s
			}
			s.SetIcon(sbi.Icon)
			s.SetText(sbi.Text)
			s.SetToolTipText(sbi.ToolTipText)
			if sbi.Width > 0 {
				s.SetWidth(sbi.Width)
			}
			if sbi.OnClicked != nil {
				s.Clicked().Attach(sbi.OnClicked)
			}
			w.StatusBar().Items().Add(s)
			w.StatusBar().SetVisible(true)
		}

		if err := w.SetSize(mw.Size.toW()); err != nil {
			return err
		}

		imageList, err := walk.NewImageList(walk.Size{16, 16}, 0)
		if err != nil {
			return err
		}
		w.ToolBar().SetImageList(imageList)

		if mw.OnDropFiles != nil {
			w.DropFiles().Attach(mw.OnDropFiles)
		}

		if mw.AssignTo != nil {
			*mw.AssignTo = w
		}

		if mw.Expressions != nil {
			for name, expr := range mw.Expressions() {
				builder.expressions[name] = expr
			}
		}
		if mw.Functions != nil {
			for name, fn := range mw.Functions {
				builder.functions[name] = fn
			}
		}

		builder.Defer(func() error {
			if mw.Visible != false {
				w.Show()
			}

			return nil
		})

		return nil
	})
}

func (mw MainWindow) Run() (int, error) {
	var w *walk.MainWindow

	if mw.AssignTo == nil {
		mw.AssignTo = &w
	}

	if err := mw.Create(); err != nil {
		return 0, err
	}

	return (*mw.AssignTo).Run(), nil
}

type StatusBarItem struct {
	AssignTo    **walk.StatusBarItem
	Icon        *walk.Icon
	Text        string
	ToolTipText string
	Width       int
	OnClicked   walk.EventHandler
}
