package viewport

import (
	"fmt"

	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/config"
	"github.com/skycoin/viscript/hypervisor/input/mouse"
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/gl"
	t "github.com/skycoin/viscript/viewport/terminal"
)

var currentTerminalModification int

const (
	TermMod_None = iota
	TermMod_Moving
	TermMod_ResizingX
	TermMod_ResizingY
	TermMod_ResizingBoth
)

// triggered both by moving **AND*** by pressing buttons
func onMouseCursorPos(m msg.MessageMousePos) {
	if config.DebugPrintInputEvents() {
		fmt.Print("msg.TypeMousePos")
		showFloat64("X", m.X)
		showFloat64("Y", m.Y)
		println()
	}

	mouse.Update(app.Vec2F{float32(m.X), float32(m.Y)})

	foc := t.Terms.GetFocusedTerminal()

	if foc == nil {
		println("onMouseCursorPos()   -   foc == nil")
		return
	}

	setPointerBasedOnPosition()

	switch currentTerminalModification {

	case TermMod_Moving:
		//high resolution delta for potentially subpixel precision resizing
		delt := mouse.GlPos.GetDeltaFrom(mouse.PrevGlPos)
		t.Terms.MoveFocusedTerminal(delt, &mouse.DeltaSinceLeftClick)
		//gl.SetHandPointer()

	case TermMod_ResizingX:
		foc.ResizeHorizontally(mouse.GlPos.X)
	case TermMod_ResizingY:
		foc.ResizeVertically(mouse.GlPos.Y)

	case TermMod_ResizingBoth:
		foc.ResizeHorizontally(mouse.GlPos.X)
		foc.ResizeVertically(mouse.GlPos.Y)

	}
}

func onMouseScroll(m msg.MessageMouseScroll) {
	if config.DebugPrintInputEvents() {
		print("msg.TypeMouseScroll")
		showFloat64("X Offset", m.X)
		showFloat64("Y Offset", m.Y)
		showBool("HoldingAlt", m.HoldingAlt)
		showBool("HoldingControl", m.HoldingControl)
		println()
	}
}

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(m msg.MessageMouseButton) {
	if config.DebugPrintInputEvents() {
		fmt.Print("msg.TypeMouseButton")
		showUInt8("Button", m.Button)
		showUInt8("Action", m.Action)
		showUInt8("Mod", m.Mod)
		println()
	}

	convertClickToTextCursorPosition(m.Button, m.Action)

	if msg.Action(m.Action) == msg.Press {
		switch msg.MouseButton(m.Button) {

		case msg.MouseButtonLeft:
			mouse.LeftButtonIsDown = true
			mouse.DeltaSinceLeftClick = app.Vec2F{0, 0}

			// // detect clicks in rects
			// if mouse.PointerIsInside(ui.MainMenu.Rect) {
			// 	respondToAnyMenuButtonClicks()
			// } else { // respond to any panel clicks outside of menu
			focusOnTopmostRectThatContainsPointer()
			// }

			currentTerminalModification = getTerminalModificationByZone()

		}
	} else if msg.Action(m.Action) == msg.Release {
		switch msg.MouseButton(m.Button) {

		case msg.MouseButtonLeft:
			mouse.LeftButtonIsDown = false
			currentTerminalModification = TermMod_None

		}
	}
}

func setPointerBasedOnPosition() {
	foc := t.Terms.GetFocusedTerminal()

	if foc == nil {
		gl.SetArrowPointer()
	} else {
		if !foc.FixedSize {
			if mouse.NearRight(foc.Bounds) &&
				mouse.NearBottom(foc.Bounds) {

				gl.SetCornerResizePointer()
				return
			}

			if mouse.NearRight(foc.Bounds) {
				gl.SetHResizePointer()
				return
			}

			if mouse.NearBottom(foc.Bounds) {
				gl.SetVResizePointer()
				return
			}
		}

		if mouse.PointerIsInside(foc.Bounds) {
			gl.SetHandPointer()
			//gl.SetIBeamPointer() //IBeam is harder to see...
			//...and probably only makes sense when you can click to type anywhere on screen.
			//Terminals currently limit the animated cursor position to be within the last 2 lines of onscreen text
		} else {
			gl.SetArrowPointer()
		}
	}
}

