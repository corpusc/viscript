package viewport

import (
	"github.com/skycoin/viscript/app"
	"github.com/skycoin/viscript/viewport/gl"
	"github.com/skycoin/viscript/viewport/terminal"
)

var (
	startMenuOpen bool
	startMenu     []MenuOption
	//current (iterators)
	buttonBounds *app.Rectangle
	charBounds   *app.Rectangle
)

type MenuOption struct {
	Name   string
	Bounds app.Rectangle
}

func drawTaskBar() {
	gl.SetColor(gl.Gray)
	drawTaskBarBackground()
	drawStartButton()
	drawTerminalButtons()
	drawStartMenu()
}

func drawTaskBarBackground() {
	buttonBounds = &app.Rectangle{
		-gl.CanvasExtents.Y + app.TaskBarHeight,
		gl.CanvasExtents.X,
		-gl.CanvasExtents.Y,
		-gl.CanvasExtents.X}

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		app.TaskBarDepth)
}

func drawStartButton() {
	if startMenuOpen {
		gl.SetColor(gl.White)
	}

	//now make buttons inset from task bar
	buttonBounds.Top -= app.TaskBarBorderSpan
	buttonBounds.Bottom += app.TaskBarBorderSpan
	buttonBounds.Left += app.TaskBarBorderSpan
	buttonBounds.Right = buttonBounds.Left + buttonBounds.Height()

	gl.Draw9SlicedRect(
		gl.Pic_GradientBorder,
		buttonBounds,
		app.TaskBarDepth)

	charBounds = getCharBoundsInsetFrom(buttonBounds)
	buttonBounds.Left += buttonBounds.Width()
	buttonBounds.Right += buttonBounds.Width()

	gl.Draw9SlicedRect(
		gl.Pic_TriangleUp,
		charBounds,
		app.TaskBarDepth)

	charBounds.Left += app.TaskBarCharWid
	charBounds.Right += app.TaskBarCharWid
}

func drawStartMenu() {
	if startMenuOpen {
		gl.SetColor(gl.Gray)

		if len(startMenu) < 1 {
			maxOptionWid := gl.CanvasExtents.X //FIX this temp insanity
			y := -gl.CanvasExtents.Y + app.TaskBarHeight
			rect := app.Rectangle{
				y + app.TaskBarHeight,
				-gl.CanvasExtents.X + maxOptionWid,
				y,
				-gl.CanvasExtents.X}
			startMenu = append(startMenu, MenuOption{"New Terminal", rect})
			rect.Top += app.TaskBarHeight
			rect.Bottom += app.TaskBarHeight
			startMenu = append(startMenu, MenuOption{"Unimplemented", rect})
		}

		for _, option := range startMenu {
			//draw option background
			gl.Draw9SlicedRect(
				gl.Pic_GradientBorder,
				&option.Bounds,
				app.TaskBarDepth)

			//draw text
			charBounds = getCharBoundsInsetFrom(&option.Bounds)
			charBounds.Right = charBounds.Left + app.TaskBarCharWid

			for i := 0; i < len(option.Name); i++ {
				gl.DrawCharAtRect(rune(option.Name[i]), charBounds, app.TaskBarDepth)
				charBounds.Left += app.TaskBarCharWid
				charBounds.Right += app.TaskBarCharWid
			}
		}
	}

}

func drawTerminalButtons() {
	for _, term := range terminal.Terms.TermMap {
		charBounds.Left = term.TaskBarButton.Left + app.TaskBarBorderSpan
		charBounds.Right = charBounds.Left + app.TaskBarCharWid

		if term.TerminalId == terminal.Terms.FocusedId {
			gl.SetColor(gl.White)
		} else {
			gl.SetColor(gl.Gray)
		}

		//draw button background
		gl.Draw9SlicedRect(
			gl.Pic_GradientBorder,
			term.TaskBarButton,
			app.TaskBarDepth)

		//prepare for id text
		textMax := term.TaskBarButton.Right - app.TaskBarBorderSpan
		//when abbreviating text, append "..." chars...
		dotWid := app.TaskBarCharWid / 2 //...but at half width

		if term.TaskBarButton.Width()-app.TaskBarBorderSpan*2 <
			float32(len(term.TaskBarButtonText))*app.TaskBarCharWid {

			textMax -= (3 * dotWid)
		}

		//draw id text
		max := len(term.TaskBarButtonText)
		for i := 0; i < max; i++ {
			if charBounds.Right <= textMax {
				gl.DrawCharAtRect(rune(term.TaskBarButtonText[i]), charBounds, app.TaskBarDepth)
			} else { //draw 3 dots
				charBounds.Right = charBounds.Left + dotWid

				for i := 0; i < 3; i++ {
					gl.DrawCharAtRect('.', charBounds, app.TaskBarDepth)
					charBounds.Left += dotWid
					charBounds.Right += dotWid
				}

				break
			}

			charBounds.Left += app.TaskBarCharWid
			charBounds.Right += app.TaskBarCharWid
		}
	}
}

func getCharBoundsInsetFrom(bounds *app.Rectangle) *app.Rectangle {
	return &app.Rectangle{
		bounds.Top - app.TaskBarBorderSpan,
		bounds.Right - app.TaskBarBorderSpan,
		bounds.Bottom + app.TaskBarBorderSpan,
		bounds.Left + app.TaskBarBorderSpan}
}
