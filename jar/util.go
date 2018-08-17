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
	"path/filepath"

	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

func isValidJar(path string) bool {
	fs := afero.NewBasePathFs(afero.NewReadOnlyFs(afero.NewOsFs()), path)

	_, err := fs.(*afero.BasePathFs).Stat(filepath.Join("/", MetadataFileName))

	if err != nil {
		dirRealPath, _ := fs.(*afero.BasePathFs).RealPath("/")
		jww.DEBUG.Printf("directory %v has no metadata file %v", dirRealPath, MetadataFileName)
		return false
	}

	return true
}

func ParseJars(repoDir string) ([]Jar, error) {
	jww.DEBUG.Printf("parsing jars from %v", repoDir)
	fs := afero.NewBasePathFs(afero.NewReadOnlyFs(afero.NewOsFs()), repoDir)
	afs := &afero.Afero{Fs: fs}

	var jars []Jar

	files, err := afs.ReadDir("/")

	if err != nil {
		jww.ERROR.Println(err)
		return jars, err
	}

	for i := range files {
		fileName, err := fs.(*afero.BasePathFs).RealPath(files[i].Name())
		j, err := NewJar(fileName)

		if err == nil {
			jww.INFO.Printf("parsed %v as Jar %v", j.Path(), j.Name())
			jars = append(jars, j)
		} else {
			jww.WARN.Println(err)
		}
	}

	return jars, err
}
