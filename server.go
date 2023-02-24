package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "time"
    "bytes"
    "syscall"
    "strings"
    "strconv"
    "os/exec"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

func english(w http.ResponseWriter, r *http.Request) {
    ocr(w, r, "")
}

func chinese(w http.ResponseWriter, r *http.Request) {
    ocr(w, r, "-l chi_sim")
}

func ocr(w http.ResponseWriter, r *http.Request, lang string) {
    contentType := r.Header.Get("Content-Type")

    if len(contentType)>=19 && contentType[0:19]=="multipart/form-data" {

        //文件上传==================================
        r.ParseMultipartForm(32 << 20)
        file, _, err := r.FormFile("image")
        if err != nil {
            res(w, 1, "image文件参数不存在", "")
            return
        }
        defer file.Close()

        tmp    := "/tmp/ocr_img_"+strconv.FormatInt(time.Now().UnixNano(), 10)
        f, err := os.Create(tmp)
        if err != nil {
            res(w, 2, err.Error(), "")
            return
        }
        defer f.Close()
        io.Copy(f, file)

        log.Println("文件上传："+tmp)





        //文件识别==================================
        var timeout int = 30
        var stdout , stderr bytes.Buffer
        command := exec.Command("/bin/bash", "-c", "/usr/bin/tesseract "+tmp+" "+tmp+" "+lang)
        //command := exec.Command("/bin/bash", "-c", "sleep 10")
        command.Stdout = &stdout
        command.Stderr = &stderr
        command.Start()
        done := make(chan error)
        go func() { done <- command.Wait() }()
        after := time.After(time.Duration(timeout) * time.Second)
        select {
            case <-after:
                        command.Process.Signal(syscall.SIGINT)
                        time.Sleep(time.Second)
                        command.Process.Kill()
                        log.Println("识别超时：", tmp)
                        res(w, 3, "识别超时", "")
                        return
            case <-done:

        }
        logout := trimOutput(stdout)
        logerr := trimOutput(stderr)

        log.Println("stdout:", logout)
        log.Println("stderr:", logerr)
/*
        if len(logerr)>0 {
            log.Println("识别失败：", tmp, logerr)
            res(w, 4, "识别失败", "")
            return
        }
*/








        //输出结果==================================
        data, err := ioutil.ReadFile(tmp+".txt")
        if err!=nil {
            res(w, 5, err.Error(), "")
            return
        }

        log.Println("文件识别成功："+tmp)

        res(w, 0, "", string(data))
    } else {
        w.WriteHeader(http.StatusForbidden)
        fmt.Fprintf(w, "Forbidden")
    }
}

func main() {
    http.HandleFunc("/en", english)
    http.HandleFunc("/cn", chinese)
    http.HandleFunc("/ping", ping)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func res(w http.ResponseWriter, code int, error string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json, err := json.Marshal(data)
    if err != nil {
        fmt.Fprintf(w, "{\"code\":-1,\"message\":\"%v\",\"response\":\"\"}", err)
        return
    }

    fmt.Fprintf(w, "{\"code\":%d,\"message\":\"%v\",\"response\":%v}", code, error, string(json))
}

func ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func trimOutput(buffer bytes.Buffer) string {
    return strings.TrimSpace(string(bytes.TrimRight(buffer.Bytes(), "\x00")))
}
