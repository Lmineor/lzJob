package main

import (
	goflag "flag"
	"fmt"
	"github.com/Lmineor/lzJob/cmd"
	"k8s.io/klog/v2"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	klog.InitFlags(nil)
	// By default klog writes to stderr. Setting logtostderr to false makes klog
	goflag.Set("logtostderr", "false")
	command := cmd.RootCmd()
	defer klog.Flush()
	if err := command.Execute();err != nil{
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}