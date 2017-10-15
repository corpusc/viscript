package terminal

import (
	"github.com/skycoin/viscript/msg"
	"strconv"
)

func (ts *TerminalStack) OnUserCommandFinalStage(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	switch cmd.Command {

	case "close_term":
		fallthrough
	case "focus":
		ts.onGivenTerminalId(cmd)
	case "defocus":
		ts.Defocus()
	case "list_terms":
		ts.commandListTerminals(receiver, cmd)
	case "new_term":
		//temporary
		//for now, we'll be testing the difference between fixed size and dynamic terminals.
		//the 1st/initial terminal will be dynamic.  new terms afterwards will all be fixed.
		ts.AddWithFixedSizeState(true)
	default:
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)
		println("OnUserCommandFinalStage()   UNHANDLED COMMAND!!!:", cmd.Command)

	}
}

//
//
//private
//
//

func (ts *TerminalStack) onGivenTerminalId(cmd msg.MessageTokenizedCommand) {
	matchedTerm := msg.TerminalId(0)
	for _, t := range ts.TermMap {
		ruledOutMatch := false
		arg := cmd.Args[0]
		tId := strconv.Itoa(int(t.TerminalId))

		//chop runes off the end
		//(if user gave more digits than an id has)
		for len(arg) > len(tId) {
			arg = arg[:len(arg)-1]
		}

		//compare each rune of user input to leftmost runes of id
		for i, c := range arg {
			if c != rune(tId[i]) {
				ruledOutMatch = true
				break
			}
		}

		if !ruledOutMatch {
			matchedTerm = t.TerminalId
			break
		}
	}

	//set new focus (or show error)
	if matchedTerm != 0 {
		println("finalStage   -   Found terminal starting with:", cmd.Args[0])

		switch cmd.Command {

		case "close_term":
			if len(ts.TermMap) < 2 {
				//TODO: print feedback in Terminal
				//TODO: print feedback in Terminal
				//TODO: print feedback in Terminal
				s := "Shouldn't close when only 1 terminal remains (UNTIL GUI IS MADE)"
				println(s)
				//st.PrintError("Shouldn't close when only 1 terminal remains (UNTIL GUI IS MADE)")
				ts.TermMap[matchedTerm].PutString(s)
				ts.TermMap[matchedTerm].NewLine()
				return
			}

			ts.Remove(matchedTerm)

		case "focus":
			ts.SetFocused(matchedTerm)

		}
	} else {
		//TODO: print feedback in Terminal
		//TODO: print feedback in Terminal
		//TODO: print feedback in Terminal
		println("ERROR!!!  \"" + cmd.Args[0] + "\" is not the beginning of any Terminal id.")
	}
}

func (ts *TerminalStack) commandListTerminals(receiver msg.TerminalId, cmd msg.MessageTokenizedCommand) {
	var m msg.MessageTerminalIds
	m.Focused = receiver

	for _, term := range ts.TermMap {
		m.TermIds = append(m.TermIds, term.TerminalId)
	}

	ts.GetFocusedTerminal().RelayToTask(msg.Serialize(msg.TypeTerminalIds, m))
}
