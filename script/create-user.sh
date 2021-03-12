#!bin/bash
sudo adduser wt --gecos "" --disabled-password
echo "wt:wt110999" | sudo chpasswd

sudo adduser syf --gecos "" --disabled-password
echo "syf:syf110999" | sudo chpasswd

sudo adduser byy --gecos "" --disabled-password
echo "byy:byy110999" | sudo chpasswd

sudo adduser zw --gecos "" --disabled-password
echo "zw:zw110999" | sudo chpasswd

sudo adduser xxm --gecos "" --disabled-password
echo "xxm:xxm110999" | sudo chpasswd

sudo adduser zh --gecos "" --disabled-password
echo "zh:zh110999" | sudo chpasswd

sudo adduser yt --gecos "" --disabled-password
echo "yt:yt110999" | sudo chpasswd

sudo adduser yl --gecos "" --disabled-password
echo "yl:yl110999" | sudo chpasswd

usermod -a -G sudo wt
usermod -a -G sudo syf
usermod -a -G sudo byy
usermod -a -G sudo zw
usermod -a -G sudo xxm
usermod -a -G sudo yt
usermod -a -G sudo yl
usermod -a -G sudo zh
