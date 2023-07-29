CWD="$(cd -P -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
cd $CWD
bash ./start-client.sh & bash ./start-server.sh