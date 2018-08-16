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

package cmd

import (
	"github.com/asicsdigital/masonjar/jar"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Create a new directory based on an existing MasonJar",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.DEBUG.Println("open called")

		jars, _ := parseJars(viper.GetString("RepoDir"))

		for i := range jars {
			jww.ERROR.Println(jars[i].Name())
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	openCmd.Flags().String("identifier", "", "Identifier for the jar to be created")
	viper.BindPFlag("JarIdentifier", openCmd.Flags().Lookup("identifier"))
}

func parseJars(repoDir string) ([]jar.MasonJar, error) {
	fs := afero.NewBasePathFs(afero.NewReadOnlyFs(afero.NewOsFs()), repoDir)
	afs := &afero.Afero{Fs: fs}

	var jars []jar.MasonJar

	files, err := afs.ReadDir("/")

	if err != nil {
		jww.ERROR.Println(err)
		return jars, err
	}

	for i := range files {
		fileName, err := fs.(*afero.BasePathFs).RealPath(files[i].Name())
		j, err := jar.NewJar(fileName)

		if err == nil {
			jww.INFO.Printf("parsed %v as MasonJar %v", j.Path(), j.Name())
			jars = append(jars, j)
		} else {
			jww.WARN.Println(err)
		}
	}

	return jars, err
}
