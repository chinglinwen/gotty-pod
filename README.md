# Gotty-pod

## Gotty-pod working with [gotty](https://github.com/chinglinwen/gotty) to provide k8s pod shell

容器shell平台，用于登录容器shell。

## 为了安全性相关考虑
- 网络功能关闭了（即不能连接数据库）
- 默认以www用户登录
- 文件系统都为只读权限
- 所有的命令都会被记录
- 权限信息和git的项目一致
- 只有git master能登录online容器