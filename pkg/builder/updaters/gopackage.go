package updaters

import (
	"fmt"
	"os"
	"path/filepath"

	protodesc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/blacklist"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

type goPackageUpdater struct{}

func NewGoPackageUpdater() *goPackageUpdater {
	return &goPackageUpdater{}
}

func (u *goPackageUpdater) Name() string {
	return "go package updater"
}

func (u *goPackageUpdater) Update(contract *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	return overwriteGoPackage(contract)
}

func overwriteGoPackage(descriptor *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	fmt.Fprintf(os.Stdout, "=== UPDATE GO PACKAGE ===\n")
	fmt.Fprintf(os.Stdout, "File name: %s\n", descriptor.GetName())
	fmt.Fprintf(os.Stdout, "Package: %s\n", descriptor.GetPackage())

	if blacklist.IsDeliveredWithProtoc(descriptor.GetName()) {
		fmt.Fprintf(os.Stdout, "Skipping: file is delivered with protoc\n")
		fmt.Fprintf(os.Stdout, "=== END UPDATE GO PACKAGE ===\n")
		return descriptor, nil
	}

	updatedGopackage, err := updateGopackage(descriptor)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error updating go package: %v\n", err)
		fmt.Fprintf(os.Stdout, "=== END UPDATE GO PACKAGE ===\n")
		return descriptor, nil
	}

	fmt.Fprintf(os.Stdout, "Original go_package: %v\n", descriptor.GetFileOptions().GetGoPackage())
	fmt.Fprintf(os.Stdout, "Updated go_package: %s\n", updatedGopackage)

	fileBuilder, err := builder.FromFile(descriptor)
	if err != nil {
		return nil, fmt.Errorf("load builder: %w", err)
	}

	fileBuilder.Options = &protodesc.FileOptions{GoPackage: &updatedGopackage}

	updatedDescriptor, err := fileBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build: %w", err)
	}

	fmt.Fprintf(os.Stdout, "=== END UPDATE GO PACKAGE ===\n")

	return updatedDescriptor, nil
}

func updateGopackage(descriptor *desc.FileDescriptor) (string, error) {
	const generatedGoPackage = "github.com/generated/gopackage"

	packageName := descriptor.GetPackage()

	if descriptor.GetFileOptions() == nil {
		return newGoPackage(packageName, generatedGoPackage), nil
	}

	goPackage := descriptor.GetFileOptions().GoPackage
	if goPackage == nil {
		return newGoPackage(packageName, generatedGoPackage), nil
	}

	if blacklist.IsDeliveredWithProtoc(*goPackage) {
		return "", fmt.Errorf("skip this contract")
	}

	return newGoPackage(packageName, *goPackage), nil
}

func newGoPackage(packageName, goPackageName string) string {
	newPackageName := strutils.ResolveNameIfCollides(packageName)
	return filepath.Join("grpc-proxy/pkg", filepath.Dir(goPackageName), strutils.ToSnakeCase(newPackageName))
}
