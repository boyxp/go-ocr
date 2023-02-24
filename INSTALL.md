
### 环境要求
centos 7

### 安装步骤
#
#### 1、依赖安装
```
yum install -y gcc gcc-c++ autoconf automake libtool libjpeg libpng libtiff zlib libjpeg-devel libpng-devel libtiff-devel zlib-devel
```
#### 2、gcc升级
```
yum install -y centos-release-scl

yum install -y devtoolset-8-gcc*

mv /usr/bin/gcc /usr/bin/gcc-4.8.5

ln -s /opt/rh/devtoolset-8/root/bin/gcc /usr/bin/gcc

mv /usr/bin/g++ /usr/bin/g++-4.8.5

ln -s /opt/rh/devtoolset-8/root/bin/g++ /usr/bin/g++


gcc --version
g++ --version
```
到这一步版本可以看到是 g++ (GCC) 8.3.1 20190311 (Red Hat 8.3.1-3)

#### 3、安装leptonica
```
wget http://www.leptonica.org/source/leptonica-1.82.0.tar.gz

tar zxf leptonica-1.82.0.tar.gz 

cd leptonica-1.82.0/

./configure

make

make install
```

#### 4.配置环境
```
vim /etc/profile
```
添加以下
```
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
export LIBLEPT_HEADERSDIR=/usr/local/include
export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
```
加载配置
```
source /etc/profile
```

#### 5、编译安装tesseract 5
github下载最新压缩包zip格式
https://github.com/tesseract-ocr/tesseract

unzip解压后进入目录开始编译(通常目录名为tesseract-main)
```
cd tesseract-main

./autogen.sh

./configure --with-extra-includes=/usr/local/include --with-extra-libraries=/usr/local/lib

make

make install
```
....漫长的等待....

#### 6、安装语言包
https://github.com/tesseract-ocr/tessdata_fast
下载语言包:eng.traineddata、chi_sim.traineddata
放到语言文件夹
```
/usr/local/share/tessdata/
```


#### 7、测试识别效果
找个带字图片input.png
```
tesseract input.png output -l chi_sim --oem 1
```
结果输出到output.txt

参考原文：https://www.jianshu.com/p/edfabeaf6ba8



