# sfs-go
Secure file sharing system based on blockchain 基于区块链的安全文件共享系统

# 开发计划
1. 文件加密 AES+椭圆曲线  ✔
2. IPFS搭建-文件上传和下载封装 ✔
3. Fabric搭建和复习（远程服务器）✔
4. 代理重加密
5. 文件签名
6. 地址生成 ✔
7. Fabric智能合约的编写
    - 地址管理 ✔
    - 文件管理 ✔
    - 文件签名验证 
    - 代理重加密 ✔
8. Fabric-SDK-Go集成
   - 正式Fabric环境的搭建 ✔
   - Fabric-SDK-Go集成 ✔
   - 合约调用 ✔

# 编译

```shell
go run -o sfs main.go
```

# 功能：

- 初始化
```shell
./sfs init
```

- 上传
```shell
./sfs upload -p 文件路径
```

- 展示
```shell
./sfs show 
```

- 分享
```shell
./sfs share -i  文件ID  -a 文件接收方地址
```

- 下载
```shell
./sfs download -i 文件ID  -d 文件存储路径
```
