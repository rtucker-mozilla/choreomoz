package libchorizo

// To run this test
// go test util_test.go util.go

import (
	"testing"
	"os/exec"
)

func TouchFile(file_path string){
	cmd := exec.Command("/usr/bin/touch", file_path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func RmFile(file_path string){
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
