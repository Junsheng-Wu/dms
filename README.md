## 简介
dms是一项负责管理非k8s环境中prometheus容器规则的服务,解决原生prometheus在非k8s环境中无法使用prometheus operator去管理告警规则的问题，
可以通过调用接口来实现告警规则的增删改查功能。

## 制作容器镜像
```make docker-build```

## 制作二进制文件
```make build```
