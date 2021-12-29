package dfuse

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type NodeServer struct {
	Driver *Driver
	Mounts map[string]NodeMount
}

type NodeMount struct {
	PoolID string
	ContID string
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.Infof("NodePublishVolume: %v", req)
	uid, err := strconv.Atoi(req.GetVolumeContext()["uid"])
	poolid := req.GetVolumeContext()["poolid"]
	targetPath := req.GetTargetPath()
	os.MkdirAll(targetPath, os.ModePerm)

	klog.Infof("uid: %d, poolid: %s, targetPath: %s", uid, poolid, targetPath)

	// Switch user to create container and mount FS successfully
	err = syscall.Setuid(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Create container
	cmd := fmt.Sprintf("/opt/daos/bin/daos cont create --pool=%s --type=POSIX | awk '{print $4}'", poolid)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		klog.Infof("create container: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	contid := strings.TrimRight(string(out), "\n")
	klog.Infof("cmd: %s", cmd)
	klog.Infof("contid: %s", contid)

	// Mount with dfuse
	cmd = fmt.Sprintf("sudo sudo -u#%d /opt/daos/bin/dfuse --mountpoint %s --pool=%s --container=%s", uid, targetPath, poolid, contid)
	out, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		klog.Infof("mount: %v, cmd: %v", err, cmd)
		klog.Infof("output: %s", out)
		cmd = fmt.Sprintf("/opt/daos/bin/daos cont destroy --pool=%s --cont=%s", poolid, contid)
		exec.Command("sh", "-c", cmd).Output()
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("output: %s", out)

	nm := NodeMount{
		PoolID: poolid,
		ContID: contid,
	}
	ns.Mounts[targetPath] = nm

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmount the volume
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	cmd := fmt.Sprintf("fusermount3 -u %s", targetPath)
	//out, err := exec.Command("sh", "-c", cmd).Output()
	exec.Command("sh", "-c", cmd).Output()

	nm := ns.Mounts[targetPath]
	contid := nm.ContID
	poolid := nm.PoolID

	cmd = fmt.Sprintf("/opt/daos/bin/daos cont destroy --pool=%s --cont=%s", poolid, contid)
	//out, err = exec.Command("sh", "-c", cmd).Output()
	exec.Command("sh", "-c", cmd).Output()

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
