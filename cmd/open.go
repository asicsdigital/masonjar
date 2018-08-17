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
	"io"
	"os"
	"path/filepath"

	"github.com/asicsdigital/masonjar/jar"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Create a new directory based on an existing Jar",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.DEBUG.Println("open called")

		jars, _ := jar.ParseJars(viper.GetString("RepoDir"))

		for i := range jars {
			j := jars[i]

			if j.Name() == viper.GetString("JarSource") {
				jww.INFO.Printf("opening jar %v", j.Name())
				viper.Set("CurrentJarName", j.Name())
				viper.Set("CurrentJarPath", j.Path())

				destFs := afero.NewOsFs()
				dfs := &afero.Afero{Fs: destFs}
				destDir := filepath.Join(viper.GetString("JarDestination"), viper.GetString("JarIdentifier"))
				viper.Set("DestRoot", destDir)

				jww.DEBUG.Printf("creating destination directory %v", destDir)
				err := destFs.(*afero.OsFs).MkdirAll(destDir, 0700)

				if err != nil {
					jww.ERROR.Printf("error creating destination directory %v: %v", destDir, err)
				}

				dirCreated, err := dfs.IsDir(destDir)

				if dirCreated {
					err = j.Walk(jarWalkFunc)

					if err != nil {
						jww.ERROR.Printf("error walking jar %v: %v", j.Path(), err)
					}
				} else {
					jww.ERROR.Panic(err)
				}
			}
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
	openCmd.Flags().String("jar", "", "Name of the jar to be used as a src")
	viper.BindPFlag("JarSource", openCmd.Flags().Lookup("jar"))

	openCmd.Flags().String("identifier", "", "Identifier for the jar to be created")
	viper.BindPFlag("JarIdentifier", openCmd.Flags().Lookup("identifier"))

	openCmd.Flags().String("destination", ".", "Path in local filesystem where jar will be created")
	viper.BindPFlag("JarDestination", openCmd.Flags().Lookup("destination"))

}

func jarWalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		jww.WARN.Printf("error walking path %v: %v", path, err)
		return filepath.SkipDir
	}

	if isSkippable(path) {
		jww.INFO.Printf("skipping %v", path)
		return nil
	}

	// no error, not skippable
	// time to actually do something with the path
	srcFs := afero.NewBasePathFs(afero.NewReadOnlyFs(afero.NewOsFs()), viper.GetString("CurrentJarPath"))

	destRoot := viper.GetString("DestRoot")
	destFs := afero.NewBasePathFs(afero.NewOsFs(), destRoot)

	return processJarPath(path, srcFs, destFs)
}

func processJarPath(path string, srcFs afero.Fs, destFs afero.Fs) error {
	jww.DEBUG.Printf("processing path %v", path)
	sfs := &afero.Afero{Fs: srcFs}
	isDir, err := sfs.IsDir(path)

	if isDir {
		fileInfo, _ := srcFs.(*afero.BasePathFs).Stat(path)
		fileMode := fileInfo.Mode()
		return destFs.(*afero.BasePathFs).Mkdir(path, fileMode)
	}

	if isTemplate(path) {
		return nil
	}

	srcFile, err := srcFs.(*afero.BasePathFs).Open(path)

	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	destFile, err := destFs.(*afero.BasePathFs).Create(path)

	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	written, err := io.Copy(destFile, srcFile)

	if err != nil {
		srcRealPath, _ := srcFs.(*afero.BasePathFs).RealPath(path)
		destRealPath, _ := destFs.(*afero.BasePathFs).RealPath(path)
		jww.DEBUG.Printf("copied %v to %v, %v bytes", srcRealPath, destRealPath, written)
	}

	err = destFile.Sync()

	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	return err
}

func isTemplate(path string) bool {
	return false
}

func isSkippable(path string) bool {
	switch path {
	case "/", filepath.Join("/", jar.MetadataFileName):
		return true
	default:
		return false
	}
}
