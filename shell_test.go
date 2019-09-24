package powershell

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/bhendo/go-powershell/backend"
)

type context struct {
	shell Shell
}

func TestLocalShell(t *testing.T) {
	c := context{}
	dir, err := ioutil.TempDir("", "pwsh-test")
	if err != nil {
		t.Errorf("error from TempDir: %v", err)
	}
	defer os.RemoveAll(dir)
	f, err := ioutil.TempFile(dir, "cowabunga")
	if err != nil {
		t.Errorf("error from Tempfile: %v", err)
	}
	f.Close()
	if err := os.Chdir(dir); err != nil {
		t.Errorf("error from Chdir: %v", err)
	}

	s, err := New(&backend.Local{})
	if err != nil {
		t.Errorf("error from New: %v", err)
	}
	c.shell = s

	t.Run("it can run commands", c.basicTest)
	t.Run("it correctly parses its boundary token", c.boundaryTest)
}

func (c *context) basicTest(t *testing.T) {
	sout, serr, err := c.shell.Execute(`echo "hi from stdout"`)
	if err != nil {
		t.Errorf("error from Execute: %v", err)
	}

	if sout != "hi from stdout\r\n" {
		t.Errorf("unexpected stdout: %q", sout)
	}
	if serr != "" {
		t.Errorf("unexpected stderr: %q", sout)
	}
}

func (c *context) boundaryTest(t *testing.T) {
	done := make(chan bool)
	for i := 0; i < 1000; i++ {
		timeout := time.After(10 * time.Second)

		go func() {
			_, _, err := c.shell.Execute("ls")
			if err != nil {
				t.Errorf("error from Execute: %v", err)
			}
			done <- true
		}()

		select {
		case <-timeout:
			t.Fatalf("Timed out waiting for command after %d iterations", i)
		case <-done:
			continue
		}
	}
}
