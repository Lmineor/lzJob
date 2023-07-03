package cmd

import "github.com/spf13/pflag"

type Flags struct {
	ConfigFile string
	*SysFlags
	*DbFlags
}



func NewFlags()*Flags {
	sysF := newSysFlags()
	dbF := newDbFlags()
	return &Flags{
		ConfigFile: "/etc/lzJob.yaml",
		SysFlags: sysF,
		DbFlags: dbF,
	}
}
func (f *Flags) AddFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&f.ConfigFile, "config-file", "c", "", "config file")
	f.SysFlags.AddFlags(fl)
	f.DbFlags.AddFlags(fl)
}
