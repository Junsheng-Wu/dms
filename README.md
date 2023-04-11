## 简介
dms是一项负责管理非k8s环境中prometheus容器规则的服务,解决原生prometheus在非k8s环境中无法使用prometheus operator去管理告警规则的问题。

## 制作容器镜像
```make docker-build```

## 制作二进制文件
```make dms-server```
