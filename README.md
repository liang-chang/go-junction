# go-junction

'bcilfw@{autoTarget}'


旧:
link 配置
'bcilfw@{UserHome}/.dubbo',
b:备份原文件夹
c:清空备份的文件夹
i:会在target文件下创建个随机指向的文件夹
l:当link路径中包有正则式时,创建整个路径的文件夹,即使最后一级文件夹不存在时
f:强制创建
w:抑制警告


新：
link 配置
'bilfw@{UserHome}/.dubbo',
b:备份原文件夹
c:清空备份的文件夹
i:隔离，会在target文件下创建个随机指向的文件夹
l:当最后一级文件夹不存在时，创建文件夹
f:强制创建全路径的文件夹，当不存在该文件夹时

w:去掉，放在全局设置