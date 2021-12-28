package dfuse

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type NodeServer struct {
	Driver *Driver
	Mounts []*NodeMount
}

type NodeMount struct {
	PoolID    string
	ContID    string
	MountPath string
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	uid, err := strconv.Atoi(req.GetVolumeContext()["uid"])
	size := req.GetVolumeContext()["size"]
	targetPath := req.GetTargetPath()

	err = syscall.Setuid(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cmd := fmt.Sprintf("dmg pool create -i --ranks=0 --scm-size=%s | grep UUID | awk '{print $3}'", size)
	out, err := exec.Command("sh", "-c", cmd).Output()
	poolid := strings.TrimRight(string(out), "\n")
	// TODO create container
	// TODO mount with dfuse

	// add mount record for unpublish
	// append(ns.Mounts, &NodeMount {
	// PoolID: poolid,
	// ContID: contid,
	// MountPath: targetPath,
	// })
	fmt.Printf("%s %s", targetPath, poolid)
	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmount the volume
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetInfo return info of the node on which this plugin is running
func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.Driver.nodeID,
	}, nil
}

// NodeGetCapabilities return the capabilities of the Node plugin
// XXX Must be implemented
func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	//return &csi.NodeGetCapabilitiesResponse{
	//	Capabilities: ns.Driver.nscap,
	//}, nil
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
					},
				},
			},
		},
	}, nil
}

func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
