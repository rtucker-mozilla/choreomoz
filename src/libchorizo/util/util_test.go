package libchorizo

// To run this test
// go test util_test.go util.go

import (
	"os/exec"
	"testing"
)

func TouchFile(file_path string) {
	cmd := exec.Command("/usr/bin/touch", file_path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func RmFile(file_path string) {
	cmd := exec.Command("rm", "-f", file_path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func TestFileExistsWithFileAbsent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	RmFile(test_filename)
	fe := FileExists(test_filename)
	if fe != false {
		t.Error("FileExists returns true when file is not present")
	}
}

func TestFileExistsWithFilePresent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	fe := FileExists(test_filename)
	if fe != true {
		t.Error("FileExists returns false when file is present")
	}

}

func TestHasScriptPathWithFileAbsent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	RmFile(test_filename)
	fe := HasScriptPath(test_filename)
	if fe != false {
		t.Error("HasSCriptFile returns true when file is not present")
	}
}

func TestHasScriptPathWithFilePresent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	fe := HasScriptPath(test_filename)
	if fe != true {
		t.Error("HasSCriptFile returns false when file is present")
	}
}

func TestHasLockFileWithFileAbsent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	RmFile(test_filename)
	fe := HasLockFile(test_filename)
	if fe != false {
		t.Error("HasLockFile returns true when file is not present")
	}
}

func TestHasLockFileWithFilePresent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	fe := HasLockFile(test_filename)
	if fe != true {
		t.Error("HasLockFile returns false when file is present")
	}
}

func TestHasStateFileWithFileAbsent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	RmFile(test_filename)
	fe := HasStateFile(test_filename)
	if fe != false {
		t.Error("HasStateFile returns true when file is not present")
	}
}

func TestHasStateFileWithFilePresent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	fe := HasStateFile(test_filename)
	if fe != true {
		t.Error("HasStateFile returns false when file is present")
	}
}

func TestWriteLockFile(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	RmFile(test_filename)
	fe := FileExists(test_filename)
	if fe != false {
		t.Error("lockfile is present but should be absent")
	}
	WriteLockFile(test_filename)
	fe = FileExists(test_filename)
	if fe != true {
		t.Error("lockfile is absent but should be present")
	}
	RmFile(test_filename)
}

func TestDeleteLockFile(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	DeleteLockFile(test_filename)
	fe := FileExists(test_filename)
	if fe != false {
		t.Error("lockfile is present but should be absent")
	}
	RmFile(test_filename)
}

func TestDeleteStateFile(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	TouchFile(test_filename)
	DeleteStateFile(test_filename)
	fe := FileExists(test_filename)
	if fe != false {
		t.Error("statefile is present but should be absent")
	}
	RmFile(test_filename)
}

func TestScriptValidatorStructConstructorFilepath(t *testing.T) {
	filepath := "/tmp/chorizo_test_script_file"
	TouchFile(filepath)
	t_sv := NewScriptValidator(filepath)
	if t_sv.Filepath != filepath {
		t.Error("Filepath in constructor not being set correctly")
	}

}

func TestScriptValidatorStructConstructorFilemode(t *testing.T) {
	filepath := "/tmp/chorizo_test_script_file"
	TouchFile(filepath)
	t_sv := NewScriptValidator(filepath)
	if t_sv.Filemode != "-rw-r--r--" {
		t.Error("Filepath in constructor not being set correctly: ", t_sv.Filemode, "asfd")
	}

}
func TestScriptValidatorStructConstructorIsExecutable(t *testing.T) {
	filepath := "/tmp/chorizo_test_script_file"
	TouchFile(filepath)
	t_sv := NewScriptValidator(filepath)
	if t_sv.IsExecutable != false {
		t.Error("IsExecutable not being set correctly: ", t_sv.Filemode, "asfd")
	}
	cmd := exec.Command("chmod", "700", filepath)
	cmd.Run()

	t_sv = NewScriptValidator(filepath)
	if t_sv.IsExecutable != true {
		t.Error("IsExecutable not being set correctly with executable set: ", t_sv.Filemode)
	}
	RmFile(filepath)
}

func TestScriptValidatorGetValidFilenameWithInvalidFilename(t *testing.T) {
	filepath := "/tmp/chorizo_test_script_file"
	TouchFile(filepath)
	t_sv := NewScriptValidator(filepath)
	if t_sv.ValidFileName != false {
		t.Error("GetValidFilename not being set correctly:", t_sv.ValidFileName)
	}
	RmFile(filepath)
}
func TestScriptValidatorGetValidFilenameWithValidFilename(t *testing.T) {
	filepath := "/tmp/1_chorizo_test_script_file"
	TouchFile(filepath)
	t_sv := NewScriptValidator(filepath)
	if t_sv.ValidFileName != true {
		t.Error("GetValidFilename not being set correctly:", t_sv.ValidFileName)
	}
	RmFile(filepath)
}
