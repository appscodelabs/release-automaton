// Code generated for package templates by go-bindata DO NOT EDIT. (@generated)
// sources:
// changelog.tpl
// release-table.tpl
// shared-changelog.tpl
// standalone-changelog.tpl
package templates

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _changelogTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\x3f\x4b\xc4\x40\x10\xc5\xfb\x7c\x8a\x47\x2e\xc5\xa5\xc8\x6e\xbc\x42\xe4\xc0\x42\xb4\xb0\xb8\xe2\x38\xb5\x12\x8b\x64\xb3\xe6\x22\xe6\x0f\xbb\x13\x10\xc6\xf9\xee\x92\x35\x87\xb1\x48\xb5\xfb\x1e\x33\xbf\xf7\x98\x0d\x98\xa1\x8e\xae\xaf\x46\x43\x87\xa6\xb3\x10\x09\xd6\xc9\x7e\xda\xc2\x07\xb9\x5d\xe8\x87\x82\x2c\xbe\x51\x4d\x4f\xbc\xcb\xf3\xeb\x2c\xbf\xca\xf2\x5d\x0c\x91\x34\x8a\x98\xe1\x8a\xae\xb6\x48\x06\xec\x6f\x03\xf7\xc3\x1a\xf2\x10\x89\x36\x1b\xbc\x32\x83\x5c\xd3\x1e\x9d\x7d\x6f\xbe\x10\xd7\x0d\x9d\xc7\x52\x99\xbe\xd5\x31\x92\x41\xbd\x9c\x0e\x10\x79\xdb\x9e\x89\x06\xbf\xd7\x9a\xf9\xcf\x4d\x17\x70\x37\xc1\x93\xe1\xd2\x69\xc6\xff\xf2\x13\xa7\x9e\x8b\x7a\x15\xa3\xdd\xbc\xa3\xa9\xa8\xf5\x72\xfe\x5f\x7d\x13\x12\x9c\xba\xef\xdb\xb6\x21\x8f\x4c\x24\x42\x16\x02\xfc\x58\x7a\x72\xc8\x71\x83\xc4\xa8\xa7\xc7\xbb\xf5\x2c\x13\xb6\x83\x77\x99\x4c\x31\xab\xb1\x9c\x2e\x33\x35\x67\x86\xed\xaa\xb5\xdf\x4f\x00\x00\x00\xff\xff\xc2\xd8\xc9\x32\xa1\x01\x00\x00")

func changelogTplBytes() ([]byte, error) {
	return bindataRead(
		_changelogTpl,
		"changelog.tpl",
	)
}

