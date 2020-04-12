### 使用goroutine的多图片上传包

#### 实现的功能
* [x] 多图片上传
* [x] 限制图片大小
* [x] 限制图片数量
* [x] 限制上传图片类型
* [x] 按年月日切割图片
* [x] 使用协程并发操作图片

<br>

#### 运行
```shell script
go get -v github.com/lujiahaoo/gin-upload
```

<br>

#### 配置
假设项目目录结构如下(缩略了一些)
```
project/
├── config
│   ├── app.yml
│   └── refreshtoken.txt
├── go.mod
├── go.sum
├── main.go
├── static
│   ├── 2020-04-12
│   │   ├── 7wyabtlgzs.jpg
│   │   ├── n514brfa9d.jpg
│   │   └── xlfs4ldl5z.jpg
│   └── thumbnail
│       └── 7wyabtlgzs.jpg
```
其中 `static`目录可以按自己需要更改

