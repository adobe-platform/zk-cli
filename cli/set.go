package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"io"
	"io/ioutil"
)

type ZkSet struct {
	version int
	File string
	Data string
	byt  []byte

}

func (zk *ZkSet) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.IntVar(&zk.version, "version", -1, `data version.  default: -1`)
	flags.StringVar(&zk.File, "input-file", "", "data read from file for set node operation")
	flags.StringVar(&zk.Data, "input-data", "", "data for set node operation")
	return flags
}

func (zk *ZkSet) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk set", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer, "set")

	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
	Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug set --as-string

	Example - no auth
	./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo7 -debug set --as-string

	  `)

}
func (zk *ZkSet) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)

	flags := flag.NewFlagSet("zk set", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	if zk.File != "" {
		bb, err := ioutil.ReadFile(zk.File)
		if err != nil {
			return nil, err
		}
		zk.byt = bb
	} else if zk.Data != "" {
		zk.byt = []byte(zk.Data)
	}

	return zk, nil
}

func (zk *ZkSet) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkSet.Execute %+v", runtime)
	return runtime.Client.Set(runtime.ZkPath, zk.byt, int32(zk.version))
}

