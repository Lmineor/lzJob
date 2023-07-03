package cmd

import "github.com/spf13/pflag"

type DbFlags struct {
	DbUrl        string
	DbPort       int32
	DbConfig     string
	Db           string
	DbUserName   string
	DbPassword   string
	MaxIdleConns int
	MaxOpenConns int
}
func newDbFlags()*DbFlags {
	return &DbFlags{
		DbUrl:        "127.0.0.1",
		DbPort:       3306,
		DbConfig:     "charset=utf8&parseTime=True&loc=Local",
		Db:           "lzJob",
		DbUserName:   "lzJob",
		DbPassword:   "123456",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
}
func (f *DbFlags) AddFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&f.DbUrl, "connection", "C", "127.0.0.1", "address of db")
	fl.Int32VarP(&f.DbPort, "db_port", "P", 3306, "port of db")
	fl.StringVarP(&f.DbConfig, "db_config", "", "charset=utf8&parseTime=True&loc=Local", "db config")
	fl.StringVarP(&f.Db, "db_name", "d", "lzJob", "db name")
	fl.StringVarP(&f.DbUserName, "db_username", "u", "root", "db username")
	fl.StringVarP(&f.DbPassword, "db_password", "e", "123456", "db password")
	fl.IntVarP(&f.MaxIdleConns, "db_max_idle_conns", "i", 10, "db_max_idle_conns")
	fl.IntVarP(&f.MaxOpenConns, "db_max_open_conns", "o", 100, "db_max_open_conns")

}
