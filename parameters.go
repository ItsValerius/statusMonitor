package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetParams(name string) map[string]string {
	functions := map[string]string{"dclp": "get_dclp", "smtp": "get_smpt_server", "port": "get_port", "receiver": "get_receiver_email"}
	path, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	workingDir := filepath.Dir(path)
	fmt.Println(workingDir)
	params := make(map[string]string)
	for key, element := range functions {
		formatedInput := fmt.Sprintf("import %s; print (%s.%s())", name, name, element)
		cmd := exec.Command("python3", "-c", formatedInput)
		cmd.Dir = workingDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		cleanedOutput := CleanOutput(out)
		params[key] = cleanedOutput

	}
	return params
}

func CleanOutput(out []byte) string {
	outString := string(out)
	outString = strings.TrimSuffix(outString, "\n\r")
	outString = strings.TrimSuffix(outString, "\r\n")
	outString = strings.TrimSuffix(outString, "\n")
	outString = strings.TrimPrefix(outString, "['")
	outString = strings.TrimSuffix(outString, "']")
	outString = strings.TrimSpace(outString)
	return outString
}
