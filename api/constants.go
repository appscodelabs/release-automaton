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

package api

const (
	Workspace      = "/tmp/workspace"
	GitHubUserKey  = "GITHUB_USER"
	GitHubTokenKey = "GITHUB_TOKEN"
	BranchMaster   = "master"
	LabelLocked    = "locked"
	LabelAutoMerge = "automerge"
	ReleasesDir    = "releases"

	StableChartRegistry     = "appscode"
	StableChartRegistryURL  = "https://charts.appscode.com/stable/"
	TestChartRegistry       = "appscode-testing"
	TestChartRegistryURL    = "https://charts.appscode.com/testing/"
	StableUIRegistry        = "bytebuilders-ui"
	StableUIRegistryURL     = "https://bundles.byte.builders/ui/"
	TestUIRegistry          = "bytebuilders-ui-dev"
	TestUIRegistryURL       = "https://raw.githubusercontent.com/bytebuilders/ui-wizards/"
	StableBundleRegistry    = "bytebuilders"
	StableBundleRegistryURL = "https://bundles.byte.builders/stable/"
	TestBundleRegistry      = "bytebuilders-testing"
	TestBundleRegistryURL   = "https://bundles.byte.builders/testing/"
)
