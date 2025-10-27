### ruleconverter

本项目类似subconverter 不过ruleconverter是针对分流规则的转换 而不是整个配置文件

使用例子：
http://127.0.0.1:8081/rule?target=mihomo&origin=adblock&url=https://easylist-downloads.adblockplus.org/easylistchina.txt

将adblock语法的规则转换为mihomo domain格式的分流规则

target(目标规则)支持的类型:

- mihomo / mihomo_domain
- mihomo_mrs
- surge / surge_module
- surge_ruleset

origin(来源规则)支持的参数:
- adblock / adguard
- hosts

### build

```shell
make
```

编译适用于在Android设备上运行的ruleconverter:
```shell
make GOOS=android GOARCH=arm64 CC=<ndk工具链clang路径>
```

在8081端口上运行:
```shell
chmod +x out/subconverter
out/subconverter --port 8081
```