package script

import (
	"fmt"
	"os"
	"strings"
)

// PEQEditorRun represents a PEQ Editor
func (s *Script) PEQEditorRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}
	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "peqeditor") {
			return
		}
		s.commandRun("docker stop peqeditor")
		s.commandRunParse("docker rm peqeditor")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "start" {
		if strings.Contains(out, "peqeditor") {
			fmt.Println("PEQ Editor is already running")
			return
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs peqeditor")
		return
	}
	fmt.Println("PEQ Editor")

	err = s.ensureImage("eqemu/peqeditor:latest")
	if err != nil {
		return
	}

	ok, err := s.SQLTableExists("peq_admin")
	if err != nil {
		return
	}
	if !ok {
		err = s.fileDownload("https://raw.githubusercontent.com/ProjectEQ/peqphpeditor/master/sql/schema.sql", "./.cache/schema.sql")
		err = s.SQLInject("./.cache/schema.sql")
		if err != nil {
			return
		}
	}
	s.commandRunParse("docker rm peqeditor")
	err = s.commandRun(fmt.Sprintf("docker run --name peqeditor -p %s:80 -e DB_HOST=mariadb -e DB_USERNAME=%s -e DB_PASSWORD=%s -e DB_NAME=%s --network=%s --detach=true eqemu/peqeditor:latest", s.PeqEditor.Port, s.Database.Username, s.Database.Password, s.Database.Name, s.Docker.Network))
	if err != nil {
		return
	}
	return
}