func focusOnTopmostRectThatContainsPointer() {
	var topmostZ float32
	var topmostId msg.TerminalId

	for _, t := range t.Terms.TermMap {
		if mouse.PointerIsInside(t.Bounds) {
			if topmostZ < t.Depth {
				topmostZ = t.Depth
				topmostId = t.TerminalId
			}
		}
	}

	if topmostZ > 0 {
		t.Terms.SetFocused(topmostId)
	}
}

func convertClickToTextCursorPosition(button, action uint8) {
	// if glfw.MouseButton(button) == glfw.MouseButtonLeft &&
	// 	glfw.Action(action) == glfw.Press {

	// 	foc := Focused

	// 	if foc.IsEditable && foc.Content.Contains(mouse.GlX, mouse.GlY) {
	// 		if foc.MouseY < len(foc.TextBodies[0]) {
	// 			foc.CursY = foc.MouseY

	// 			if foc.MouseX <= len(foc.TextBodies[0][foc.CursY]) {
	// 				foc.CursX = foc.MouseX
	// 			} else {
	// 				foc.CursX = len(foc.TextBodies[0][foc.CursY])
	// 			}
	// 		} else {
	// 			foc.CursY = len(foc.TextBodies[0]) - 1
	// 		}
	// 	}
	// }
}

func respondToAnyMenuButtonClicks() {
	// for _, bu := range ui.MainMenu.Buttons {
	// 	if mouse.PointerIsInside(bu.Rect.Rectangle) {
	// 		bu.Activated = !bu.Activated

	// 		switch bu.Name {
	// 		case "Run":
	// 			if bu.Activated {
	// 				//script.Digest(true)
	// 			}
	// 			break
	// 		case "Testing Tree":
	// 			if bu.Activated {
	// 				//script.Digest(true)
	// 				//script.MakeTree()
	// 			} else { // deactivated
	// 				// remove all terminals with trees
	// 				b := t.Terms[:0]
	// 				for _, t := range t.Terms {
	// 					if len(t.Trees) < 1 {
	// 						b = append(b, t)
	// 					}
	// 				}
	// 				t.Terms = b
	// 				//fmt.Printf("len of b (from Terms) after removing ones with trees: %d\n", len(b))
	// 				//fmt.Printf("len of Terms: %d\n", len(Terms))
	// 			}
	// 			break
	// 		}

	// 		app.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
	// 	}
	// }
}

// the rest of these funcs are almost identical, just top 2 vars customized (and string format)
func showBool(s string, x bool) {
	fmt.Printf("   [%s: %t]", s, x)
}

func showUInt8(s string, x uint8) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showSInt32(s string, x int32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showUInt32(s string, x uint32) {
	fmt.Printf("   [%s: %d]", s, x)
}

func showFloat64(s string, f float64) {
	fmt.Printf("   [%s: %.1f]", s, f)
}

func getTerminalModificationByZone() int {
	foc := t.Terms.GetFocusedTerminal()

	if foc == nil {
		println("onMouseCursorPos()   -   foc == nil")
		return TermMod_None
	}

	if !foc.FixedSize {
		if mouse.NearRight(foc.Bounds) &&
			mouse.NearBottom(foc.Bounds) {

			return TermMod_ResizingBoth
		}

		if /****/ mouse.NearRight(foc.Bounds) {
			return TermMod_ResizingX
		} else if mouse.NearBottom(foc.Bounds) {
			return TermMod_ResizingY
		}
	}

	if mouse.PointerIsInside(foc.Bounds) {
		return TermMod_Moving
	} else {
		return TermMod_None
	}
}
