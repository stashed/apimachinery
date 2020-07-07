/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package restic

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	api_v1alpha1 "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	api_v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	"github.com/stretchr/testify/assert"
	storage "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

var (
	localRepoDir      string
	scratchDir        string
	secretDir         string
	targetPath        string
	password          = "password"
	fileName          = "some-file"
	fileContent       = "hello stash"
	stdinPipeCommand  = Command{Name: "echo", Args: []interface{}{"hello"}}
	stdoutPipeCommand = Command{Name: "cat"}
)

func setupTest(tempDir string) (*ResticWrapper, error) {
	localRepoDir = filepath.Join(tempDir, "repo")
	scratchDir = filepath.Join(tempDir, "scratch")
	secretDir = filepath.Join(tempDir, "secret")
	targetPath = filepath.Join(tempDir, "target")

	if err := os.MkdirAll(localRepoDir, 0777); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(scratchDir, 0777); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(secretDir, 0777); err != nil {
		return nil, err
	}
	err := ioutil.WriteFile(filepath.Join(secretDir, RESTIC_PASSWORD), []byte(password), os.ModePerm)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(targetPath, 0777); err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filepath.Join(targetPath, fileName), []byte(fileContent), os.ModePerm)
	if err != nil {
		return nil, err
	}

	setupOpt := SetupOptions{
		Provider:    storage.ProviderLocal,
		Bucket:      localRepoDir,
		SecretDir:   secretDir,
		ScratchDir:  scratchDir,
		EnableCache: false,
	}

	w, err := NewResticWrapper(setupOpt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func cleanup(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		log.Errorln(err)
	}
}

