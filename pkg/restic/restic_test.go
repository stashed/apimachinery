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

	"github.com/stretchr/testify/assert"
	"gomodules.xyz/pointer"
	"k8s.io/klog/v2"
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

var testTargetRef = api_v1beta1.TargetRef{
	APIVersion: "test.stash.appscode.com",
	Kind:       "UnitTest",
	Name:       "unit-test-demo",
}

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
		klog.Errorln(err)
	}
}

func TestInitializeRepository(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}
	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryAlreadyExist_AfterInitialization(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}
	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}
	repoExist := w.RepositoryAlreadyExist()
	assert.Equal(t, true, repoExist)
}

func TestRepositoryAlreadyExist_WithoutInitialization(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}
	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	repoExist := w.RepositoryAlreadyExist()
	assert.Equal(t, false, repoExist)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
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
	backupOut, err := w.RunBackup(backupOpt, testTargetRef)
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
	restoreOut, err := w.RunRestore(restoreOpt, testTargetRef)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	backupOpt := BackupOptions{
		StdinPipeCommands: []Command{stdinPipeCommand},
		StdinFileName:     fileName,
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("backup output:", backupOut)

	dumpOpt := DumpOptions{
		FileName:           fileName,
		StdoutPipeCommands: []Command{stdoutPipeCommand},
	}
	dumpOut, err := w.Dump(dumpOpt, testTargetRef)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	w.config.IONice = &ofst.IONiceSettings{
		Class:     pointer.Int32P(2),
		ClassData: pointer.Int32P(3),
	}
	w.config.Nice = &ofst.NiceSettings{
		Adjustment: pointer.Int32P(12),
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
	backupOut, err := w.RunBackup(backupOpt, testTargetRef)
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
	restoreOut, err := w.RunRestore(restoreOpt, testTargetRef)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	w.config.IONice = &ofst.IONiceSettings{
		Class:     pointer.Int32P(2),
		ClassData: pointer.Int32P(3),
	}
	w.config.Nice = &ofst.NiceSettings{
		Adjustment: pointer.Int32P(12),
	}

	backupOpt := BackupOptions{
		StdinPipeCommands: []Command{stdinPipeCommand},
		StdinFileName:     fileName,
		RetentionPolicy: api_v1alpha1.RetentionPolicy{
			Name:     "keep-last-1",
			KeepLast: 1,
			Prune:    true,
			DryRun:   false,
		},
	}
	backupOut, err := w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("backup output:", backupOut)

	dumpOpt := DumpOptions{
		FileName:           fileName,
		StdoutPipeCommands: []Command{stdoutPipeCommand},
	}
	dumpOut, err := w.Dump(dumpOpt, testTargetRef)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, testTargetRef, 2)
	if err != nil {
		t.Error(err)
	}
	// verify each host status
	for i := range backupOutput.BackupTargetStatus.Stats {
		assert.Equal(t, backupOutput.BackupTargetStatus.Stats[i].Phase, api_v1beta1.HostBackupSucceeded)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, testTargetRef, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host backup has succeeded
	for i := range backupOutput.BackupTargetStatus.Stats {
		assert.Equal(t, backupOutput.BackupTargetStatus.Stats[i].Phase, api_v1beta1.HostBackupSucceeded)
	}

	// run parallel restore
	restoreOptions, err := newParallelRestoreOptions(tempDir)
	if err != nil {
		t.Error(err)
	}
	restoreOutput, err := w.RunParallelRestore(restoreOptions, testTargetRef, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host has been restored successfully
	for i := range restoreOutput.RestoreTargetStatus.Stats {
		assert.Equal(t, restoreOutput.RestoreTargetStatus.Stats[i].Phase, api_v1beta1.HostRestoreSucceeded)
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

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
	}

	backupOpts := newParallelBackupOptions()
	backupOutput, err := w.RunParallelBackup(backupOpts, testTargetRef, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host backup has succeeded
	for i := range backupOutput.BackupTargetStatus.Stats {
		assert.Equal(t, backupOutput.BackupTargetStatus.Stats[i].Phase, api_v1beta1.HostBackupSucceeded)
	}

	// run parallel dump
	dumpOptions := newParallelDumpOptions()

	dumpOutput, err := w.ParallelDump(dumpOptions, testTargetRef, 2)
	if err != nil {
		t.Error(err)
	}

	// verify that all host has been restored successfully
	for i := range dumpOutput.RestoreTargetStatus.Stats {
		t.Logf("Host: %s, Phase: %s",
			dumpOutput.RestoreTargetStatus.Stats[i].Hostname,
			dumpOutput.RestoreTargetStatus.Stats[i].Phase,
		)
		assert.Equal(t, dumpOutput.RestoreTargetStatus.Stats[i].Phase, api_v1beta1.HostRestoreSucceeded)
	}
}

func TestIncludeExcludePattern(t *testing.T) {
	retentionPolicy := api_v1alpha1.RetentionPolicy{
		Name:     "keep-last-1",
		KeepLast: 1,
		Prune:    true,
		DryRun:   false,
	}

	testCases := []struct {
		name              string
		backupOpt         BackupOptions
		restoreOpt        RestoreOptions
		sourceFileNames   []string
		restoredFileNames []string
	}{
		{
			name: "normal backup and restore",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-1", "file-2", "file-3"},
		},
		{
			name: "exclude one file during backup",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
				Exclude:         []string{"file-1"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-2", "file-3"},
		},
		{
			name: "exclude multiple files during backup",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
				Exclude:         []string{"file-1", "file-2"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-3"},
		},
		{
			name: "include one file during restore",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
			},
			restoreOpt: RestoreOptions{
				Include: []string{"file-1"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-1"},
		},
		{
			name: "include multiple files during restore",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
			},
			restoreOpt: RestoreOptions{
				Include: []string{"file-1", "file-2"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-1", "file-2"},
		},
		{
			name: "exclude one file during restore",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
			},
			restoreOpt: RestoreOptions{
				Exclude: []string{"file-1"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-2", "file-3"},
		},
		{
			name: "exclude multiple files during restore",
			backupOpt: BackupOptions{
				RetentionPolicy: retentionPolicy,
			},
			restoreOpt: RestoreOptions{
				Exclude: []string{"file-1", "file-2"},
			},
			sourceFileNames:   []string{"file-1", "file-2", "file-3"},
			restoredFileNames: []string{"file-3"},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "stash-unit-test-")
			if err != nil {
				t.Error(err)
				return
			}

			w, err := setupTest(tempDir)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup(tempDir)

			// Initialize Repository
			err = w.InitializeRepository()
			if err != nil {
				t.Error(err)
			}

			// create the source files
			err = os.Remove(filepath.Join(targetPath, fileName))
			if err != nil {
				t.Error(err)
				return
			}
			for _, name := range test.sourceFileNames {
				err = ioutil.WriteFile(filepath.Join(targetPath, name), []byte(fileContent), 0777)
				if err != nil {
					t.Error(err)
					return
				}
			}
			test.backupOpt.BackupPaths = []string{targetPath}

			_, err = w.RunBackup(test.backupOpt, testTargetRef)
			if err != nil {
				t.Error(err)
				return
			}

			// delete target then restore
			if err = os.RemoveAll(targetPath); err != nil {
				t.Error(err)
				return
			}
			test.restoreOpt.RestorePaths = []string{targetPath}

			_, err = w.RunRestore(test.restoreOpt, testTargetRef)
			if err != nil {
				t.Error(err)
				return
			}

			var restoredFiles []string
			err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					restoredFiles = append(restoredFiles, info.Name())
				}
				return nil
			})
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.restoredFileNames, restoredFiles)
		})
	}
}

func TestApplyRetentionPolicy(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
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
	// take two backup
	_, err = w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	_, err = w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	// apply retention policy
	repoStats, err := w.ApplyRetentionPolicies(backupOpt.RetentionPolicy)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(1), repoStats.SnapshotCount)
	assert.Equal(t, int64(1), repoStats.SnapshotsRemovedOnLastCleanup)
}
func TestVerifyRepositoryIntegrity(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "stash-unit-test-")
	if err != nil {
		t.Error(err)
	}

	w, err := setupTest(tempDir)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(tempDir)

	// Initialize Repository
	err = w.InitializeRepository()
	if err != nil {
		t.Error(err)
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
	// take two backup
	_, err = w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	_, err = w.RunBackup(backupOpt, testTargetRef)
	if err != nil {
		t.Error(err)
	}
	// apply retention policy
	repoStats, err := w.VerifyRepositoryIntegrity()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, *repoStats.Integrity)
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
			Host:               "host-0",
			FileName:           filepath.Join(targetPath, fileName),
			StdoutPipeCommands: []Command{stdoutPipeCommand},
		},
		{
			Host:               "host-1",
			FileName:           filepath.Join(targetPath, fileName),
			StdoutPipeCommands: []Command{stdoutPipeCommand},
		},
		{
			Host:               "host-2",
			FileName:           filepath.Join(targetPath, fileName),
			StdoutPipeCommands: []Command{stdoutPipeCommand},
		},
	}
}
