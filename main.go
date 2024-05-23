package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bareclone [repository]",
	Short: "bareclone lets you setup a bare repo with git worktrees",
	Long:  "An automated way to setup git worktrees with a bare repo",
	Args:  cobra.ExactArgs(1),
	Run:   rootCommand,
}

type model struct {
	repoUrl      string
	hiddenFolder string
}

var state = model{}

func main() {
	setupFlags()
	_ = rootCmd.Execute()
}

func setupFlags() {
	rootCmd.Flags().StringVarP(&state.hiddenFolder, "folder", "f", ".bare", "Name of the folder where to store all git files")
}

func rootCommand(cmd *cobra.Command, args []string) {
	// this is validated to exist by ExactArgs(1)
	state.repoUrl = strings.TrimSpace(args[0])

	fmt.Println("Clonning repo...")
	if err := clone(state.repoUrl, state.hiddenFolder); err != nil {
		os.Exit(1)
	}

	fmt.Println("Setting up .git file...")
	if err := setupGitFile(state.repoUrl, state.hiddenFolder); err != nil {
		fmt.Printf("Ups! Looks like we got an error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Fixing branches configuration...")
	if err := configBranches(state.repoUrl); err != nil {
		os.Exit(1)
	}
}

func getRepoName(repoURL string) string {
	parts := strings.Split(strings.TrimSuffix(repoURL, ".git"), "/")
	return parts[len(parts)-1]
}

func clone(repoUrl, hiddenFolder string) error {
	repoName := getRepoName(repoUrl)
	cmd := exec.Command("git", "clone", "--bare", repoUrl, fmt.Sprintf("%s/%s", repoName, hiddenFolder))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func setupGitFile(repoURL, hiddenFolder string) error {
	repoName := getRepoName(repoURL)
	content := fmt.Sprintf("gitdir: ./%s\n", hiddenFolder)
	return os.WriteFile(fmt.Sprintf("./%s/.git", repoName), []byte(content), 0600)
}

func configBranches(repoURL string) (err error) {
	repoName := getRepoName(repoURL)

	dir := fmt.Sprintf("./%s", repoName)
	if err = os.Chdir(dir); err != nil {
		return fmt.Errorf("repository folder %s does not exist: %w", repoName, err)
	}

	defer func() {
		chErr := os.Chdir("..")
		if chErr != nil {
			err = chErr
		}
	}()

	cmd := exec.Command("git", "config", "remote.origin.fetch", `+refs/heads/*:refs/remotes/origin/*`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
