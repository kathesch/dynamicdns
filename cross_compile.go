package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	platforms := []string{"darwin", "linux"}
	archs := []string{"arm64"}

	err := os.Mkdir("bin", 0755)
	if err != nil {
		fmt.Println("Error", err)
	}

	for _, platform := range platforms {
		for _, arch := range archs {
			outputName := fmt.Sprintf("dynamicdns-%s-%s", platform, arch)
			cmd := exec.Command(
				"go", "build", "-o", "bin/"+outputName, "-ldflags", "-s -w", "cmd/dynamicdns.go",
			)
			cmd.Env = append(os.Environ(), "GOOS="+platform, "GOARCH"+arch)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error building for %s/%s: %v\n", platform, arch, err)
				return
			}
		}
	}
}
