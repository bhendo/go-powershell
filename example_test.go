package powershell_test

import (
	"fmt"

	ps "github.com/simonjanss/go-powershell"
	"github.com/simonjanss/go-powershell/backend"
	"github.com/simonjanss/go-powershell/middleware"
)

func ExampleShell() {
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		panic(err)
	}

	// prepare remote session configuration
	config := middleware.NewSessionConfig()
	config.ComputerName = "remote-pc-1"

	// create a new shell by wrapping the existing one in the session middleware
	session, err := middleware.NewSession(shell, config)
	if err != nil {
		panic(err)
	}
	defer session.Exit()  // will also close the underlying ps shell!
	defer session.Close() // will disconnect the current session

	// everything run via the session is run on the remote machine
	stdout, _, err := session.Execute("Get-WmiObject -Class Win32_Processor")
	if err != nil {
		panic(err)
	}

	fmt.Println(stdout)
}
