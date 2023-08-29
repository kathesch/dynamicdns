import os
import subprocess

projectname = "dynamicdns"
platforms = [("linux", "arm64"), ("darwin", "arm64")]

for os_name, arch in platforms:
    output = f"{projectname}-{os_name}-{arch}"

    env_vars = os.environ.copy()
    env_vars["GOOS"] = os_name
    env_vars["GOARCH"] = arch

    result = subprocess.run(
        [
            "go",
            "build",
            "-o",
            f"bin/{output}",
            "cmd/dynamicdns.go",
        ],
        env=env_vars,
    )

    if result.returncode != 0:
        print(f"Failed to build for {os_name}/{arch}")
        exit(1)

print("Build completed for all platforms")
