// Copyright © 2018 Steve Huff <steve.huff@asics.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package jar

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

const MetadataFileName = "metadata.json"

type MasonJar interface {
	Name() string
	Path() string
}

type masonJar struct {
	name string
	path string
}

func (j *masonJar) Name() string {
	return j.name
}

func (j *masonJar) Path() string {
	return j.path
}

func NewJar(path string) (*masonJar, error) {
	if !isValidJar(path) {
		return nil, fmt.Errorf("%v is not a valid MasonJar directory", path)
	}

	j := new(masonJar)
	j.path = path

	_, name := filepath.Split(path)

	j.name = name
	return j, nil
}

func isValidJar(path string) bool {
	fs := afero.NewBasePathFs(afero.NewReadOnlyFs(afero.NewOsFs()), path)

	_, err := fs.(*afero.BasePathFs).Open(filepath.Join("/", MetadataFileName))

	if err != nil {
		dirRealPath, _ := fs.(*afero.BasePathFs).RealPath("/")
		jww.WARN.Printf("directory %v has no metadata file %v", dirRealPath, MetadataFileName)
		return false
	}

	return true
}
