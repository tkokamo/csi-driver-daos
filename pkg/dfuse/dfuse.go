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
