package main

import (
	"crab/crabe"
	"crab/exec"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func maisn() {

	var executor = exec.CommandExecutor{}
	var rootCmd = &cobra.Command{
		Use:   "crab",
		Run: func(cmd *cobra.Command, args []string) {
			var template string
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				bytes, _ := ioutil.ReadAll(os.Stdin)
				template = strings.TrimSpace(string(bytes))
			}
			if template == "" {
				if len(args) == 0 {
					Error(cmd, args, errors.New("missing args"))
				}
				bytes, err := ioutil.ReadFile(args[0])
				if err != nil {
					Error(cmd, args, err)
				}
				template = string(bytes)
			}
			fmt.Println(template)
			fmt.Println("----")
			// TODO cue
			// TODO kubectl

			output, err := executor.ExecuteCommandWithCombinedOutput("kubectl", "version", "--client", "-o", "json")
			fmt.Println(output, err)
		},
	}

	var apiServerURL string
	rootCmd.Flags().StringVarP(&apiServerURL, "server", "s", "", "the address and port of the Kubernetes API server")
	// _ = rootCmd.MarkFlagRequired("server")

	var context string
	rootCmd.Flags().StringVarP(&context, "context", "c", "", "the context of this instance")
	_ = rootCmd.MarkFlagRequired("context")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func Error(cmd *cobra.Command, args []string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}

func main() {
	crabe.CheckCUEVersion()
}


