package cli

import (
	"github.com/behance/go-logging/log"
	"flag"
	"io"
	"fmt"
)

type ZkGet struct {
	convertToString bool //AclList
	convertToHex    bool //AclList
}

func ( zk *ZkGet) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.BoolVar(&zk.convertToString, "as-string", false, `convert the node value to a string`)
	flags.BoolVar(&zk.convertToHex, "as-hex", false, `convert the node value to a hex`)
	return flags
}

func (zk *ZkGet) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk get", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer,"get")

	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer,`
	Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug get --as-string
	INFO[0000] result "something else"

	Example - no auth
	./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo7 -debug get --as-string
	INFO[0000] result "wide open"
	  `)

}
func (zk *ZkGet) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)

	flags := flag.NewFlagSet("zk get", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	return zk, nil
}

func (zk *ZkGet) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkGet.Execute %+v", runtime)
	data, stat, err := runtime.Client.Get(runtime.ZkPath)
	if err != nil {
		return nil, err
	}
	log.Debug(stat)
	if zk.convertToString {
		return string(data), nil
	} else if zk.convertToHex{
		return fmt.Sprintf("%X",data),nil
	}
	return data,nil
}


