package builder

import (
	"fmt"
	"os"

	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/blacklist"
)

type protoUpdater interface {
	Name() string
	Update(contract *desc.FileDescriptor) (*desc.FileDescriptor, error)
}

func UpdateContracts(contracts []*desc.FileDescriptor, updaters ...protoUpdater) ([]*desc.FileDescriptor, error) {
	var updatedDescriptors []*desc.FileDescriptor

	fmt.Fprintf(os.Stdout, "=== UPDATE CONTRACTS ===\n")
	fmt.Fprintf(os.Stdout, "Total contracts to process: %d\n", len(contracts))

	for i, descriptor := range contracts {
		goPackage := protoGoPackage(descriptor)

		fmt.Fprintf(os.Stdout, "Contract %d: %s\n", i+1, descriptor.GetName())
		fmt.Fprintf(os.Stdout, "  go_package: %s\n", goPackage)

		if blacklist.IsGoogleAPIContract(goPackage) {
			fmt.Fprintf(os.Stdout, "  Skipping: Google API contract\n")
			updatedDescriptors = append(updatedDescriptors, descriptor)
			continue
		}

		if blacklist.IsDeliveredWithProtoc(goPackage) {
			fmt.Fprintf(os.Stdout, "  Skipping: Delivered with protoc\n")
			updatedDescriptors = append(updatedDescriptors, descriptor)
			continue
		}

		fmt.Fprintf(os.Stdout, "  Will update with %d updaters\n", len(updaters))

		updatedDescriptor := descriptor

		var err error
		for _, updater := range updaters {
			updatedDescriptor, err = updater.Update(updatedDescriptor)
			if err != nil {
				return nil, fmt.Errorf("update proto descriptor with %s: %w", updater.Name(), err)
			}
		}

		updatedDescriptors = append(updatedDescriptors, updatedDescriptor)
	}

	fmt.Fprintf(os.Stdout, "=== END UPDATE CONTRACTS ===\n")

	return updatedDescriptors, nil
}

// TODO separate common functions
func protoGoPackage(descriptor *desc.FileDescriptor) string {
	if descriptor.GetFileOptions() == nil {
		return ""
	}

	goPackage := descriptor.GetFileOptions().GoPackage

	if goPackage == nil {
		return ""
	}

	return *goPackage
}
