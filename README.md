

docker build -t proxy:v0.1 .

docker run --rm --name proxy -p 5333:53 proxy:v0.1

dig -p 5333 +tcp +short @localhost google.com
