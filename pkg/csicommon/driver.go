/*
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

package csicommon

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/opencurve/curvefs-csi/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

type CSIDriver struct {
	Name    string
	NodeID  string
	Version util.VersionInfo
	cap     []*csi.ControllerServiceCapability
	vc      []*csi.VolumeCapability_AccessMode
}

// Creates a NewCSIDriver object. Assumes vendor version is equal to driver version &
// does not support optional driver plugin info manifest field. Refer to CSI spec for more details.
func NewCSIDriver(name string, version util.VersionInfo, nodeID string) *CSIDriver {
	klog.Infof("Driver: %v version info %v nodeId", name, version, nodeID)
	if name == "" {
		klog.Errorf("Driver name missing")
		return nil
	}

	if nodeID == "" {
		klog.Errorf("NodeID missing")
		return nil
	}

	if version == (util.VersionInfo{}) {
		klog.Errorf("Version argument missing")
		return nil
	}

	driver := CSIDriver{
		Name:    name,
		Version: version,
		NodeID:  nodeID,
	}

	return &driver
}

func (d *CSIDriver) ValidateControllerServiceRequest(
	c csi.ControllerServiceCapability_RPC_Type,
) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range d.cap {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s", c))
}

func (d *CSIDriver) AddControllerServiceCapabilities(
	cl []csi.ControllerServiceCapability_RPC_Type,
) {
	var csc []*csi.ControllerServiceCapability

	for _, c := range cl {
		klog.Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}

	d.cap = csc

	return
}

func (d *CSIDriver) AddVolumeCapabilityAccessModes(
	vc []csi.VolumeCapability_AccessMode_Mode,
) []*csi.VolumeCapability_AccessMode {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		klog.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.vc = vca
	return vca
}

func (d *CSIDriver) GetVolumeCapabilityAccessModes() []*csi.VolumeCapability_AccessMode {
	return d.vc
}