func changelogTpl() (*asset, error) {
	bytes, err := changelogTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "changelog.tpl", size: 417, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _releaseTableTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\x41\x4b\x03\x31\x10\x85\xef\xf9\x15\x8f\xda\x83\x1e\x52\x62\x0f\x1e\x04\x0f\xd2\xca\x0a\x2e\x2a\x85\x7a\x59\x3c\xac\xcd\x50\x17\x4a\x02\xc9\xee\x29\x33\xff\x5d\x22\xc6\x68\x69\xf7\xb2\xbc\x37\x93\x79\x7c\xef\x02\x29\x61\xf1\x1a\xbc\x9d\x76\x63\x3b\x38\x82\x08\x36\x74\xa0\x3e\x52\x54\x8a\x4f\x8d\xdf\x28\xc4\xc1\x3b\x70\x59\xc4\xba\x1f\x09\x8c\x6d\xa4\x80\x66\x1a\x6c\x16\xab\xcf\xde\xed\xe9\xe0\xf7\x60\x3c\x4d\x1f\x14\x1c\x8d\x14\xeb\x63\xc5\xfa\xfc\x07\xc6\x59\xf9\x57\x1c\xed\x15\x53\xa5\x84\x90\xe3\x31\x0f\xb8\xbd\x03\x16\x05\x09\x5a\x44\x31\xba\x94\x30\x0f\xc5\x85\xc8\xfb\xe5\x3f\x67\xbb\x69\x21\x72\x85\x6f\xfe\x6a\xff\x70\xda\xfc\x9b\x2d\x8d\xb9\xd1\xe6\x5a\x9b\xe5\x2c\xb7\xc2\xe8\x2a\x7f\x39\xb7\xf6\xbb\x58\x6f\x75\xab\xc7\xfb\xe7\xe6\xa1\x7d\x69\xca\xfc\xb7\xa4\xe3\xc0\xda\x58\x29\x2c\x47\x64\x2e\x72\x16\x22\x5f\x01\x00\x00\xff\xff\xfc\x53\xa3\x37\xb8\x01\x00\x00")

func releaseTableTplBytes() ([]byte, error) {
	return bindataRead(
		_releaseTableTpl,
		"release-table.tpl",
	)
}

func releaseTableTpl() (*asset, error) {
	bytes, err := releaseTableTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "release-table.tpl", size: 440, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sharedChangelogTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x92\xc1\x6e\xf3\x20\x0c\xc7\xef\x3c\x85\x15\x7d\x3d\x12\x3e\xed\x88\xb4\xc3\x54\x4d\xdd\xa1\xdb\xaa\xbd\x40\x85\x12\xb7\xf5\x44\x20\x02\xda\x1e\x98\xdf\x7d\x0a\x4d\xbb\x4c\x5a\x2e\x3d\x70\xb0\xb1\x7f\xfe\xfb\x0f\x39\x4b\xf8\x77\xc2\x10\xc9\x3b\xd0\x8f\x10\xb1\x3b\x61\x80\xfa\x03\x2d\x9a\x88\x20\x99\x85\x94\x52\x24\x4a\x16\x35\x2c\x0f\xc6\xed\xd1\xfa\x3d\x7c\x41\xce\x50\x6f\x82\x6f\x8f\x4d\x5a\x93\x43\x60\x16\x2d\xc6\x26\x50\x9f\xc8\xbb\x49\xad\xe8\xd0\x1d\xb5\x00\xe8\x2f\xd5\xdb\x9c\xc1\xfa\xf3\x30\xe6\x77\xff\x70\x51\xe5\x5c\x8f\x7a\x98\x2b\x60\x1e\x1a\x01\xa8\x45\x97\x68\x47\x18\x34\x34\x57\xb0\x9c\x03\x95\x16\x67\xba\xa9\xe2\x92\xeb\x4d\x40\x97\x34\x9c\xd1\x36\xbe\xc3\x92\x3b\x23\xed\x0f\x49\x43\xce\x7d\x20\x97\x76\x50\x2d\xda\xc5\xff\x87\x72\xaa\x9b\x39\xf5\xab\xf9\xf4\x61\x12\x92\x9b\x86\x1b\x93\x9a\x03\xb3\xb8\xae\x78\x19\x3e\xab\x6f\x70\x64\xac\xb9\xcb\x14\x11\xb1\x19\x5c\xde\x16\x10\xb5\x3f\x1b\x1d\x83\xd5\xa0\x46\x68\x54\x73\x54\xf5\x17\x55\x8d\x10\x75\x73\x58\x09\x63\xc9\x44\x8c\xc3\x2b\xc8\xbb\xb9\xcb\x97\xa7\xb7\xd5\xf3\xfa\x7d\xa5\xca\x67\x12\x39\x27\xec\x7a\x6b\x12\x42\x75\x9b\x55\xa7\xde\x56\x50\x33\x8b\xef\x00\x00\x00\xff\xff\xce\xe8\x59\x1d\x95\x02\x00\x00")

func sharedChangelogTplBytes() ([]byte, error) {
	return bindataRead(
		_sharedChangelogTpl,
		"shared-changelog.tpl",
	)
}

func sharedChangelogTpl() (*asset, error) {
	bytes, err := sharedChangelogTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "shared-changelog.tpl", size: 661, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _standaloneChangelogTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\xc1\x4e\xf3\x30\x0c\xc7\xef\x79\x0a\xab\xfa\x76\x6c\xfb\x89\x63\x24\x0e\x68\x42\xe3\x30\x60\xe2\x05\xa6\x28\xf5\x36\xa3\x34\xa9\x92\x6c\x3b\x18\xbf\x3b\x4a\x37\x46\x11\x1b\x87\x1e\x5c\xff\x6d\xff\xf4\x53\x98\x6b\xf8\x77\xc0\x98\x28\x78\xd0\xf7\x90\xb0\x3f\x60\x84\xe6\x0d\x1d\x9a\x84\x50\x8b\xa8\xba\xae\x55\xa6\xec\x50\xc3\x7c\x67\xfc\x16\x5d\xd8\xc2\x07\x30\x43\xb3\x8a\xa1\xdb\xdb\xbc\x24\x8f\x20\xa2\x3a\x4c\x36\xd2\x90\x29\xf8\x49\x56\xf5\xe8\xf7\x5a\x01\x74\xc1\xa6\x35\x33\x54\xcc\xcd\xf9\xa6\x48\x05\x22\xa5\x09\x40\x1d\xfa\x4c\x1b\xc2\xa8\xc1\x7e\x0d\xd7\xcc\xe0\xc2\xb1\x30\xfd\x3c\x56\x1a\x17\x4c\x91\x71\x83\x37\xfd\x14\xf2\x5a\x64\x30\x11\x7d\xd6\x70\x44\x67\x43\x8f\xe3\xbf\x23\xd2\x76\x97\x35\x30\x0f\x91\x7c\xde\x40\x35\xeb\x66\xff\xef\xc6\xaf\xba\xe8\x69\x9e\xcd\x7b\x88\x93\x92\xfc\xb4\x5c\x99\x6c\x77\x22\x6a\x38\x51\xae\x4f\x2c\xb7\xe8\x47\x27\xe7\xcc\x4d\x2d\x2a\xa1\x2d\x2e\xd7\x63\x98\xba\x6f\xea\x7d\x74\x1a\xda\x32\xd8\x5e\x1b\x6c\xcf\xb9\xd6\xde\x50\xd1\x2a\xe3\xc8\x24\x4c\xc5\x7c\xfd\xd7\xa6\xf9\xd3\xc3\xcb\xe2\x71\xf9\xba\xf8\xb5\xa1\x3c\x0b\xc5\x9c\xb1\x1f\x9c\xc9\x08\xd5\xe5\x58\x93\x07\x57\x41\x23\xa2\x3e\x03\x00\x00\xff\xff\x19\x56\x0b\xd5\x5f\x02\x00\x00")

func standaloneChangelogTplBytes() ([]byte, error) {
	return bindataRead(
		_standaloneChangelogTpl,
		"standalone-changelog.tpl",
	)
}

func standaloneChangelogTpl() (*asset, error) {
	bytes, err := standaloneChangelogTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "standalone-changelog.tpl", size: 607, mode: os.FileMode(420), modTime: time.Unix(1573722179, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"changelog.tpl":            changelogTpl,
	"release-table.tpl":        releaseTableTpl,
	"shared-changelog.tpl":     sharedChangelogTpl,
	"standalone-changelog.tpl": standaloneChangelogTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"changelog.tpl":            {changelogTpl, map[string]*bintree{}},
	"release-table.tpl":        {releaseTableTpl, map[string]*bintree{}},
	"shared-changelog.tpl":     {sharedChangelogTpl, map[string]*bintree{}},
	"standalone-changelog.tpl": {standaloneChangelogTpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0o755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
