package libchorizo

// To run this test
// go test util_test.go util.go

import (
	"testing"
	"os/exec"
)

func TestFileExistsWithFileAbsent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	cmd := exec.Command("rm", test_filename)
	cmd.Run()
	fe := FileExists(test_filename)
	if fe != false {
		t.Error("FileExists returns true when file is not present")
	}


}
func TestFileExistsWithFilePresent(t *testing.T) {
	test_filename := "/tmp/chorizo_test_file_shouldnt_ever_be_useful"
	cmd := exec.Command("/usr/bin/touch", test_filename)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fe := FileExists(test_filename)
	if fe != true {
		t.Error("FileExists returns false when file is present")
	}

}