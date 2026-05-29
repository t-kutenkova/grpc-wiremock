package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

func ParseProtoFiles(protoFiles []string) (protocontract.SetOfContracts, error) {
	fmt.Fprintf(os.Stdout, "=== PARSE PROTO FILES ===\n")
	fmt.Fprintf(os.Stdout, "Proto files count: %d\n", len(protoFiles))
	for i, f := range protoFiles {
		fmt.Fprintf(os.Stdout, "  File %d: %s\n", i+1, f)
	}

	var contracts protocontract.SetOfContracts

	protoPath := filepath.Dir(sliceutils.FirstOf(protoFiles))

	fmt.Fprintf(os.Stdout, "Proto path: %s\n", protoPath)
	fmt.Fprintf(os.Stdout, "Extra proto paths: %s\n", environment.TmpWellKnownProtosDir)

	descriptors, err := parse([]string{protoPath, environment.TmpWellKnownProtosDir}, protoFiles...)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Parsed descriptors count: %d\n", len(descriptors))

	for _, descriptor := range descriptors {
		convertedContract, errConvert := protocontract.FromProtoDescriptor(descriptor, protoPath)
		if errConvert != nil {
			return nil, fmt.Errorf("convert proto: %w", errConvert)
		}

		contracts = append(contracts, convertedContract)
	}

	fmt.Fprintf(os.Stdout, "=== END PARSE PROTO FILES ===\n")

	return contracts, nil
}

func ParseProtoDir(fs afero.Fs, protoPaths []string) (protocontract.SetOfContracts, error) {
	fmt.Fprintf(os.Stdout, "=== PARSE PROTO DIR ===\n")
	fmt.Fprintf(os.Stdout, "Proto paths count: %d\n", len(protoPaths))
	for i, p := range protoPaths {
		fmt.Fprintf(os.Stdout, "  Path %d: %s\n", i+1, p)
	}

	var contracts protocontract.SetOfContracts

	for _, protoPath := range protoPaths {
		files, err := fsutils.GatherMatchedEntriesInDir(fs, protoPath, onlyProtoFiles)
		if err != nil {
			return nil, fmt.Errorf("get valid files: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Path %s: found %d files\n", protoPath, len(files))

		descriptors, err := parse([]string{protoPath, environment.TmpWellKnownProtosDir}, files...)
		if err != nil {
			return nil, fmt.Errorf("parse proto: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Path %s: parsed %d descriptors\n", protoPath, len(descriptors))

		for _, descriptor := range descriptors {
			convertedContract, errConvert := protocontract.FromProtoDescriptor(descriptor, protoPath)
			if errConvert != nil {
				return nil, fmt.Errorf("convert proto: %w", errConvert)
			}

			contracts = append(contracts, convertedContract)
		}
	}

	fmt.Fprintf(os.Stdout, "=== END PARSE PROTO DIR ===\n")

	return contracts, nil
}

func parse(protoPaths []string, protoFiles ...string) ([]*desc.FileDescriptor, error) {
	fmt.Fprintf(os.Stdout, "=== PARSE PROTO ===\n")
	fmt.Fprintf(os.Stdout, "Proto paths:\n")
	for _, p := range protoPaths {
		fmt.Fprintf(os.Stdout, "  - %s\n", p)
	}
	fmt.Fprintf(os.Stdout, "Proto files:\n")
	for _, f := range protoFiles {
		fmt.Fprintf(os.Stdout, "  - %s\n", f)
	}

	resolvedHeaders, err := protoparse.ResolveFilenames(protoPaths, protoFiles...)
	if err != nil {
		return nil, fmt.Errorf("resolve names: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Resolved headers:\n")
	for _, h := range resolvedHeaders {
		fmt.Fprintf(os.Stdout, "  - %s\n", h)
	}

	parser := &protoparse.Parser{
		IncludeSourceCodeInfo: true,
		ImportPaths:           protoPaths,
	}

	descriptors, err := parser.ParseFiles(resolvedHeaders...)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Parsed %d descriptors\n", len(descriptors))
	for _, d := range descriptors {
		fmt.Fprintf(os.Stdout, "  - %s (package: %s)\n", d.GetName(), d.GetPackage())
	}
	fmt.Fprintf(os.Stdout, "=== END PARSE PROTO ===\n")

	return descriptors, nil
}

func onlyProtoFiles(info os.FileInfo) bool {
	if !info.IsDir() &&
		strings.Contains(info.Name(), ".proto") {
		return true
	}
	return false
}
