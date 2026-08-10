/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and valkey-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package operator

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sap/component-operator-runtime/pkg/manifests"

	operatorv1alpha1 "github.com/sap/valkey-operator/api/v1alpha1"
)

// The chart refuses images which are not the ones it was tested with, unless allowInsecureImages is set;
// and pull secrets are passed as a global value, so that they apply to all images.
func TestImageParameters(t *testing.T) {
	g := NewWithT(t)

	transformer, err := manifests.NewTemplateParameterTransformer(data, "data/parameters.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	parameters, err := transformer.TransformParameters("testns", "test", &operatorv1alpha1.ValkeySpec{
		Replicas: 1,
		Image: &operatorv1alpha1.ImageProperties{
			Registry:    "registry.example.com",
			Repository:  "mirror/valkey",
			PullSecrets: []string{"my-pull-secret"},
		},
	})
	g.Expect(err).NotTo(HaveOccurred())
	values := parameters.ToUnstructured()

	g.Expect(values["image"]).To(HaveKeyWithValue("repository", "mirror/valkey"))
	g.Expect(values["global"]).To(HaveKeyWithValue("imageRegistry", "registry.example.com"))
	g.Expect(values["global"]).To(HaveKeyWithValue("imagePullSecrets", []any{"my-pull-secret"}))
	g.Expect(values["global"]).To(HaveKeyWithValue("security", map[string]any{"allowInsecureImages": true}))

	parameters, err = transformer.TransformParameters("testns", "test", &operatorv1alpha1.ValkeySpec{Replicas: 1})
	g.Expect(err).NotTo(HaveOccurred())
	values = parameters.ToUnstructured()

	g.Expect(values).NotTo(HaveKey("global"))
}
