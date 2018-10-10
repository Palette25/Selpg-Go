## Selpg-Go作业说明
作业介绍：selpg是实用的Linux命令行程序，相关详细介绍：[何为Selpg?](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

作业要求：参考C语言版命令行程序Selpg，采用Go语言完成Selpg-Go程序

## Selpg-Go实现
* selpg-go接受输入文件，通过命令行参数控制，实现文本文件的固定范围终端打印，向打印机传输数据打印等宫功能。

* selpg-go需要的终端输入格式如下：
`selpg-go -sstart_page -eendpage [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]`

* 其中，`-s`后接打印页面范围的起始页面下标(此处下标从1开始), `-e`后接打印的最后一页下标,。`-f`和`-l`是互斥的参数, 其中前者后不接任何数值，表示将文本文件中的换页符作为换页标记;后者接每页的行数，规划每一页的范围。'-d'接指定的打印机名称, `in_filename`用于指定输入文件名。

## Selpg-Go测试
1. 测试程序第一步测试换行符分页模式，分别检查：正常运行，输入起始页码大于总页数，结束页码大于总页数。

2. 第二步检测换页符分页模式。

3. 第三步检测打印机功能，分别检查：打印机不存在时错误打印，打印机存在时负责打印的子进程的标准输入是否正确。