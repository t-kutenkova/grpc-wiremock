package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/grpc2http"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
)

func main() {
	command, err := buildGenerateCommand()
	if err != nil {
		log.Fatalln("create cli command:", err.Error())
	}
	if err = command.Execute(); err != nil {
		log.Fatalln("execute cli command:", err.Error())
	}
}

func removeTmpDirs() {
	if err := environment.RemoveProcessTmpDirs(afero.NewOsFs()); err != nil {
		log.Println("remove tmp dirs:", err.Error())
	}
}

type flags struct {
	inputPath  string
	outputPath string
	baseURL    string
}

func buildGenerateCommand() (*cobra.Command, error) {
	var (
		args    flags
		rootCmd = &cobra.Command{
			RunE: func(cmd *cobra.Command, _ []string) error {
				defer removeTmpDirs()

				gen := grpc2http.NewProxyGen(args.inputPath, args.outputPath, args.baseURL, os.Stdout)
				return gen.Generate(cmd.Context())
			},
		}
	)

	rootCmd.Flags().StringVar(&args.inputPath, "input", "", "Directory with contract files")
	rootCmd.Flags().StringVarP(&args.outputPath, "output", "o", "generated_proxy", "Directory for generated proxy")
	rootCmd.Flags().StringVar(&args.baseURL, "base-url", "http://localhost:8080", "Proxy URL")

	if err := rootCmd.MarkFlagRequired("input"); err != nil {
		return nil, fmt.Errorf("make input flag persistent: %w", err)
	}

	return rootCmd, nil
}
