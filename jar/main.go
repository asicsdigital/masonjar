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
	"github.com/spf13/viper"
)

type Jar interface {
	Name() string
	Path() string
	Prefix() string
	Metadata() *viper.Viper
	ParseMetadata(string) (*viper.Viper, error)
	Walk(filepath.WalkFunc) error
}

type MasonJar struct {
	name     string
	path     string
	metadata *viper.Viper
}

func (j *MasonJar) Name() string {
	return j.name
}

func (j *MasonJar) Path() string {
	return j.path
}

func (j *MasonJar) Metadata() *viper.Viper {
	return j.metadata
}

func (j *MasonJar) Prefix() string {
	return j.Metadata().GetString("prefix")
}

func (j *MasonJar) Walk(walkFn filepath.WalkFunc) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), j.Path())
	afs := &afero.Afero{Fs: fs}

	err := afs.Walk("/", walkFn)

	return err
}

func NewJar(path string) (*MasonJar, error) {
	j := new(MasonJar)
	j.path = path

	_, name := filepath.Split(path)

	j.name = name

	metadata, err := j.ParseMetadata(MetadataFileName)

	if err != nil {
		return nil, fmt.Errorf("%v is not a valid Jar directory", path)
	}

	j.metadata = metadata
	return j, err
}
