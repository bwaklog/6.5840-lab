cd main
go build -buildmode=plugin ../mrapps/wc.go
rm mr-out-*

if [ "$1" == "-s" ]; then
    go run mrsequential.go wc.so pg-*.txt
elif [ "$1" == "-c" ]; then
    go run mrcoordinator.go pg-*.txt
else
    echo "Usage: ./test.sh [-s|-c]"
fi

cat mr-out-* | sort | more
