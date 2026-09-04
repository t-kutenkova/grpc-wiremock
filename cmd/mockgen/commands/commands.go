package commands

import (
	"log"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
)

func CreateCommandRoot() *cobra.Command {
	return mockgenCommand(
		mockgenFirstRun(),
	)
}

func removeTmpDirs() {
	if err := environment.RemoveProcessTmpDirs(afero.NewOsFs()); err != nil {
		log.Println("remove tmp dirs:", err.Error())
	}
}
