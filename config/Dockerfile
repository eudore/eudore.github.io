FROM scratch
ADD main /
CMD ["/main"]

# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go
#  -ldflags "-s -w" 