func TestBackupRestoreDirs(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	backupOpt := BackupOptions{
		BackupPaths: []string{targetPath},
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(backupOut)

	// delete target then restore
	if err = os.RemoveAll(targetPath); err != nil {
		t.Error(err)
	}
	restoreOpt := RestoreOptions{
		RestorePaths: []string{targetPath},
	}
	restoreOut, err := w.RunRestore(restoreOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(restoreOut)

	// check file
	fileContentByte, err := ioutil.ReadFile(filepath.Join(targetPath, fileName))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, fileContent, string(fileContentByte))
}

func TestBackupRestoreStdin(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	backupOpt := BackupOptions{
		StdinPipeCommand: stdinPipeCommand,
		StdinFileName:    fileName,
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("backup output:", backupOut)

	dumpOpt := DumpOptions{
		FileName:          fileName,
		StdoutPipeCommand: stdoutPipeCommand,
	}
	dumpOut, err := w.Dump(dumpOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("dump output:", dumpOut)
}

func TestBackupRestoreWithScheduling(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	w.config.IONice = &ofst.IONiceSettings{
		Class:     types.Int32P(2),
		ClassData: types.Int32P(3),
	}
	w.config.Nice = &ofst.NiceSettings{
		Adjustment: types.Int32P(12),
	}

	backupOpt := BackupOptions{
		BackupPaths: []string{targetPath},
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(backupOut)

	// delete target then restore
	if err = os.RemoveAll(targetPath); err != nil {
		t.Error(err)
	}
	restoreOpt := RestoreOptions{
		RestorePaths: []string{targetPath},
	}
	restoreOut, err := w.RunRestore(restoreOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(restoreOut)

	// check file
	fileContentByte, err := ioutil.ReadFile(filepath.Join(targetPath, fileName))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, fileContent, string(fileContentByte))
}

func TestBackupRestoreStdinWithScheduling(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	w.config.IONice = &ofst.IONiceSettings{
		Class:     types.Int32P(2),
		ClassData: types.Int32P(3),
	}
	w.config.Nice = &ofst.NiceSettings{
		Adjustment: types.Int32P(12),
	}

	backupOpt := BackupOptions{
		StdinPipeCommand: stdinPipeCommand,
		StdinFileName:    fileName,
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("backup output:", backupOut)

	dumpOpt := DumpOptions{
		FileName:          fileName,
		StdoutPipeCommand: stdoutPipeCommand,
	}
	dumpOut, err := w.Dump(dumpOpt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("dump output:", dumpOut)
}

func TestRunParallelBackup(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	// write large (100Mb) sample  file
	largeContent := make([]byte, 104857600)
	fileContent = string(largeContent)

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, 2)
	if err != nil {
		t.Error(err)
	}

	// verify repository stats
	assert.Equal(t, *backupOutput.RepositoryStats.Integrity, true)
	assert.Equal(t, backupOutput.RepositoryStats.SnapshotCount, int64(3))

	// verify each host status
	for i := range backupOutput.HostBackupStats {
		assert.Equal(t, backupOutput.HostBackupStats[i].Phase, api_v1beta1.HostBackupSucceeded)
	}
}

func TestRunParallelRestore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	// write large (100Mb) sample  file
	largeContent := make([]byte, 104857600)
	fileContent = string(largeContent)

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host backup has succeeded
	for i := range backupOutput.HostBackupStats {
		assert.Equal(t, backupOutput.HostBackupStats[i].Phase, api_v1beta1.HostBackupSucceeded)
	}

	// run parallel restore
	restoreOptions, err := newParallelRestoreOptions(tempDir)
	if err != nil {
		t.Error(err)
	}
	restoreOutput, err := w.RunParallelRestore(restoreOptions, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host has been restored successfully
	for i := range restoreOutput.HostRestoreStats {
		assert.Equal(t, restoreOutput.HostRestoreStats[i].Phase, api_v1beta1.HostRestoreSucceeded)
	}

	// verify that restored file contents are identical to the backed up file
	for i := range restoreOptions {
		// check file
		restoredFileContent, err := ioutil.ReadFile(filepath.Join(restoreOptions[i].Destination, targetPath, fileName))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, fileContent, string(restoredFileContent))
	}
}

func TestRunParallelDump(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	// write large (100Mb) sample  file
	largeContent := make([]byte, 104857600)
	fileContent = string(largeContent)

	defer cleanup(tempDir)
	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host backup has succeeded
	for i := range backupOutput.HostBackupStats {
		assert.Equal(t, backupOutput.HostBackupStats[i].Phase, api_v1beta1.HostBackupSucceeded)
	}

	// run parallel dump
	dumpOptions := newParallelDumpOptions()

	dumpOutput, err := w.ParallelDump(dumpOptions, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host has been restored successfully
	for i := range dumpOutput.HostRestoreStats {
		t.Logf("Host: %s, Phase: %s", dumpOutput.HostRestoreStats[i].Hostname, dumpOutput.HostRestoreStats[i].Phase)
		assert.Equal(t, dumpOutput.HostRestoreStats[i].Phase, api_v1beta1.HostRestoreSucceeded)
	}
}

func newParallelBackupOptions() []BackupOptions {
	return []BackupOptions{
		{
			Host:        "host-0",
			BackupPaths: []string{targetPath},
			RetentionPolicy: api_v1alpha1.RetentionPolicy{
				Name:     "keep-last-1",
				KeepLast: 1,
				Prune:    true,
				DryRun:   false,
			},
		},
		{
			Host:        "host-1",
			BackupPaths: []string{targetPath},
			RetentionPolicy: api_v1alpha1.RetentionPolicy{
				Name:     "keep-last-1",
				KeepLast: 1,
				Prune:    true,
				DryRun:   false,
			},
		},
		{
			Host:        "host-2",
			BackupPaths: []string{targetPath},
			RetentionPolicy: api_v1alpha1.RetentionPolicy{
				Name:     "keep-last-1",
				KeepLast: 1,
				Prune:    true,
				DryRun:   false,
			},
		},
	}
}

func newParallelRestoreOptions(tempDir string) ([]RestoreOptions, error) {
	if err := os.MkdirAll(filepath.Join(tempDir, "host-0"), 0777); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(tempDir, "host-1"), 0777); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(tempDir, "host-2"), 0777); err != nil {
		return nil, err
	}

	return []RestoreOptions{
		{
			Host:         "host-0",
			SourceHost:   "",
			RestorePaths: []string{targetPath},
			Destination:  filepath.Join(tempDir, "host-0"),
		},
		{
			Host:         "host-1",
			SourceHost:   "",
			RestorePaths: []string{targetPath},
			Destination:  filepath.Join(tempDir, "host-1"),
		},
		{
			Host:         "host-2",
			SourceHost:   "",
			RestorePaths: []string{targetPath},
			Destination:  filepath.Join(tempDir, "host-2"),
		},
	}, nil
}

func newParallelDumpOptions() []DumpOptions {

	return []DumpOptions{
		{
			Host:              "host-0",
			FileName:          filepath.Join(targetPath, fileName),
			StdoutPipeCommand: stdoutPipeCommand,
		},
		{
			Host:              "host-1",
			FileName:          filepath.Join(targetPath, fileName),
			StdoutPipeCommand: stdoutPipeCommand,
		},
		{
			Host:              "host-2",
			FileName:          filepath.Join(targetPath, fileName),
			StdoutPipeCommand: stdoutPipeCommand,
		},
	}
}
