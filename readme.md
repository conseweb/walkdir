
# 把指定目录(包括子目录)生成目录树哈希写入到指定文件里面.

# 编译：

~~~
$ make  # 编译二进制
~~~

# 调用 :
在dthash工具项目所在根目录执行
    sudo ./build/dthash --input=yourdirectory --output=yourresultfile
    input是所要处理(生成目录树哈希)的目录
    output是生成目录树的结果文件

# 结果文件格式 :
结果文件中每一行代表一个处理过的文件,每行中的内容用逗号隔开
每行字段含义按顺序分别是:
文件名,文件的sha1哈希值,文件大小(byte)

#测试
     1.使用docker统一测试环境
     2.利用docker/Dockerfile生成测试环境镜像
     3.sh docker/container/dthash.sh 生成测试环境容器
     4.在测试环境容器中，$GOPATH下找到walkdir工程进行测试


