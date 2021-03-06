/*
Copyright 2016 The Kubernetes Authors.

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

package registry

import (
	"fmt"
	"k8s.io/kops/pkg/client/simple"
	api "k8s.io/kops/pkg/apis/kops"
)

func CreateClusterConfig(clientset simple.Clientset, cluster *api.Cluster, groups []*api.InstanceGroup) error {
	// Check for instancegroup Name duplicates before writing
	{
		names := map[string]bool{}
		for i, ns := range groups {
			if ns.Name == "" {
				return fmt.Errorf("InstanceGroup #%d did not have a Name", i+1)
			}
			if names[ns.Name] {
				return fmt.Errorf("Duplicate InstanceGroup Name found: %q", ns.Name)
			}
			names[ns.Name] = true
		}
	}

	_, err := clientset.Clusters().Create(cluster)
	if err != nil {
		return err
	}

	for _, ig := range groups {
		_, err = clientset.InstanceGroups(cluster.Name).Create(ig)
		if err != nil {
			return fmt.Errorf("error writing updated instancegroup configuration: %v", err)
		}
	}

	return nil
}
