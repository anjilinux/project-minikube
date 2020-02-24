/*
Copyright 2019 The Kubernetes Authors All rights reserved.

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

package node

import (
	"errors"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/machine"
)

const (
	cacheImages         = "cache-images"
	waitUntilHealthy    = "wait"
	cacheImageConfigKey = "cache"
	containerRuntime    = "container-runtime"
	embedCerts          = "embed-certs"
	keepContext         = "keep-context"
	mountString         = "mount-string"
	createMount         = "mount"
	waitTimeout         = "wait-timeout"
)

// Add adds a new node config to an existing cluster.
func Add(cc *config.ClusterConfig, n config.Node) error {
	cc.Nodes = append(cc.Nodes, n)
	err := config.SaveProfile(cc.Name, cc)
	if err != nil {
		return err
	}

	err = Start(*cc, n, nil)

	return err
}

// Delete stops and deletes the given node from the given cluster
func Delete(cc config.ClusterConfig, name string) error {
	_, index, err := Retrieve(&cc, name)
	if err != nil {
		return err
	}

	if err != nil {
		glog.Warningf("Failed to stop node %s. Will still try to delete.", name)
	}

	api, err := machine.NewAPIClient()
	if err != nil {
		return err
	}

	err = machine.DeleteHost(api, name)
	if err != nil {
		return err
	}

	cc.Nodes = append(cc.Nodes[:index], cc.Nodes[index+1:]...)
	return config.SaveProfile(viper.GetString(config.MachineProfile), &cc)
}

// Retrieve finds the node by name in the given cluster
func Retrieve(cc *config.ClusterConfig, name string) (*config.Node, int, error) {
	for i, n := range cc.Nodes {
		if n.Name == name {
			return &n, i, nil
		}
	}

	return nil, -1, errors.New("Could not find node " + name)
}
