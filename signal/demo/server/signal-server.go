package main

import (
	"github.com/skycoin/viscript/signal"
	"github.com/skycoin/viscript/msg"
	"strings"
	"bufio"
	"os"
	"log"
	"fmt"
	"strconv"
)

func main() {
	server := signal.Init("127.0.0.1:7999")
	server.Run()
	showHelp()
	promptCycle(server)
}

func promptCycle(server *signal.MonitorServer) {
	for {
		newCommand, args := inputFromCli()
		if newCommand == "" {
			continue
		}
		dispatcher(server, strings.ToLower(newCommand), args)
	}
}

func inputFromCli() (command string, args []string) {
	command = ""
	args = []string{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	splitInput := strings.Fields(input)
	if len(splitInput) == 0 {
		return
	}

	command = strings.Trim(splitInput[0], " ")
	if len(splitInput) > 1 {
		args = splitInput[1:]
	}
	return
}

func dispatcher(server *signal.MonitorServer, cmd string, args []string) {
	log.Println("new command:" + cmd)
	log.Println(args[0:])



	switch cmd {

	case "help":
		showHelp()

	case "ping":
		s, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println(err)
		}
		appId := uint32(s)
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    appId,
			Payload:  msg.Serialize(msg.TypePing, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(appId, serializedCommand)

	case "shutdown":
		s, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println(err)
		}
		appId := uint32(s)
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    appId,
			Payload:  msg.Serialize(msg.TypeShutdown, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(appId, serializedCommand)

	case "res_usage":
		s, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			log.Println(err)
		}
		appId := uint32(s)
		msgUserCommand := msg.MessageUserCommand{
			Sequence: 1,
			AppId:    appId,
			Payload:  msg.Serialize(msg.TypeResourceUsage, msg.MessagePing{})}

		serializedCommand := msg.Serialize(msg.TypeUserCommand, msgUserCommand)

		signal.Monitor.Send(appId, serializedCommand)

	case "add_node":
		server.AddSignalNodeConn(args[0], args[1])

	case "list_nodes":
		server.ListNodes()


	default:
		log.Println("Unknown user command:")

	}

}

func showHelp() {
	fmt.Printf("> help\t\t\tShow list of commands.\n\n")
	fmt.Printf("> ping <id>\t\tPing app with choosen id.\n\n")
	fmt.Printf("> shutdown <id>\t\tKill app with choosen id.\n\n")
	fmt.Printf("> res_usage <id>\tShow cpu and memory stats.\n\n")
	fmt.Printf("> add_node <ip> <port>\tShow cpu and memory stats.\n\n")
	fmt.Printf("> list_nodes\t\tShow list of runnig apps.\n\n")
}
