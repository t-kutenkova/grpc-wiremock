package proxy

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/errgroup"
	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

// GeneratePackages runs proto compiler with 'go' and 'go-grpc'
// plugins and places generated golang packages to the proxy template.
func (g proxyGenerator) GeneratePackages(ctx context.Context, contracts contract.SetOfContracts, logs io.Writer) error {
	packagesPath := filepath.Join(g.output, "pkg")
	return g.generatePackages(ctx, contracts, packagesPath, logs)
}

func (g proxyGenerator) generatePackages(ctx context.Context, contracts contract.SetOfContracts, output string, logs io.Writer) error {
	// Log all contracts being compiled for debugging proto registration issues
	if _, err := fmt.Fprintln(logs, "=== PACKAGES GENERATION ==="); err != nil {
		return fmt.Errorf("print: %w", err)
	}
	if _, err := fmt.Fprintln(logs, "Total contracts to generate:", len(contracts)); err != nil {
		return fmt.Errorf("print: %w", err)
	}
	for i, c := range contracts {
		if _, err := fmt.Fprintln(logs, fmt.Sprintf("Contract %d:", i+1)); err != nil {
			return fmt.Errorf("print: %w", err)
		}
		if _, err := fmt.Fprintln(logs, "  HeaderPath:", c.HeaderPath); err != nil {
			return fmt.Errorf("print: %w", err)
		}
		if len(c.ImportsPaths) > 0 {
			if _, err := fmt.Fprintln(logs, "  ImportsPaths:"); err != nil {
				return fmt.Errorf("print: %w", err)
			}
			for _, imp := range c.ImportsPaths {
				if _, err := fmt.Fprintln(logs, "    -", imp); err != nil {
					return fmt.Errorf("print: %w", err)
				}
			}
		}
	}
	if _, err := fmt.Fprintln(logs, "=== END PACKAGES GENERATION ==="); err != nil {
		return fmt.Errorf("print: %w", err)
	}

	errG, errCtx := errgroup.WithContext(ctx)
	errG.SetLimit(len(contracts))

	for _, contractToGenerate := range contracts {
		contractToGenerate := contractToGenerate

		errG.Go(func() error {
			tmpDirForContract, err := afero.TempDir(g.fs, environment.TmpGeneratedPackagesDir, "")
			if err != nil {
				return fmt.Errorf("create tmp dir for contract '%s': %w", contractToGenerate.HeaderPath, err)
			}

			if err = g.compiler.CompileToGo(errCtx, contractToGenerate, tmpDirForContract, logs); err != nil {
				return fmt.Errorf("compile '%s': %w", contractToGenerate.HeaderPath, err)
			}

			targetPkgDir := filepath.Join(tmpDirForContract, pkgDir)

			if err = fsutils.CopyDir(g.fs, g.fs, targetPkgDir, output, true); err != nil {
				return fmt.Errorf("move generated contracts into proxy folder: %w", err)
			}

			return nil
		})
	}

	if err := errG.Wait(); err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	return nil
}
