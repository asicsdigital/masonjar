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
	"fmt"
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
	Short: "Create a new directory based on an existing jar",
	Long: `Create a new directory by making a copy of an existing jar.

Required parameters are --jar (which must match one of the jar names output by
"masonjar list") and -identifier (a unique identifier for the copy of the jar).`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.DEBUG.Println("open called")

		jars, _ := jar.ParseJars(viper.GetString("RepoDir"))

		targetJar := viper.GetString("JarSource")

		if !jar.MatchJar(targetJar, jars, jarWalkFunc) {
			jww.ERROR.Printf("Unable to find a jar matching '%v'.  Use `masonjar list` to list available jars.", targetJar)
			os.Exit(1)
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
	openCmd.Flags().String("jar", "", "Name of the jar to be used as a source (required)")
	openCmd.MarkFlagRequired("jar")
	viper.BindPFlag("JarSource", openCmd.Flags().Lookup("jar"))

	openCmd.Flags().String("identifier", "", "Identifier for the jar to be created (required)")
	openCmd.MarkFlagRequired("identifier")
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

	metadata := viper.Get("CurrentJarMetadata").(*viper.Viper)

	if isTemplate(path, metadata) {
		return nil
	}

	err = jar.CopyFile(path, srcFs, destFs)

	return err
}

func isTemplate(path string, metadata *viper.Viper) bool {
	filename := filepath.Base(path)
	template_spec := fmt.Sprintf("%s.%s", "templates", filename)

	if !metadata.IsSet(template_spec) {
		jww.DEBUG.Printf("no template specification for %v", template_spec)
		return false
	}

	jww.INFO.Printf("found template specification for %v", template_spec)
	return true
}

func isSkippable(path string) bool {
	metadataMatch, err := filepath.Match(fmt.Sprintf("/%s.*", jar.MetadataFileName), path)

	if err != nil {
		jww.WARN.Println(err)
		return false
	}

	if metadataMatch {
		return true
	}

	switch path {
	case "/", "/templates":
		return true
	default:
		return false
	}
}
