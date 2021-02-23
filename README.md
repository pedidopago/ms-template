
## Install MAGE

```sh
mkdir -p ~/temp
cd ~/temp
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
cd ..
rm -rf mage
```