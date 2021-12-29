/*
Copyright 2021 Takuya Okamoto.
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
