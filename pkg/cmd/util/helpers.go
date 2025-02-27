/*
Copyright 2020 The Kruise Authors.

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

package util

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

const (
	DefaultErrorExitCode = 1
)

func Print(msg string) {
	if klog.V(2).Enabled() {
		klog.FatalDepth(2, msg)
	}
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
}

func fatal(msg string, code int) {
	if klog.V(2).Enabled() {
		klog.FatalDepth(2, msg)
	}
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

func CheckErr(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}
	fatal(msg, DefaultErrorExitCode)
}

func AddFieldManagerFlagVar(cmd *cobra.Command, p *string, defaultFieldManager string) {
	cmd.Flags().StringVar(p, "field-manager", defaultFieldManager, "Name of the manager used to track field ownership.")
}

func PatchSubResource(RESTClient resource.RESTClient, resource, subResource, namespace, name string, namespaceScoped bool, pt types.PatchType, data []byte, options *metav1.PatchOptions) (runtime.Object, error) {
	if options == nil {
		options = &metav1.PatchOptions{}
	}
	return RESTClient.Patch(pt).
		NamespaceIfScoped(namespace, namespaceScoped).
		Resource(resource).
		SubResource(subResource).
		Name(name).
		VersionedParams(options, metav1.ParameterCodec).
		Body(data).
		Do(context.TODO()).
		Get()
}

