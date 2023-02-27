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

import "github.com/fvbock/endless"
import "github.com/joho/godotenv"

func init() {
    log.Println("init...")
    godotenv.Overload()
}

func english(w http.ResponseWriter, r *http.Request) {
    ocr(w, r, "--oem 1 -l eng quiet")
}

func chinese(w http.ResponseWriter, r *http.Request) {
    ocr(w, r, "--oem 1 -l chi_sim quiet")
}

func ocr(w http.ResponseWriter, r *http.Request, param string) {
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
        var timeout int = 10 //超时
        if os.Getenv("timeout")!="" {
            set, err := strconv.Atoi(os.Getenv("timeout"))
            if err==nil {
                timeout=set
            }
        }


        var stdout, stderr bytes.Buffer
        command := exec.Command("/bin/bash", "-c", "/usr/local/bin/tesseract "+tmp+" "+tmp+" "+param)
        command.Stdout = &stdout
        command.Stderr = &stderr
        command.Start()

        //监听超时
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

    server  := endless.NewServer(":9090", nil)
    server.BeforeBegin = func(add string) {
        pid := syscall.Getpid()
        log.Println("pid:",pid)
        con := []byte(strconv.Itoa(pid))
        err := ioutil.WriteFile("pid", con, 0644)
        if err != nil {
            log.Fatal(err)
        }
    }

    err := server.ListenAndServe()
    if err != nil {
        log.Println(err)
    }

    log.Println("Server stopped")

    os.Exit(0)
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
