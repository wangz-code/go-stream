# go-stream
替换nginx 的stream模块转发, 主要是nginx编译模块不太方便



## 缘起
 ~假如你有一台国外VPS + 一台国内的轻量云,  
1. 当你直接ping vps 延迟和丢包惨不忍睹,  当你用轻量云ping vps 结果还过得去
这个时候你就会想, 如果通过轻量云转发一下流量 由: 我==>vps  变成  我==> 轻量云 ==> vps

你可能会觉得这不就是多此一举, 而且带宽也被限制到轻量云, 这里面的还是有区别的, 一般云厂商都有自己的国外云主机或者自己的国外专线 避免家庭带宽转来转去导致丢包,  其次云主机出海外很正常,多少应该有点作用. 不然天天被墙云厂商还怎么活🤕, 

用了好几天小流量跑起来还是很舒服的比直连感觉上快很多, 大流量还是走直连吧


## nginx 需要手动编译 stream 模块

参考 [使用国内服务器中转流量](https://v2xtls.org/%E4%BD%BF%E7%94%A8%E5%9B%BD%E5%86%85%E6%9C%8D%E5%8A%A1%E5%99%A8%E4%B8%AD%E8%BD%AC%E6%B5%81%E9%87%8F/)
```sh

# 转发流
stream {
    # vps增加的配置
    server {
        listen 1024;  # 1-65535的任意一个数字，无需与境外服务器的端口号相同
        proxy_pass 1.1.1.1:443; # 用境外ip和端口号替换
    }
}
```
## 有firewalld 和 Nginx  为什么还要go
xray-core 是用golang 编写的,可能是师出同门会好一点吧(心理作用) 😏😏😏
我主要是方便迁移, 轻量云到期了 直接一起打包, 脚本一把梭哈~  不用单独去记编译nginx 还是 有firewalld

## 编译至 linux
``` CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-stream -trimpath -ldflags "-s -w -buildid=" main.go ```

