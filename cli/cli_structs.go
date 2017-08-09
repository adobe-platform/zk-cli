package cli

import (
	"fmt"
	"github.com/behance/go-logging/log"
)

type ArgList []string

// String - provide string helper
func (i *ArgList) String() string {
	return fmt.Sprintf("%s", *i)
}

// Set - Value interface
func (i *ArgList) Set(value string) error {
	log.Debugf("Args.Set %s", value)
	*i = append(*i, value)
	return nil
}
