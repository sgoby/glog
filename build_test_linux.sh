# export GOARCH=amd64
# export GOOS=linux
GOPATH=/mnt/hgfs/golang
echo $GOPATH
go test -c -i -o /mnt/hgfs/golang/src/demo.com/glog/Test_log demo.com/glog
