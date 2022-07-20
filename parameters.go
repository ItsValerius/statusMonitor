package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetParams(name string) map[string]string {
	functions := map[string]string{"dclp": "get_dclp", "smtp": "get_smpt_server", "port": "get_port", "receiver": "get_receiver_email"}
	params := make(map[string]string)
	for key, element := range functions {
		formatedInput := fmt.Sprintf("import %s; print (%s.%s())", name, name, element)
		cmd := exec.Command("python3", "-c", formatedInput)
		fmt.Println(cmd.Args)
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
