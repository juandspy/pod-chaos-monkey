package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	namespace     string // target namespace
	labelSelector string // label selector for listing the pods

	chaosCmd = &cobra.Command{
		Use:   "pod-chaos-monkey",
		Short: "pod-chaos-monkey - a simplistic implementation of Chaos engineering for Kubernetes",
		Long: `pod-chaos-monkey is a CLI tool used to kill a random pod in the given namespace. 
    
The default namespace is "default", but it can be changed via CLI args.
It is also possible to filter the candidates by using a label selector.`,
		Run: func(cmd *cobra.Command, args []string) {
			NewPodDeleter().DeleteRandomPod(namespace, labelSelector)
		},
	}
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	chaosCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Target namespace")
	chaosCmd.Flags().StringVar(&labelSelector, "label-selector", "app!=pod-chaos-monkey", "Label selector to select the pods to be deleted")

}

func RunChaos() {
	if err := chaosCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("error executing the CLI command")
	}
}
