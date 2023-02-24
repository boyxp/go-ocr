
# 安装tesseract4

yum -y install yum-utils

yum-config-manager --add-repo https://download.opensuse.org/repositories/home:/Alexander_Pozdnyakov/CentOS_7/
rpm --import https://build.opensuse.org/projects/home:Alexander_Pozdnyakov/public_key

yum --showduplicates list tesseract

yum install tesseract


下载中文训练数据

https://github.com/tesseract-ocr/tessdata

chi_sim.traineddata

chi_sim_vert.traineddata

放到 /usr/share/tesseract/4/tessdata/

