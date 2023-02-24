package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "time"
    "strconv"
    "os/exec"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

func ocr(w http.ResponseWriter, r *http.Request) {
    contentType := r.Header.Get("Content-Type")

    if len(contentType)>=19 && contentType[0:19]=="multipart/form-data" {
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

        out, err := exec.Command("/bin/bash", "-c", "/usr/bin/tesseract "+tmp+" "+tmp+" -l chi_sim").Output()
        if err != nil {
            log.Println("文件识别失败："+tmp+"\t"+err.Error())
            res(w, 3, err.Error(), string(out))
            return
        }

        data, err := ioutil.ReadFile(tmp+".txt")
        if err!=nil {
            res(w, 4, err.Error(), "")
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
    http.HandleFunc("/", ocr)
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

