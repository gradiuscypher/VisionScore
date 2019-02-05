package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("command", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command Bar"
		v.Editable = true
		v.Wrap = false

		if _, err := setCurrentViewOnTop(g, "command"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("output", 0, 0, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "output"
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("output2", 0, maxY/2+1, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "output2"
		v.Editable = false
		v.Wrap = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func sendCommand(g *gocui.Gui, v *gocui.View) error {
	content := v.Buffer()
	out, err := g.View("output")
	if err != nil {
		return err
	}

	fmt.Fprintln(out, content)

	pos_x, _ := v.Cursor()
	v.MoveCursor(-pos_x, 0, true)
	v.Clear()
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, sendCommand); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
