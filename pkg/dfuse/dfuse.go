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
package dfuse

type Dfuse struct {
	config Config
}

type Config struct {
	DriverName string
	Version    string
	Endpoint   string
	NodeID     string
}

func (d *Driver) Run() {
	d.ns = NewNodeServer(d)
	s := NewNonBlockingGRPCServer()
	s.Start(d.endpoint,
		NewDefaultIdentityServer(d),
		nil, //NewControllerServiceCapability(d),
		d.ns,
		false,
	)
	s.Wait()
}

func NewNodeServer(d *Driver) *NodeServer {
	return &NodeServer{
		Driver: d,
		Mounts: make(map[string]NodeMount),
	}
}

func NewDfuseDriver(cfg Config) *Driver {
	//klog.V(2).Infof("Driver: %v version: %v", driverName, driverVersion)

	return &Driver{
		name:     cfg.DriverName,
		version:  cfg.Version,
		nodeID:   cfg.NodeID,
		endpoint: cfg.Endpoint,
		//cap:      map[csi.VolumeCapability_AccessMode_Mode]bool{},
		//perm:     perm,
	}
}
