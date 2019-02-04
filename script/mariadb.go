package script

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	//used for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// MariaDBRun represents a database
func (s *Script) MariaDBRun(args []string) (err error) {
	err = s.prepare()
	if err != nil {
		return
	}

	_, err = os.Stat(s.Database.Directory)
	if err != nil {
		err = os.Mkdir(s.Database.Directory, 0777)
		if err != nil {
			err = errors.Wrap(err, "failed to make database directory")
			return
		}
		err = os.Chmod(s.Database.Directory, 0777)
		if err != nil {
			err = errors.Wrap(err, "failed to chmod database directory")
			return
		}

	}

	out, err := s.commandRunParse("docker ps")
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "pull" {
		err = s.mariadbReleaseDownload()
		if err != nil {
			return
		}
		err = s.binPull()
		if err != nil {
			return
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "stop" {
		if !strings.Contains(out, "mariadb") {
			return
		}
		s.commandRun("docker stop mariadb")
		s.commandRunParse("docker rm mariadb")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		s.commandRun("docker logs mariadb")
		return
	}
	fmt.Println("MariaDB")

	if len(os.Args) > 1 && os.Args[1] == "dump" {
		if len(args) < 1 {
			err = fmt.Errorf("must provide destination")
			return
		}
		err = s.SQLDump(args[0])
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "inject" {
		if len(args) < 1 {
			err = fmt.Errorf("must provide source")
			return
		}

		err = s.SQLInject(args[0])
		return
	}
	if s.isPortListening("0.0.0.0:3306") {
		return
	}
	s.commandRunParse("docker rm mariadb")
	err = s.commandRun(fmt.Sprintf(`docker run -v %s:/bitnami/mariadb -p 3306:3306 -e MARIADB_DATABASE=%s -e MARIADB_USER=%s -e MARIADB_PASSWORD=%s -e MARIADB_ROOT_PASSWORD=%s -e MARIADB_REPLICATION_MODE=master -e ALLOW_EMPTY_PASSWORD=no --network=%s --detach=true --name=mariadb bitnami/mariadb:latest`, s.Database.Directory, s.Database.Name, s.Database.Username, s.Database.Password, s.Database.RootPassword, s.Docker.Network))
	if err != nil {
		return
	}
	for counter := 0; counter < 15; counter++ {
		<-time.After(1 * time.Second)
		if s.isPortListening("0.0.0.0:3306") {
			return
		}
	}
	err = fmt.Errorf("timeout waiting for mariadb to be up")
	return
}

// SQLInject will inject a .sql file to the database
func (s *Script) SQLInject(src string) (err error) {
	if !s.isPortListening("0.0.0.0:3306") {
		err = fmt.Errorf("please start mariadb before injecting")
		return
	}
	fi, err := os.Stat(src)
	if err != nil {
		err = errors.Wrapf(err, "could not find %s", src)
		return
	}
	if !fi.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file (%q)", fi.Name(), fi.Mode().String())
		return
	}

	switch filepath.Ext(src) {
	case ".sql":
		err = fileGZip(src, "./.cache/tmp.gz")
		if err != nil {
			err = errors.Wrapf(err, "failed to copy/compress %s to cache", fi.Name())
			return
		}
	case ".gz":
		err = fileCopy(src, "./.cache/tmp.gz")
		if err != nil {
			err = errors.Wrapf(err, "failed to copy %s to cache", fi.Name())
			return
		}
	default:
		err = fmt.Errorf("%s is an invalid sql file, must end with .sql or .gz", src)
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get working directory")
		return
	}
	src = fmt.Sprintf("%s/.cache/tmp.gz", wd)

	dir, err := filepath.Abs(filepath.Dir(src))
	if err != nil {
		err = errors.Wrap(err, "failed to get absolute path")
		return
	}
	cmds := strings.Split(fmt.Sprintf("docker run -v %s:/inject --entrypoint  --network=%s bitnami/mariadb:latest sh -c", dir, s.Docker.Network), " ")
	cmds = append(cmds, fmt.Sprintf("gunzip < /inject/tmp.gz | mysql -h mariadb -u root -p%s %s", s.Database.RootPassword, s.Database.Name))
	fmt.Printf("Injecting mariadb... (This may take a while)\n")
	err = s.commandRunSplit(cmds)
	if err != nil {
		err = errors.Wrapf(err, "database inject failed.")
		return
	}
	return
}

// SQLDump dumps the database to a file
func (s *Script) SQLDump(dst string) (err error) {
	if !s.isPortListening("0.0.0.0:3306") {
		err = s.MariaDBRun([]string{"start"})
	}
	fmt.Printf("Dumping mariadb to %s... (This may take a while)\n", dst)

	f, err := os.Create(dst)
	if err != nil {
		err = errors.Wrapf(err, "could not write to %s", dst)
		return
	}
	f.Close()
	err = os.Remove(dst)
	if err != nil {
		err = errors.Wrapf(err, "could not write test %s", dst)
		return
	}

	dir, err := filepath.Abs(filepath.Dir(dst))
	if err != nil {
		err = errors.Wrap(err, "failed to get absolute path")
		return
	}

	err = s.commandRun(fmt.Sprintf("docker run -v %s:/dump --entrypoint  --network=%s bitnami/mariadb:latest mysqldump -h mariadb -u %s -p%s %s --result-file=/dump/%s", dir, s.Docker.Network, s.Database.Username, s.Database.Password, s.Database.Name, f.Name()))
	if err != nil {
		return
	}

	return
}

// SQLTableExists returns true when it does
func (s *Script) SQLTableExists(table string) (ok bool, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", s.Database.Username, s.Database.Password, s.Database.Name))
	if err != nil {
		err = errors.Wrap(err, "failed to open mysql")
		return
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SHOW TABLES LIKE '%s'", table))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		return
	}
	if rows.Next() {
		ok = true
		return
	}
	return
}

func (s *Script) mariadbReleaseDownload() (err error) {
	switch strings.ToLower(s.Bin.ReleaseType) {
	case "github":
		err = s.githubReleaseDownload()
		if err != nil {
			return
		}
	case "gitea":
		err = s.giteaReleaseDownload("database", s.Database.ReleaseURL, s.Database.ReleaseUser, s.Database.ReleaseRepo, s.Bin.ReleaseAccessToken, s.Database.Directory)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknown database release type: %s", s.Database.ReleaseType)
	}
	return
}
