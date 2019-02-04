package script

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *Script) crash(name string) (err error) {
	var lastFile os.FileInfo

	err = filepath.Walk(fmt.Sprintf("%scrash/", s.Bin.Directory), func(path string, f os.FileInfo, err error) error {
		if !strings.Contains(f.Name(), name) {
			return nil
		}

		if lastFile != nil && lastFile.ModTime().Before(f.ModTime()) {
			return nil
		}
		lastFile = f
		return nil
	})
	if lastFile == nil {
		fmt.Println("no crash dumps detected for", name)
		return
	}
	if s.IsVerbose {
		fmt.Println("latest crash is", lastFile.Name())
	}
	commands := strings.Split(fmt.Sprintf("docker run -v %s:/src --network=%s eqemu/server sh -c", s.Bin.Directory, s.Docker.Network), " ")
	commands = append(commands, fmt.Sprintf(`gdb --batch --quiet -ex "thread apply all bt" -ex "quit" /src/world /src/crash/%s | gzip -9 > /src/crash.gz`, lastFile.Name()))

	err = s.commandRunSplit(commands)
	if err != nil {
		err = nil
		return
	}
	fmt.Println("crash.gz saved at", s.Bin.Directory)
	return
}
