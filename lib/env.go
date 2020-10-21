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
	"net/url"
	"path"
	"strings"
)

func RepoURL2EnvKey(repoURL string) string {
	if !strings.Contains(repoURL, "://") {
		repoURL = "https://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	return toEnvKey(path.Join(u.Path, "tag"))
}

func Key2EnvKey(key string) string {
	return toEnvKey(path.Join(key, "version"))
}

func toEnvKey(key string) string {
	key = strings.Trim(key, "/")
	key = strings.ReplaceAll(key, "/", "_")
	key = strings.ReplaceAll(key, "-", "_")
	return strings.ToUpper(key)
}
