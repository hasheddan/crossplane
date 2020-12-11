/*
Copyright 2020 The Crossplane Authors.

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

package v1beta1

import (
	"github.com/crossplane/crossplane/apis/pkg/meta"
)

// ConvertV1beta1ConfigurationToInternal converts a v1beta1 Configuration to
// its internal representation.
func ConvertV1beta1ConfigurationToInternal(in *Configuration, out *meta.Configuration) {
	out.ObjectMeta = in.ObjectMeta
	ConvertV1beta1CrossplaneConstraintsToInternal(in.Spec.Crossplane, out.Spec.Crossplane)
	ConvertV1beta1DependenciesToInternal(in.Spec.DependsOn, out.Spec.DependsOn)
}

// ConvertV1beta1ProviderToInternal converts a v1beta1 Provider to its
// internal representation.
func ConvertV1beta1ProviderToInternal(in *Provider, out *meta.Provider) {
	out.ObjectMeta = in.ObjectMeta
	out.Spec.Controller = meta.ControllerSpec{
		Image:              in.Spec.Controller.Image,
		PermissionRequests: in.Spec.Controller.PermissionRequests,
	}
	ConvertV1beta1CrossplaneConstraintsToInternal(in.Spec.Crossplane, out.Spec.Crossplane)
	ConvertV1beta1DependenciesToInternal(in.Spec.DependsOn, out.Spec.DependsOn)
}

// ConvertV1beta1CrossplaneConstraintsToInternal converts v1beta1
// CrossplaneConstraints to its internal representation.
func ConvertV1beta1CrossplaneConstraintsToInternal(in *CrossplaneConstraints, out *meta.CrossplaneConstraints) {
	if in != nil {
		out = &meta.CrossplaneConstraints{
			Version: in.Version,
		}
	}
}

// ConvertV1beta1DependenciesToInternal converts v1beta1 Dependencies to its
// internal representation.
func ConvertV1beta1DependenciesToInternal(in []Dependency, out []meta.Dependency) {
	out = make([]meta.Dependency, len(in))
	for i, d := range in {
		out[i] = meta.Dependency{
			Provider:      d.Provider,
			Configuration: d.Configuration,
			Version:       d.Version,
		}
	}
}
