package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sebest/xff"
)

var (
	srvPort int
	gitHost string
	CWD, _  = os.Getwd()
)

func init() {
	flag.StringVar(&gitHost, "githost", "github.com", "git host prefix")
	flag.IntVar(&srvPort, "p", 4000, "Listen port")
}

func Md5str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func BuildHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println(params)
	scheme := "http"
	if gitHost == "github.com" {
		scheme = "https"
	}
	repoURL := fmt.Sprintf(scheme+"://"+gitHost+"/%s/%s.git", params["owner"], params["repo"])
	goPath := filepath.Join("tmp/",
		"tmp-repo-"+Md5str(repoURL)+"-"+params["ref"])
	folder := filepath.Join(goPath,
		"src", gitHost, params["owner"], params["repo"])
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
		http.Error(w, err.Error()+"\n"+output.String(), 502)
		return
	}
	http.ServeFile(w, r, target)
	exec.Command("scripts/add-repo.sh", params["owner"]+"/"+params["repo"]).Run()
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func ScriptHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("grins.sh"))
	tmpl.Execute(w, r.Host)
}

func LuckyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	out, _ := exec.Command("scripts/get-repo.sh", vars["name"]).Output()
	io.WriteString(w, string(out))
}

func main() {
	flag.Parse()

	m := mux.NewRouter()
	m.HandleFunc("/", Homepage)
	m.HandleFunc("/lucky/{name}", LuckyHandler)
	m.Handle("/grins.sh", xff.Handler(http.HandlerFunc(ScriptHandler)))
	m.Handle("/{owner}/{repo}/{ref}/{goos}/{arch}", xff.Handler(Gzip(http.HandlerFunc(BuildHandler))))

	log.Printf("Listening on *:%d", srvPort)
	http.ListenAndServe(":"+strconv.Itoa(srvPort), m)
}
