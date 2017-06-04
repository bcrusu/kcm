package main

import (
	"flag"

	"github.com/bcrusu/kcm/cmd"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	if err := cmd.RootCmd.Execute(); err != nil {
		glog.V(2).Info(err)
	}
}
