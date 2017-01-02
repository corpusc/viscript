package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/gfx"
	"github.com/corpusc/viscript/msg"
	"github.com/corpusc/viscript/script"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// apparently every time this is fired, a mouse position event is ALSO fired
func onMouseButton(
	w *glfw.Window,
	b glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	if action == glfw.Press {
		switch glfw.MouseButton(b) {
		case glfw.MouseButtonLeft:
			// respond to clicks in ui rectangles
			if gfx.MouseCursorIsInside(ui.MainMenu.Rect) {
				respondToAnyMenuButtonClicks()
			} else { // respond to any panel clicks outside of menu
				for _, pan := range gfx.Rend.Panels {
					if pan.ContainsMouseCursor() {
						pan.RespondToMouseClick()
					}
				}
			}
		}
	}

	// build message
	//content := append(getByteOfUInt8(uint8(b)), getByteOfUInt8(uint8(action))...)
	//content = append(content, getByteOfUInt8(uint8(mod))...)
	//dispatchWithPrefix(content, msg.TypeMouseButton)

	//MessageMouseButton
	var m msg.MessageMouseButton
	m.Button = uint8(b)
	m.Action = uint8(action)
	m.Mod = uint8(mod)
	DispatchEvent(msg.TypeMouseButton, m)
}

func respondToAnyMenuButtonClicks() {
	for _, bu := range ui.MainMenu.Buttons {
		if gfx.MouseCursorIsInside(bu.Rect) {
			bu.Activated = !bu.Activated

			switch bu.Name {
			case "Run":
				if bu.Activated {
					script.Process(true)
				}
				break
			case "Testing Tree":
				if bu.Activated {
					script.Process(true)
					script.MakeTree()
				} else { // deactivated
					// remove all panels with trees
					b := gfx.Rend.Panels[:0]
					for _, pan := range gfx.Rend.Panels {
						if len(pan.Trees) < 1 {
							b = append(b, pan)
						}
					}
					gfx.Rend.Panels = b
					//fmt.Printf("len of b (from gfx.Rend.Panels) after removing ones with trees: %d\n", len(b))
					//fmt.Printf("len of gfx.Rend.Panels: %d\n", len(gfx.Rend.Panels))
				}
				break
			}

			gfx.Con.Add(fmt.Sprintf("%s toggled\n", bu.Name))
		}
	}
}
