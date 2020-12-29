
## Install MAGE

```sh
mkdir -p ~/tmp
cd ~/tmp
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
cd ..
rm -rf mage
```