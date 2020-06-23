/*
Copyright AppsCode Inc.

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

package lib

import (
	"fmt"
	"testing"
)

// ref: https://gist.github.com/inotnako/c4a82f6723f6ccea5d83c5d3689373dd
// ref: https://github.com/keighl/metabolize
// ref: https://github.com/rsc/go-import-redirector/blob/master/main.go#L134
//ref: https://github.com/appscodelabs/gh-release-automation-testing/issues/22
func TestDetectVCSRoot(t *testing.T) {
	// res, _ := http.Get("https://stash.appscode.dev/cli?go-get=1")
	// res, _ := http.Get("https://k8s.io/api?go-get=1")
	// "https://github.com/cloudevents/sdk-go/blob/master/samples/kafka?go-get=1"

	r, err := DetectVCSRoot("stash.appscode.dev/cli")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
