package parse_update_script

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
)

type UpdateScript struct {
	ScriptType           string
	ScriptOrder          string
	ScriptRequired       string
	ScriptExitMax        string
	ScriptDescription    string
	ScriptDocs           string
	ScriptExitCodeReboot string
	FilePath             string
}

type UpdateScriptFile struct {
	FilePath    string
	FileContent []byte
}

func extract_script_param(param_type string, byte_arr [][]byte) string {
	regex_string := fmt.Sprintf("^#\\s?%s\\s?:\\s?(.*)$", param_type)
	scrypt_type_regexp, _ := regexp.Compile(regex_string)
	for i := 0; i < len(byte_arr); i++ {
		line_string := string(byte_arr[i])
		groups := scrypt_type_regexp.FindAllStringSubmatch(line_string, -1)
		if len(groups) > 0 {
			return (groups[0][1])
		}
	}
	return ("")
}

func (us *UpdateScript) ParseScript(uf *UpdateScriptFile) {
	byte_arr := bytes.Split(uf.FileContent, []byte("\n"))
	us.ScriptType = extract_script_param("Script-Type", byte_arr)
	us.ScriptOrder = extract_script_param("Script-Order", byte_arr)
	us.ScriptExitMax = extract_script_param("Script-Exit-Max", byte_arr)
	us.ScriptDescription = extract_script_param("Script-Description", byte_arr)
	us.ScriptDocs = extract_script_param("Script-Docs", byte_arr)
	us.ScriptExitCodeReboot = extract_script_param("Script-Exit-Code-Reboot", byte_arr)
	us.FilePath = uf.FilePath
}

func ReadFile(update_script_file *UpdateScriptFile) error {
	file_content, err := ioutil.ReadFile(update_script_file.FilePath)
	update_script_file.FileContent = file_content
	return err
}
