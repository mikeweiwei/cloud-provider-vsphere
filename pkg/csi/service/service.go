/*
Copyright 2018 The Kubernetes Authors.

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

package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/rexray/gocsi"
	log "github.com/sirupsen/logrus"

	"k8s.io/cloud-provider-vsphere/pkg/csi/service/fcd"
)

const (
	// Name is the name of this CSI SP.
	Name = "io.k8s.cloud-provider-vsphere.vsphere"

	// Name of FCD API
	APIFCD = "FCD"

	defaultAPI = APIFCD
)

var (
	api = defaultAPI
)

// Service is a CSI SP and idempotency.Provider.
type Service interface {
	csi.IdentityServer
	csi.NodeServer
	GetController() csi.ControllerServer
	BeforeServe(context.Context, *gocsi.StoragePlugin, net.Listener) error
}

type service struct {
	cs csi.ControllerServer
}

// New returns a new Service.
func New() Service {
	// check which API to use

	api = os.Getenv(EnvAPI)
	if api == "" {
		api = defaultAPI
	}
	if strings.EqualFold(APIFCD, api) {
		return &service{
			cs: fcd.New(),
		}
	}
	return &service{}
}

func (s *service) GetController() csi.ControllerServer {
	return s.cs
}

func (s *service) BeforeServe(
	ctx context.Context, sp *gocsi.StoragePlugin, lis net.Listener) error {

	defer func() {
		fields := map[string]interface{}{
			"api": api,
		}

		log.WithFields(fields).Infof("configured: %s", Name)
	}()

	if s.cs == nil {
		return fmt.Errorf("Invalid API: %s", api)
	}

	return nil
}
