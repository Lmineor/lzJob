package cmd

import "github.com/spf13/pflag"

type SysFlags struct {
	Addr string
	Port int32
}

func newSysFlags()*SysFlags {
	return &SysFlags{
		Addr: "127.0.0.1",
		Port: 9191,
	}
}
func (f *SysFlags) AddFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&f.Addr, "address", "a", "127.0.0.1", "address of api")
	fl.Int32VarP(&f.Port, "port", "p", 9191, "port of api")
}
