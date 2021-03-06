package msg

import (
	"math/rand"
)

var TaskIdGlobal TaskId = 1 //sequential
var ExternalAppIdGlobal ExternalAppId = 1

type TaskId uint64
type ExternalAppId uint64
type TerminalId uint64

//
//
//
func NextTaskId() TaskId {
	TaskIdGlobal += 1
	return TaskIdGlobal
}

func NextExternalAppId() ExternalAppId {
	ExternalAppIdGlobal += 1
	return ExternalAppIdGlobal
}

func RandTerminalId() TerminalId {
	return (TerminalId)(rand.Int63())
}
