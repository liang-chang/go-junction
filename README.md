# go-junction
windows 下创建 目录链接(junction)的工具

## 使用方法
````
例出所有配置及当前应用的文件夹
go-junction.exe --config=config.toml --action=list

检查当前配置
go-junction.exe --config=config.toml --action=check

按配置文件应用配置
go-junction.exe --config=config.toml --action=make

恢复
go-junction.exe --config=config.toml --action=recovery

````

## 示例配置
````
[config]
#link文件夹重命名备份，默认为false
backupLinkFolder = false

#清空备份的文件夹，默认为true
clearBackupFolder = true

#当target文件不存在时，默认为true
createTargetFolder = true

#当没有匹配到link文件夹时是否警告，默认为true
warnLinkNotExist = true

#所有的target文件
targetFolders=[
'v:/useless/Z/Z[0-9]',
]

[pathAlias]
#build in path variable
# UserHome
# Temp
useless='V:/useless/'
chromeCache='V:/chrome_cache'
firefoxCache='V:/firefox_cache'
appDataRoaming='{UserHome}/AppData/Roaming'
appDataLocal='{UserHome}/AppData/Local'
QQDataHome='{UserHome}/AppData/Roaming/Tencent'

[[symbolic]]
target = '{useless}/Z'
link = [
'fil@v:/|log.|/tt',
]

[[symbolic]]
target = '<auto>'
link = [
'bcilf@v:/|log.|/tt',
]

````
## targetFolders 配置解释

````
#多个数值支持通配符
targetFolders=[
'v:/useless/?',
'c:/[a-b]*'
]
````
<pre>
pattern:
    { term }
term:
    '*'                                  匹配0或多个非路径分隔符的字符
    '?'                                  匹配1个非路径分隔符的字符
    '[' [ '^' ] { character-range } ']'  字符组（必须非空）
    c                                    匹配字符c（c != '*', '?', '\\', '['）
    '\\' c                               匹配字符c
character-range:
    c           匹配字符c（c != '\\', '-', ']'）
    '\\' c      匹配字符c
    lo '-' hi   匹配区间[lo, hi]内的字符
    
</pre>
详细见 filepath.match 的单元测试结果

[https://golang.org/src/path/filepath/match_test.go](https://golang.org/src/path/filepath/match_test.go)


## pathAlias
内置两个变量 UserHome 和 Temp 

win7 下<br>
UserHome=C:\Users\yourName\

支持变量嵌套，即<br>
a='c:/'<br>
b='{a}/c'<br>
也是可以的

## target 配置 
target = '{useless}/Z'<br/>
指向 {useless}/Z 的文件夹

target = '\<auto\>'<br/>
自动在所有可用的target目录下挑选一个文件夹，下一个可用的文件夹


## link 配置

'bcilf@v:/|log.|/tt',

@之扣的路径 v:/|log.|/tt , |log.| 表示这一块是个正则式，竖线中间的是表达式，
这里表示的是匹配以log开头的文个

旧:
link 配置<br/>
'bcilfw@{UserHome}/.dubbo',<br/>
b:备份原文件夹<br/>
c:清空备份的文件夹<br/>
i:会在target文件下创建个随机指向的文件夹<br/>
l:当link路径中包有正则式时,创建整个路径的文件夹,即使最后一级文件夹不存在时<br/>
f:强制创建<br/>
w:抑制警告<br/>


新：
link 配置<br/>
'bcilf@{UserHome}/.dubbo',<br/>
b:备份原文件夹<br/>
c:清空备份的文件夹<br/>
i:隔离，会在target文件下创建个随机指向的文件夹<br/>
l:当最后一级文件夹不存在时，创建文件夹<br/>
f:强制创建全路径的文件夹，当不存在该文件夹时<br/>
~~w:抑制警告~~