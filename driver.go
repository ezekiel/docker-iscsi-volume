package main

import (
	"docker-iscsi-volume/iscsi"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
)

type iscsiDriver struct {
	baseMountPath string
	m             *sync.Mutex
}

func ISCSIVolumeDriver(root string) iscsiDriver {
	return iscsiDriver{
		baseMountPath: root,
		m:             &sync.Mutex{},
	}
}

func (driver iscsiDriver) Create(request *volume.CreateRequest) error {
	return nil
}

func (d iscsiDriver) mountpoint(name string) string {
	return filepath.Join(d.baseMountPath, name)
}

func (driver iscsiDriver) List() (*volume.ListResponse, error) {
	return &volume.ListResponse{}, nil
}

func (driver iscsiDriver) Get(request *volume.GetRequest) (*volume.GetResponse, error) {
	return &volume.GetResponse{}, nil
}

func (driver iscsiDriver) Remove(request *volume.RemoveRequest) error {
	// logout all mountpoints
	return nil
}

func (driver iscsiDriver) Path(request *volume.PathRequest) (*volume.PathResponse, error) {
	return &volume.PathResponse{Mountpoint: driver.mountpoint(request.Name)}, nil
}

func (driver iscsiDriver) Mount(request *volume.MountRequest) (*volume.MountResponse, error) {

	m := driver.mountpoint(request.Name)

	log.Printf("Mounting volume %s on %s\n", request.Name, m)

	//Create a temp folder in mountPath.
	os.Mkdir(m, os.ModeDir)
	plugin := iscsi.NewISCSIPlugin()
	err := plugin.LoginTarget("10.0.2.15:3260", "cc")
	if err == nil {
		return &volume.MountResponse{Mountpoint: m}, nil
	}
	//Login te target
	//Mount logic.

	return &volume.MountResponse{}, fmt.Errorf("no such volume")
}

func (driver iscsiDriver) Unmount(request *volume.UnmountRequest) error {

	m := driver.mountpoint(request.Name)
	log.Printf("Unmount volume %s on %s\n", request.Name, m)

	//Create a temp folder in mountPath.
	//Login te target
	//Mount logic.

	return nil
}

func (driver iscsiDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}
