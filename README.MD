# shuZhiNet
熟知网手机版后端

提供的API见[api.md](api.md)。

## 需要自行安装的依赖

- redis

## 安装go依赖包
会安装gin（用于实现文件修改时的热重载）和项目中用到的包。
```bash
./script/install-dependency.sh
```

## 启动开发服务器
会设置必要的环境变量，并启动redis和服务器
```bash
./script/dev-server.sh
```
