package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/franela/goreq"
)

func httpGetString(url string) (body string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(data)), err
}

func saveFolder() string {
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(strings.Split(gopath, ";")[0], "bin")
	}
	return "."
}

func getBinary(u url.URL, binName string) {
	log.Println("Request:", u.String())
	res, err := goreq.Request{
		Method:      "POST",
		Uri:         u.String(),
		Compression: goreq.Gzip(),
	}.Do()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Println("StatusCode:", res.StatusCode)
		io.Copy(os.Stdout, res.Body)
		return
	}
	defer res.Body.Close()
	saveDir := saveFolder()
	savePath := filepath.Join(saveDir, binName)
	log.Println("Save to", savePath)
	filefd, err := os.Create(savePath)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(filefd, res.Body)
	filefd.Close()
	os.Chmod(savePath, 0755)
	return
}

func figureBinName(repo string) string {
	base := filepath.Base(repo)
	if runtime.GOOS == "windows" {
		base += ".exe"
	}
	return base
}

func installAction(ctx *cli.Context) {
	server := ctx.GlobalString("server")
	ref := ctx.String("ref")
	if len(ctx.Args()) == 0 {
		log.Fatal("Need args")
	}
	repo := ctx.Args().Get(0)
	log.Println(repo)

	u := url.URL{
		Scheme: "http",
		Host:   server,
	}

	// guess repo
	if !strings.Contains(repo, "/") {
		luckyURL := u
		luckyURL.Path = "/lucky/" + repo
		var err error
		repo, err = httpGetString(luckyURL.String())
		if err != nil {
			log.Fatal(err)
		}
		if repo == "" {
			log.Fatal("Guess repo failed, require fullname <owner/reponame>")
		}
	}
	log.Println("use repo:", repo)
	binURL := u
	binURL.Path = fmt.Sprintf("%s/%s/%s/%s", repo, ref, runtime.GOOS, runtime.GOARCH)
	getBinary(binURL, figureBinName(repo))
}

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Name = "grcli"
	app.Usage = "gobuild3 client"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "grget server address",
			Value:  "grget.shengxiang.me",
			EnvVar: "GRGET_SERVER_ADDR",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "install",
			Usage: "install builded file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ref",
					Value: "master",
				},
			},
			Action: installAction,
		},
	}
}

func main() {
	app.RunAndExitOnError()
}
