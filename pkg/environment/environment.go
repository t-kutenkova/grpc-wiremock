package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

const (
	NginxConfigsPath          = "/etc/nginx/http.d"
	SupervisordConfigsPath    = "/etc/supervisord/mocks"
	SupervisordConfigsDirPath = "/etc/supervisord"
	SupervisordMainConfigPath = "/etc/supervisord/supervisord.conf"

	DefaultWiremockConfigPath = "/home/mock"

	AnnotationsPath     = "google/api/annotations.proto"
	AnnotationsHttpPath = "google/api/http.proto"

	DefaultCertificatesPath = "/etc/ssl"

	TrustedCertificatePath = "/etc/ssl/certs/ca-certificates.crt"

	CAKeyFile  = "mock/share/mockCA.key"
	CACertFile = "mock/share/mockCA.crt"

	CertKeyFile  = "mock/mock.key"
	CertCertFile = "mock/mock.crt"
)

// Tmp directories are unique per process. Mocks and proxy generators run in
// parallel, share the same tmp root and wipe it on start (see CleanTmpDirs),
// so with fixed paths they used to delete each other's files mid-run.
var (
	TmpUnifiedContractsDir     = processTmpDir("unified-contracts")
	TmpOverwrittenContractsDir = processTmpDir("overwritten-contracts")
	TmpGeneratedPackagesDir    = processTmpDir("generated-packages")

	TmpWellKnownProtosDir  = processTmpDir("proto-includes")
	TmpAnnotationProtosDir = processTmpDir("proto-annotations")
)

func processTmpDir(name string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", name, os.Getpid()))
}

func contractsTmpDirs() []string {
	return []string{
		TmpUnifiedContractsDir,
		TmpGeneratedPackagesDir,
		TmpOverwrittenContractsDir,
	}
}

func DumpProtos(fs afero.Fs) error {
	protoToCopy := map[string]string{
		"proto-includes":    TmpWellKnownProtosDir,
		"proto-annotations": TmpAnnotationProtosDir,
	}

	staticFS := static.FromEmbed()

	for sourcePath, targetPath := range protoToCopy {
		if err := fsutils.CopyDir(staticFS, fs, sourcePath, targetPath, true); err != nil {
			return fmt.Errorf("copy protos from embed fs: %w", err)
		}
	}

	return nil
}

func CleanTmpDirs(fs afero.Fs) error {
	if err := fsutils.RemoveTmpDirs(fs, contractsTmpDirs()...); err != nil {
		return fmt.Errorf("remove tmp dir: %w", err)
	}

	return nil
}

// RemoveProcessTmpDirs drops everything this process created in the tmp root.
// Unlike CleanTmpDirs it does not recreate the directories, so it is meant to
// be deferred by a command: without it every run would leave its dirs behind.
func RemoveProcessTmpDirs(fs afero.Fs) error {
	tmpDirs := append(contractsTmpDirs(), TmpWellKnownProtosDir, TmpAnnotationProtosDir)

	for _, path := range tmpDirs {
		if err := fs.RemoveAll(path); err != nil {
			return fmt.Errorf("remove tmp dir '%s': %w", path, err)
		}
	}

	return nil
}

func CleanConfigs(fs afero.Fs, path string) error {
	match := func(info os.FileInfo) bool {
		if info.IsDir() {
			return false
		}

		const (
			mockPrefix = "mock-"
			mockSuffix = ".conf"
		)

		return strings.HasPrefix(info.Name(), mockPrefix) &&
			strings.HasSuffix(info.Name(), mockSuffix)
	}

	entries, err := fsutils.GatherMatchedEntriesInDir(fs, path, match)
	if err != nil {
		return fmt.Errorf("gather matched entries: %w", err)
	}

	for _, entryPath := range entries {
		if err = fs.Remove(entryPath); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
	}

	return nil
}

func IsCAExists(fs afero.Fs, output string) bool {
	certFilePath := filepath.Join(output, CACertFile)

	_, err := fs.Stat(certFilePath)
	if err == nil {
		return true
	}

	return false
}
