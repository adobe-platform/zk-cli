package main

import (
	"fmt"
	log "github.com/behance/go-logrus"
	"os"
	"strings"
	//"errors"
	"encoding/json"
	cli "github.com/adobe-platform/zk-cli/cli"
)

type commandMap map[string]cli.CommandParse

var commands commandMap

// initialize the top level command map
func init() {
	commands = commandMap{
		"create": cli.CommandParse(new(cli.CreateZkNode)),
		"rm":     cli.CommandParse(new(cli.ZkRm)),
		"ls":     cli.CommandParse(new(cli.ZkLs)),
		"get":    cli.CommandParse(new(cli.ZkGet)),
		"getAcl": cli.CommandParse(new(cli.ZkGetAcl)),
		"setAcl": cli.CommandParse(new(cli.ZkSetAcl)),
	}
}

func usage(msg string) {
	if msg != "" {
		log.Errorf(" %s ", msg)
	}
	options := []string{
		"create",
		"rm",
		"get",
		"getAcl",
		"setAcl",
	}
	fmt.Fprintf(os.Stderr, `USAGE

	 %s <global-options>  {%s|help} [<action options>|help]

COMMANDS:
	 `, os.Args[0], strings.Join(options, "|"))
	fmt.Fprintln(os.Stderr, "")
	for _, action := range options {
		commands[action].Usage(os.Stderr)
	}
	fmt.Fprintln(os.Stderr, `
GLOBAL OPTIONS:
 The zk target path is taken --zk-hosts for most commands
		`)

	runtime := &cli.Runtime{}
	runtime.Usage(os.Stderr)
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}
func In(targ string, list []string) (retval bool) {
	for _, cc := range list {
		if cc == targ {
			return true
		}
	}
	return retval
}
func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) == 1 {
		usage("")
	}
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	index := -1
	var action string
	for v, value := range os.Args {
		if In(value, keys) {
			index = v
			action = value
			break
		}
	}
	runtime := &cli.Runtime{}
	if index != -1 {
		var commonArgs []string
		if index > 1 {
			commonArgs = os.Args[1:index]
		} else {
			log.Debugf("No command args used\n")
		}
		log.Debugf("commonArgs %+v action: %s\n", commonArgs, action)
		if _, err := runtime.Parse(commonArgs); err != nil {
			usage(err.Error())
		} else if action == "" {
			usage("missing action")
		} else if commands[action] == nil {
			usage(fmt.Sprintf("'%s' command not defined", action))
		}
		if runtime.Debug {
			log.SetLevel(log.DebugLevel)
		}
		log.Debugf("runtime: %#v", runtime)
		var executorArgs []string
		if len(os.Args) > (index + 1) {
			executorArgs = os.Args[index+1:]
		}
		log.Debugf("executorArgs %+v", executorArgs)
		if action == "help" {
			usage("your help:")
		} else if executor, err := commands[action].Parse(executorArgs); err != nil {
			log.Fatalf("%s failed because %+v", action, err)
		} else {
			if result, err2 := executor.Execute(runtime); err2 != nil {
				log.Fatalf("action %s execution failed because %+v", action, err2)
			} else {
				log.Debugf("Result type: %T", result)

				switch result.(type) {
				case json.RawMessage:
					var f interface{}
					by := result.(json.RawMessage)
					if err := json.Unmarshal(by, &f); err != nil {
						log.Infof(string(by))
					} else {
						if b2, err2 := json.MarshalIndent(f, "", "  "); err2 != nil {
							log.Infof(string(by))
						} else {
							log.Infof(string(b2))
						}
					}
				default:
					if bb, err7 := json.MarshalIndent(result, "", "  "); err7 == nil {
						log.Infof("result %s\n", (string(bb)))
					}
				}
			}
		}
	} else {
		if len(os.Args) > 1 {
			usage("You need to include a verb")
		} else {
			usage("Nothing to do.  You need to choose an action\n")
		}
	}
}
