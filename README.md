# 基于 tesseract 5的go文字识别服务

环境要求：Centos 7

## 安装
[安装 tesseract 5](INSTALL.md)

## 启动
```
sh manage.sh start
```

## 测试
```
curl 127.0.0.1:9090/cn -F "image=@test.png"
```

## 使用
可以解析域名，以POST方式发送image图片文件，得到返回的json结果如下

```
{
    "code": 0,
    "message": "",
    "response": "哈喽"
}
```
code=0为成功，其他为失败


