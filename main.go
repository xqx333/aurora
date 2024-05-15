package main

import (
    "aurora/initialize"
    "embed"
    "encoding/json"
    "io/fs"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"

    "github.com/acheong08/endless"
    "github.com/joho/godotenv"
)

//go:embed web/*
var staticFiles embed.FS

func main() {
    gin.SetMode(gin.ReleaseMode)
    router := initialize.RegisterRouter()
    subFS, err := fs.Sub(staticFiles, "web")
    if err != nil {
        log.Fatal(err)
    }
    router.StaticFS("/web", http.FS(subFS))

    // 添加 /ip 路由
    router.GET("/ip", func(c *gin.Context) {
        // 发送 HTTP GET 请求到 ipify API
        resp, err := http.Get("https://api.ipify.org?format=json")
        if err != nil {
            log.Println("无法获取出口 IP:", err)
            c.String(http.StatusInternalServerError, "无法获取出口 IP")
            return
        }
        defer resp.Body.Close()

        // 读取响应体
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("无法读取响应体:", err)
            c.String(http.StatusInternalServerError, "无法读取响应体")
            return
        }

        // 解析 JSON 响应
        var data map[string]string
        err = json.Unmarshal(body, &data)
        if err != nil {
            log.Println("无法解析 JSON:", err)
            c.String(http.StatusInternalServerError, "无法解析 JSON")
            return
        }

        // 获取出口 IP
        ip := data["ip"]

        // 将出口 IP 写入响应
        c.String(http.StatusOK, "容器的出口 IP: "+ip)
    })

    _ = godotenv.Load(".env")
    host := os.Getenv("SERVER_HOST")
    port := os.Getenv("SERVER_PORT")
    tlsCert := os.Getenv("TLS_CERT")
    tlsKey := os.Getenv("TLS_KEY")

    if host == "" {
        host = "0.0.0.0"
    }
    if port == "" {
        port = os.Getenv("PORT")
        if port == "" {
            port = "8080"
        }
    }

    if tlsCert != "" && tlsKey != "" {
        _ = endless.ListenAndServeTLS(host+":"+port, tlsCert, tlsKey, router)
    } else {
        _ = endless.ListenAndServe(host+":"+port, router)
    }
}
