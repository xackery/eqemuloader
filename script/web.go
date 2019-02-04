package script

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// WebRun represents the web command
func (s *Script) WebRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}
	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "web") {
			return
		}
		s.commandRun("docker stop web")
		s.commandRunParse("docker rm web")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "start" {
		if strings.Contains(out, "web") {
			fmt.Println("RebuildEQ website service is already running")
			return
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs web")
		return
	}
	fmt.Println("RebuildEQ website service")

	err = s.ensureImage("jonsnowd3n/rebuildeq-web:latest")
	if err != nil {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "invalid working directory")
		return
	}

	_, err = os.Stat("./bin/web/application/logs")
	if err != nil {
		err = os.Mkdir("./bin/web/application/logs", 0777)
		if err != nil {
			err = errors.Wrap(err, "failed to make bin/web/application/logs directory")
			return
		}
	}
	err = os.Chmod("./bin/web/application/logs", 0777)
                if err != nil {
                        err = errors.Wrap(err, "failed to chmod bin/web/application/logs directory")
                        return
                }
	err = os.Chmod("./bin/web/application/cache", 0777)
                if err != nil {
                        err = errors.Wrap(err, "failed to chmod bin/web/application/cache directory")
                        return
                }
	s.commandRunParse("docker rm hugo")
	fmt.Println("Updating the changelog using hugo...")
	err = s.commandRun(fmt.Sprintf("docker run --rm --name hugo -v %s/bin/web/hugo:/hugo -w /hugo -v %s/bin/web:/var/www --network=%s jonsnowd3n/rebuildeq-web:latest hugo -d /var/www/html/changelog", wd, wd, s.Docker.Network))
	if err != nil {
		err = errors.Wrap(err, "failed to update changelog using hugo")
		return
	}

	s.commandRunParse("docker rm web")
	err = s.commandRun(fmt.Sprintf("docker run --name web -v %s/bin/web:/var/www -p %s:80 -e DB_HOST=mariadb -e DB_USERNAME=%s -e DB_PASSWORD=%s -e DB_NAME=%s --network=%s --detach=true jonsnowd3n/rebuildeq-web:latest", wd, s.Web.Port, s.Database.Username, s.Database.Password, s.Database.Name, s.Docker.Network))
	if err != nil {
		return
	}
	return
}
