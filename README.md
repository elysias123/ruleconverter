### ruleconverter

本项目类似subconverter 不过ruleconverter是针对分流规则的转换 而不是整个配置文件

使用例子：
http://127.0.0.1:8081/rule?target=mihomo&origin=adblock&url=https://easylist-downloads.adblockplus.org/easylistchina.txt

将adblock语法的规则转换为mihomo domain格式的分流规则


#### target(目标规则)支持类型

| 类型    | 作为源类型 | 作为目标类型 | 参数 |
| ------ | :---: | :----: | ------ |
| adblock / adguard |   ✓   |   ×   | adblock |
| hosts |   ✓   |    ×   | hosts |
| mihomo domain格式 |  ×  |  ✓   | mihomo / mihomo_domain |
| surge module 格式 |   ×   |    ✓   | surge / surge_module |
| surge ruleset 格式 |   ×   |    ✓   | surge_ruleset |

> 支持导入多个url链接(使用,隔开) 并且会进行去重

### build

```shell
make
```

编译适用于在Android设备上运行的ruleconverter:
```shell
make GOOS=android GOARCH=arm64 CC=<ndk工具链clang路径>
```

### 或者前往[release](https://github.com/elysias123/ruleconverter/releases)直接下载编译好的可执行文件运行

---

### 如何运行

在8081端口上运行:
```shell
chmod +x out/subconverter
out/subconverter --port 8081
```
在Linux等系统上建议使用systemd服务在后台持久化运行

