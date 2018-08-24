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
	"strings"

	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

func MatchJar(target string, jars []Jar, walkFunc filepath.WalkFunc) bool {
	jww.DEBUG.Printf("matching %v against %v jars", target, len(jars))

	matchedJar := false

	for i := range jars {
		j := jars[i]

		if j.Name() == target {
			jww.INFO.Printf("opening jar %v", j.Name())
			matchedJar = true
			viper.Set("CurrentJarName", j.Name())
			viper.Set("CurrentJarPath", j.Path())
			viper.Set("CurrentJarMetadata", j.Metadata())

			destFs := afero.NewOsFs()
			dfs := &afero.Afero{Fs: destFs}
			destDir := filepath.Join(viper.GetString("JarDestination"), strings.Join([]string{j.Prefix(), viper.GetString("JarIdentifier")}, ""))
			viper.Set("DestRoot", destDir)

			dirExists, err := dfs.DirExists(destDir)
			if !dirExists {
				jww.DEBUG.Printf("creating destination directory %v", destDir)
				err := destFs.(*afero.OsFs).MkdirAll(destDir, 0700)

				if err != nil {
					jww.ERROR.Printf("error creating destination directory %v: %v", destDir, err)
				}
			}

			err = j.Walk(walkFunc)

			if err != nil {
				jww.ERROR.Printf("error walking jar %v: %v", j.Path(), err)
			}
		}
	}

	return matchedJar
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
			jww.INFO.Printf("parsed %v as jar %v", j.Path(), j.Name())
			jars = append(jars, j)
		} else {
			jww.WARN.Println(err)
		}
	}

	return jars, err
}
