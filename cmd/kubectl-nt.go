package main

import (
	"os"

	"github.com/spf13/cobra"
	cmdnt "github.com/sukimoyoi/kubectl-nt/pkg/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	cmd := NewCmdNeoTree()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}

func NewCmdNeoTree() *cobra.Command {
	// type options struct {
	// 	kubeConfig string
	// 	server     string
	// }
	// o = &options{}

	cmd := &cobra.Command{
		Use:   "kubectl-nt",
		Short: "get tree of Kubernetes resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			// err := cmdnt.GetTree("pvc", "myvol-block-myapp-sts-0", "default")
			err := cmdnt.NeoTree(args)
			return err
		},
		// SilenceErrors: true,
		// SilenceUsage:  true,
	}
	// TODO: enable to use specific kubeconfig
	// cmd.Flags().StringVar(&o.kubeConfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	// cmd.Flags().StringVar(&o.server, "server", "", "The address and port of the Kubernetes API server")
	return cmd
}
