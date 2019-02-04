package script

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (s *Script) setup() (err error) {
	var value string
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the EQEMU Loader program. Review each step during this setup wizard")
	fmt.Println("prod: Production Server")
	fmt.Println("dev: Development Server")
	for {
		fmt.Print("What is the plan with this server? (prod, dev): ")
		value, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		value = strings.Replace(value, "\n", "", -1)
		switch value {
		case "prod":
			fmt.Println("prod is not yet available with this tool")
			continue
		case "dev":
		default:
			fmt.Println("invalid choice, please choose dev or prod")
			continue
		}
		s.Global.Stage = value
		break
	}
	fmt.Println("Server Long Name is used to identify your server on the server selection screen.")
	for {
		fmt.Print("Enter a server long name (e.g. Shin's Server): ")
		value, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		value = strings.Replace(value, "\n", "", -1)
		if len(value) < 3 {
			fmt.Println("name must be at least 3 characters long")
			continue
		}
		s.World.LongName = value
		break
	}
	fmt.Println("Server Short Name is used for storing ui and character profiles.")
	fmt.Println("It should be small with no spaces.")
	for {
		fmt.Print("Enter a server short name: ")
		value, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		value = strings.Replace(value, "\n", "", -1)
		if len(value) < 3 {
			fmt.Println("name must be at least 3 characters long")
			continue
		}
		s.World.ShortName = value
		break
	}
	fmt.Println("A database password is used to access your database.")
	for {
		fmt.Print("Enter a database password: ")
		value, err = r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		value = strings.Replace(value, "\n", "", -1)
		if len(value) < 3 {
			fmt.Println("password must be at least 3 characters long")
			continue
		}
		s.Database.Password = value
		break
	}

        fmt.Println("A database root password is used to administrate your database.")
        for {
                fmt.Print("Enter a database root password: ")
                value, err = r.ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        continue
                }
                value = strings.Replace(value, "\n", "", -1)
                if len(value) < 3 {
                        fmt.Println("root password must be at least 3 characters long")
                        continue
                }
                s.Database.RootPassword = value
                break
        }



	fmt.Println("loader.conf file generated.")
	fmt.Println("It is advised to review this this file before running this program again")

	return
}
