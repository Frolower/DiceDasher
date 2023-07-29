CWD="$(cd -P -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
screen -dmS DiceDasher $CWD/start-all.sh 