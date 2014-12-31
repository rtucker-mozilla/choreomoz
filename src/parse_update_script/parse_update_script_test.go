package parse_update_script

import (
  "testing"
  "io/ioutil"
  "os"
  "fmt"
)
var script_headers_string = `
#!/bin/bash
### SCRIPT START ###
# Script-Type: update
# Script-Exit-Code-Reboot: 99
# Script-Order: 10
# Script-Exit-Max: 128
# Script-Description: Test Description
# Script-Docs: http://docs.org
echo "asdf"`
var script_headers = []byte(script_headers_string)

func TestUpdateScriptStruct_ScriptType(t *testing.T) {
  fmt.Println("Starting tests of parse_update_script")
  var us UpdateScript
  us.ScriptType = "update"
  if us.ScriptType != "update"{
    t.Error("Unable to set ScriptType")
  }
}
func TestUpdateScriptStruct_ScriptOrder(t *testing.T) {
  var us UpdateScript
  us.ScriptOrder = "10"
  if us.ScriptOrder != "10"{
    t.Error("Unable to set ScriptOrder")
  }
}
func TestUpdateScriptStruct_ScriptRequired(t *testing.T) {
  var us UpdateScript
  us.ScriptRequired = "True"
  if us.ScriptRequired != "True"{
    t.Error("Unable to set ScriptRequired")
  }
}
func TestUpdateScriptStruct_ScriptExitMax(t *testing.T) {
  var us UpdateScript
  us.ScriptExitMax = "100"
  if us.ScriptExitMax != "100"{
    t.Error("Unable to set ScriptExitMax")
  }
}
func TestUpdateScriptStruct_ScriptDescription(t *testing.T) {
  var us UpdateScript
  us.ScriptDescription = "Test Description Here"
  if us.ScriptDescription != "Test Description Here"{
    t.Error("Unable to set ScriptDescription")
  }
}
func TestUpdateScriptStruct_ScriptDocs(t *testing.T) {
  var us UpdateScript
  us.ScriptDocs = "Test Docs Here"
  if us.ScriptDocs != "Test Docs Here"{
    t.Error("Unable to set ScriptDocs")
  }
}
func TestUpdateScriptStruct_ScriptExitCodeReboot(t *testing.T) {
  var us UpdateScript
  us.ScriptExitCodeReboot = "99"
  if us.ScriptExitCodeReboot != "99"{
    t.Error("Unable to set ScriptExitCodeReboot")
  }
}
func TestCanReadFile(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  filecontent := []byte("Test\nFile\nContent")
  err := ioutil.WriteFile(filepath, filecontent, 0644)
  if err != nil {
    t.Error("Could not write file")
  }
  var uf UpdateScriptFile
  uf.FilePath = filepath
  read_err := ReadFile(&uf)
  if read_err != nil {
    t.Error("Cannot read file contents")
  }
  if string(uf.FileContent) != string(filecontent) {
    t.Error("File content doesn't match")
  }
  os.Remove(filepath)
}
func TestCanParseScriptTypeWithNoSpaceBetweenHash(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  var no_space_script_headers = []byte("#Script-Type: update\n")
  ioutil.WriteFile(filepath, no_space_script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile_err := ReadFile(&uf)
  if ReadFile_err != nil {
    panic(ReadFile_err)
  }
  us.ParseScript(&uf)
  if us.ScriptType != "update"{
    t.Error("Unable to parse ScriptType")
  }
  os.Remove(filepath)

}
func TestCanParseScriptType(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile_err := ReadFile(&uf)
  if ReadFile_err != nil {
    panic(ReadFile_err)
  }
  us.ParseScript(&uf)
  if us.ScriptType != "update"{
    t.Error("Unable to parse ScriptType")
  }
  os.Remove(filepath)

}
func TestCanParseScriptOrder(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile(&uf)
  //fmt.Println(string(uf.FileContent))
  us.ParseScript(&uf)
  if us.ScriptOrder != "10"{
    t.Error("Unable to parse ScriptOrder")
  }
  os.Remove(filepath)

}
func TestCanParseScriptExitMax(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile(&uf)
  //fmt.Println(string(uf.FileContent))
  us.ParseScript(&uf)
  if us.ScriptExitMax != "128"{
    t.Error("Unable to parse ScriptExitMax")
  }
  os.Remove(filepath)

}
func TestCanParseScriptDescription(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile(&uf)
  //fmt.Println(string(uf.FileContent))
  us.ParseScript(&uf)
  if us.ScriptDescription != "Test Description"{
    t.Error("Unable to parse ScriptDescription")
  }
  os.Remove(filepath)

}
func TestCanParseScriptDocs(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile(&uf)
  //fmt.Println(string(uf.FileContent))
  us.ParseScript(&uf)
  if us.ScriptDocs != "http://docs.org"{
    t.Error("Unable to parse ScriptDocs")
  }
  os.Remove(filepath)

}
func TestCanParseScriptExitCodeReboot(t *testing.T) {
  filepath := "/tmp/testfile.txt"
  ioutil.WriteFile(filepath, script_headers, 0644)
  var us UpdateScript
  var uf UpdateScriptFile
  uf.FilePath = filepath
  ReadFile(&uf)
  //fmt.Println(string(uf.FileContent))
  us.ParseScript(&uf)
  if us.ScriptExitCodeReboot != "99"{
    fmt.Println(us.ScriptExitCodeReboot)
    t.Error("Unable to parse ScriptExitCodeReboot")
  }
  os.Remove(filepath)

}
