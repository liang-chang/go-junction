# go-junction


## target 配置 
target = '{useless}/Z'<br/>
指向 {useless}/Z 的文件夹

target = '\<auto\>'<br/>
自动在所有可用的target目录下挑选一个文件夹，下一个可用的文件夹


## link 配置

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