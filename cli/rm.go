package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"io"
)

type ZkRm struct {
	recursive bool //AclList
}

func (zk *ZkRm) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.BoolVar(&zk.recursive, "recursive", false, `convert the node value to a string`)
	return flags
}

func (zk *ZkRm) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk rm", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer, "rm")

	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
	Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug rm  --recursive
	  `)

}
func (zk *ZkRm) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)

	flags := flag.NewFlagSet("zk rm", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	return zk, nil
}

func (zz *ZkRm) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkRm.Execute %+v", runtime)
	children, err := func() ([]string, error) {
		if zz.recursive {
			return GetChildren(runtime.ZkPath, runtime.Client)
		}
		return []string{runtime.ZkPath}, nil
	}()
	if err != nil {
		return nil, err
	}
	for ii := len(children) - 1; ii >= 0; ii-- {
		leaf := children[ii]
		children = children[:ii]
		if err = Delete(leaf, runtime.Client); err != nil {
			log.Errorf("Error deleting %s", leaf)
		}
		log.Debugf("deleted %s", leaf)
	}
	return GenericResult{Success: true, Message: fmt.Sprintf("%s deleted successfully", runtime.ZkPath)}, nil
}
