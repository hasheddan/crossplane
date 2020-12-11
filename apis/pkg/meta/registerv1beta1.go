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

package meta

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	V1beta1Group   = "meta.pkg.crossplane.io"
	V1beta1Version = "v1beta1"
)

var (
	// V1beta1SchemeGroupVersion is group version used to register these objects
	V1beta1SchemeGroupVersion = schema.GroupVersion{Group: V1beta1Group, Version: V1beta1Version}

	// V1beta1SchemeBuilder is used to add go types to the GroupVersionKind scheme
	V1beta1SchemeBuilder = &scheme.Builder{GroupVersion: V1beta1SchemeGroupVersion}

	// V1beta1AddToScheme adds all registered types to scheme
	V1beta1AddToScheme = V1beta1SchemeBuilder.AddToScheme
)

func init() {
	V1beta1SchemeBuilder.Register(&Configuration{})
	V1beta1SchemeBuilder.Register(&Provider{})
}
