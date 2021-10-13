package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalSessions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "norminet_total_sessions",
		Help: "Number of scans sessions initiated",
	})
	totalScans = promauto.NewCounter(prometheus.CounterOpts{
		Name: "norminet_total_scans",
		Help: "Number of files scanned",
	})
	totalUsers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "norminet_total_users",
		Help: "Number of different users",
	})
)

func normScan(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return w
	}

	params := r.URL.Query()
	fileMap, fileOk := params["filename"]
	userMap, userOk := params["username"]
	if !fileOk || len(fileMap) <= 0 || !userOk || len(userMap) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.ErrMissingFile.Error()))
		return w
	}

	remoteName := userMap[0] + "@" + strings.Split(r.RemoteAddr, ":")[0]
	fileName := fmt.Sprintf("./temp/%s/%s", remoteName, fileMap[0])

	checkError(createDirIfNotExist("./temp/" + remoteName))
	err = ioutil.WriteFile(fileName, data, 0677)
	if err != nil {
		logError(err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return w
	}

	output, _ := exec.Command("norminette", fileName).Output()
	_, _ = w.Write(output)
	_ = os.Remove(fileName)
	return w
}

func normHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		totalScans.Inc()
		w = normScan(w, r)
	} else if r.Method == "GET" {
		totalSessions.Inc()
		version, err := exec.Command("norminette", "-v").Output()
		if err != nil {
			logError(err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write(version)
	}
}

func main() {
	checkError(createDirIfNotExist("./temp"))
	http.HandleFunc("/norm", normHandler)
	http.Handle("/metrics", promhttp.Handler())
	checkError(http.ListenAndServe(":8080", nil))
}
