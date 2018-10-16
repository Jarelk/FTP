package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"trackio"

	"github.com/jlaffaye/ftp"
)

func main() {
	var (
		user = flag.String("user", "admin", "Username for login")
		pass = flag.String("pass", "123456", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "localhost", "IP-adress")
	)
	fmt.Println("Connecting to " + *host + ":" + fmt.Sprint(*port))
	client, err := ftp.Dial(*host + ":" + fmt.Sprint(*port))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Authenticating")
	if err := client.Login(*user, *pass); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Attempting to upload")
	//DownloadAndPrint(client)
	UploadFile(client)
	client.Quit()
}

func DownloadAndPrint(client *ftp.ServerConn) {
	entries, _ := client.List("hello.txt")

	for _, entry := range entries {
		name := entry.Name
		reader, err := client.Retr(name)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Reading downloaded file %v\n", name)
		b := make([]byte, 128)
		var downloaded uint64 = 0
		for {
			n, err := reader.Read(b)
			downloaded += uint64(n)
			fmt.Printf("\r%v%%", (100*downloaded)/entry.Size)
			//fmt.Printf("%q", b[:n])
			if err == io.EOF {
				break
			}
		}
		closeErr := reader.Close()
		if closeErr != nil {
			fmt.Println(closeErr)
			panic(closeErr)
		}

	}
}

func UploadFile(client *ftp.ServerConn) {
	ex, erros := os.Executable()
	if erros != nil {
		fmt.Println(erros)
		panic(erros)
	}

	path := path.Join(path.Dir(ex), "/files/hello.txt")
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fileinfo, err2 := f.Stat()
	if err2 != nil {
		fmt.Println(err2)
		panic(err2)
	}

	tr := trackio.NewReader(f)

	fmt.Println("Starting upload")

	//tick := 0.0
	//t := time.Now()
	go func() {
		for {
			prc := percent(float64(tr.N()), float64(fileinfo.Size()))
			fmt.Printf("\r%v%% Uploaded", int64(prc))
			if prc == 100.0 {
				break
			}
		}
	}()

	client.Stor("/hello.txt", tr)
}

//Returns % of n / size
func percent(n, size float64) float64 {
	if n == 0 {
		return 0
	}
	if n >= size {
		return 100
	}
	return 100.0 * (n / size)
}
