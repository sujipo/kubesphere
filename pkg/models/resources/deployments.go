/*

 Copyright 2019 The KubeSphere Authors.

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
package resources

import (
	"kubesphere.io/kubesphere/pkg/informers"
	"kubesphere.io/kubesphere/pkg/server/params"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/api/apps/v1"
)

type deploymentSearcher struct {
}

func (*deploymentSearcher) get(namespace, name string) (interface{}, error) {
	return informers.SharedInformerFactory().Apps().V1().Deployments().Lister().Deployments(namespace).Get(name)
}

func deploymentStatus(item *v1.Deployment) string {
	if item.Spec.Replicas != nil {
		if item.Status.ReadyReplicas == 0 && *item.Spec.Replicas == 0 {
			return StatusStopped
		} else if item.Status.ReadyReplicas == *item.Spec.Replicas {
			return StatusRunning
		} else {
			return StatusUpdating
		}
	}
	return StatusStopped
}

// Exactly Match
func (*deploymentSearcher) match(kv map[string]string, item *v1.Deployment) bool {
	for k, v := range kv {
		switch k {
		case Status:
			if deploymentStatus(item) != v {
				return false
			}
		default:
			if !match(k, v, item.ObjectMeta) {
				return false
			}
		}
	}
	return true
}

func (*deploymentSearcher) fuzzy(kv map[string]string, item *v1.Deployment) bool {

	for k, v := range kv {
		if !fuzzy(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (s *deploymentSearcher) compare(a, b *v1.Deployment, orderBy string) bool {
	switch orderBy {
	case UpdateTime:
		aLastUpdateTime := s.lastUpdateTime(a)
		bLastUpdateTime := s.lastUpdateTime(b)
		if aLastUpdateTime.Equal(bLastUpdateTime) {
			return strings.Compare(a.Name, b.Name) <= 0
		}
		return aLastUpdateTime.Before(bLastUpdateTime)
	default:
		return compare(a.ObjectMeta, b.ObjectMeta, orderBy)
	}
}

func (s *deploymentSearcher) lastUpdateTime(deployment *v1.Deployment) time.Time {
	lastUpdateTime := deployment.CreationTimestamp.Time
	for _, condition := range deployment.Status.Conditions {
		if condition.LastUpdateTime.After(lastUpdateTime) {
			lastUpdateTime = condition.LastUpdateTime.Time
		}
	}
	return lastUpdateTime
}

func (s *deploymentSearcher) search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	deployments, err := informers.SharedInformerFactory().Apps().V1().Deployments().Lister().Deployments(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*v1.Deployment, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = deployments
	} else {
		for _, item := range deployments {
			if s.match(conditions.Match, item) && s.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			tmp := i
			i = j
			j = tmp
		}
		return s.compare(result[i], result[j], orderBy)
	})

	r := make([]interface{}, 0)
	for _, i := range result {
		r = append(r, i)
	}
	return r, nil
}
