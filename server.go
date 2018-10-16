package main

import(
	"fmt"
	"github.com/goftp/server"
	"github.com/goftp/file-driver"
	"os"
	"path"
	"log"
	"flag"
)

type MyNotYetSecureAuth struct{
	Name	string
	Pass 	string
}


func main() {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		user = flag.String("user", "admin", "Username for login")
		pass = flag.String("pass", "123456", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "localhost", "Port")
	)

	factory := &filedriver.FileDriverFactory{
		RootPath: path.Join(path.Dir(ex) + "/files/"),
		Perm: server.NewSimplePerm("root", "root"),
	}
	  
	opts := &server.ServerOpts{
		Factory: 	factory,
		Port: 		*port,
		Hostname:	*host,
		Auth:		&MyNotYetSecureAuth{Name: *user, Pass: *pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	
	server  := server.NewServer(opts)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func (auth *MyNotYetSecureAuth) CheckPasswd(userhash, passhash string) (bool, error){
	if userhash != auth.Name || passhash != auth.Pass{
		return false, nil
	}
	return true, nil
}