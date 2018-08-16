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
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Get the latest masonjar definitions",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.DEBUG.Println("update called")

		repoDir := viper.GetString("RepoDir")
		err := pullRepo(repoDir, viper.GetString("RepoRemote"))

		switch err {
		case git.NoErrAlreadyUpToDate:
			jww.INFO.Println(err)
			err = nil
		case git.ErrRepositoryNotExists:
			jww.INFO.Println(err)
			err = cloneRepo(repoDir, viper.GetString("RepoUrl"))
		default:
			jww.ERROR.Println(err)
		}

		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().String("repository", "", "Git repo containing masonjar definitions (default is https://github.com/asicsdigital/masonjars)")
	viper.SetDefault("RepoUrl", "https://github.com/asicsdigital/masonjars")
	viper.BindPFlag("RepoUrl", updateCmd.Flags().Lookup("repository"))

	updateCmd.Flags().String("remote", "", "Remote of Git repo containing masonjar definitions (default is 'origin')")
	viper.SetDefault("RepoRemote", "origin")
	viper.BindPFlag("RepoRemote", updateCmd.Flags().Lookup("remote"))
}

func cloneRepo(destDir string, repoUrl string) error {
	jww.DEBUG.Println("cloneRepo called")
	_, err := git.PlainClone(destDir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stderr,
	})

	jww.DEBUG.Println("cloneRepo returned")
	return err
}

func pullRepo(destDir string, repoRemote string) error {
	jww.DEBUG.Println("pullRepo called")
	r, err := git.PlainOpen(destDir)

	if err != nil {
		return err
	}

	w, err := r.Worktree()

	err = w.Pull(&git.PullOptions{
		RemoteName: repoRemote,
		Progress:   os.Stderr,
	})

	jww.DEBUG.Println("pullRepo returned")
	return err
}
