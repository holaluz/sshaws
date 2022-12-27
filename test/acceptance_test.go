package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestSshAws(t *testing.T) {
	t.Run("sshaws compiles", func(t *testing.T) {
		binDir := "/tmp/sshawstest"
		binPath := fmt.Sprintf("%s/sshaws", binDir)

		os.Setenv("GOARCH", "amd64")

		os.Setenv("GOOS", "windows")
		buildWindows := exec.Command("go", "build", "-o", binPath+".windows", "../cmd/sshaws/main.go")
		out, err := buildWindows.CombinedOutput()
		if err != nil {
			t.Error(string(out))
		}

		os.Setenv("GOOS", "darwin")
		buildMacOS := exec.Command("go", "build", "-o", binPath+".darwin", "../cmd/sshaws/main.go")
		out, err = buildMacOS.CombinedOutput()
		if err != nil {
			t.Error(string(out))
		}

		os.Setenv("GOOS", "linux")
		buildLinux := exec.Command("go", "build", "-o", binPath+".linux", "../cmd/sshaws/main.go")
		out, err = buildLinux.CombinedOutput()
		if err != nil {
			t.Error(string(out))
		}

		err = os.RemoveAll(binDir)
		if err != nil {
			t.Error(err)
		}
	})
}
