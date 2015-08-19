package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	srvPort = flag.Int("p", 4000, "Listen port")
	CWD, _  = os.Getwd()
)

func Md5str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func BuildHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println(params)
	repoURL := fmt.Sprintf("https://github.com/%s/%s", params["owner"], params["repo"])
	goPath := filepath.Join("tmp/",
		"tmp-repo-"+Md5str(repoURL)+"-"+params["ref"])
	folder := filepath.Join(goPath,
		"src", "github.com", params["owner"], params["repo"])
	os.RemoveAll(goPath) // need clean first

	// clone folder
	c := exec.Command("git",
		"clone", "--depth", "5", "--branch", params["ref"],
		repoURL, folder)
	c.Dir = "."
	c.Env = append(c.Env, "GOPATH="+filepath.Join(CWD, "tmp"))
	c.Env = append(c.Env, "GOROOT="+os.Getenv("GOROOT"))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	target := filepath.Join(CWD, "tmp/dist", params["repo"])
	c = exec.Command("godep", "go", "build", "-o", target)
	c.Env = append(c.Env,
		"PATH="+os.Getenv("PATH"),
		"GOPATH="+filepath.Join(CWD, goPath),
		"GOROOT="+os.Getenv("GOROOT"),
		"GOOS="+params["goos"],
		"GOARCH="+params["arch"])
	//log.Println(c.Env)
	c.Dir = folder
	output := bytes.NewBuffer(nil)
	c.Stdout = io.MultiWriter(output, os.Stdout)
	c.Stderr = io.MultiWriter(output, os.Stderr)
	if err := c.Run(); err != nil {
		http.Error(w, output.String(), 502)
		return
	}
	http.ServeFile(w, r, target)
}

func main() {
	flag.Parse()
	m := mux.NewRouter()
	m.HandleFunc("/{owner}/{repo}/{ref}/{goos}/{arch}", BuildHandler)
	log.Printf("Listening on *:%d", *srvPort)
	http.ListenAndServe(":"+strconv.Itoa(*srvPort), m)
}
