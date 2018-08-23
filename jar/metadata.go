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
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const MetadataFileName = "metadata"

func (j *MasonJar) ParseMetadata(filename string) (*viper.Viper, error) {
	path := j.Path()

	jww.DEBUG.Printf("parsing metadata for jar %v (path: %v, filename: %v)", j.Name(), path, filename)

	config := viper.New()

	config.SetConfigName(filename)
	config.AddConfigPath(j.Path())

	err := config.ReadInConfig()

	if err != nil {
		jww.WARN.Printf("unable to parse metadata: %v", err)
		return nil, err
	}

	return config, err
}
