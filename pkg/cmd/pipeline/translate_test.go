// Copyright 2021 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pipeline

import (
	"context"
	"testing"

	"github.com/okteto/okteto/pkg/model"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_translateConfigMap(t *testing.T) {
	ctx := context.Background()
	namespace := "test"
	cmap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TranslatePipelineName("test"),
			Namespace: namespace,
			Labels:    map[string]string{},
		},
		Data: map[string]string{
			statusField: DeployedStatus,
		},
	}
	fakeClient := fake.NewSimpleClientset(cmap)
	var tests = []struct {
		name    string
		status  string
		appName string
	}{
		{
			name:    "existing cmap",
			status:  DeployedStatus,
			appName: "test",
		},
		{
			name:    "existing cmap overwrite status",
			status:  ErrorStatus,
			appName: "test",
		},
		{
			name:    "not found cmap",
			status:  ProgressingStatus,
			appName: "not-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &CfgData{
				Name:      tt.appName,
				Namespace: namespace,
				Status:    tt.status,
			}
			cfg, err := TranslateConfigMapAndDeploy(ctx, data, fakeClient)

			assert.Nil(t, err)
			assert.Equal(t, cfg.Data[statusField], tt.status)
		})
	}
}

func Test_AddDevAnnotations(t *testing.T) {
	ctx := context.Background()
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "deployment",
			Namespace:   "namespace",
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
	}
	fakeClient := fake.NewSimpleClientset(d)
	t.Setenv(model.GithubRepositoryEnvVar, "git-repo")
	manifest := &model.Manifest{
		Namespace: "namespace",
		Dev: model.ManifestDevs{
			"deployment": &model.Dev{
				Name: "deployment",
			},
			"autocreate": &model.Dev{
				Autocreate: true,
			},
		},
	}
	AddDevAnnotations(ctx, manifest, fakeClient)
	d, _ = fakeClient.AppsV1().Deployments("namespace").Get(ctx, "deployment", metav1.GetOptions{})
	assert.Equal(t,
		d.Annotations,
		map[string]string{
			model.OktetoRepositoryAnnotation: "git-repo",
			model.OktetoDevNameAnnotation:    "deployment",
		},
	)
}
