package main

import (
	"flag"
	"fmt"
	//"os"
	//"github.com/golang/glog"
	"github.com/tkokamo/csi-driver-daos/pkg/dfuse"
)

func main() {
	cfg := dfuse.Config{}
	flag.StringVar(&cfg.DriverName, "drivername", "dfuse.csi.k8s.io", "name of the driver")
	flag.StringVar(&cfg.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.Version, "version", "0.0.1", "version of the driver")
	flag.StringVar(&cfg.NodeID, "nodeid", "", "node ID")
	flag.Parse()
	fmt.Printf("%v", cfg.Endpoint)

	driver := dfuse.NewDfuseDriver(cfg)
	driver.Run()
}
