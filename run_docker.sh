#! /bin/bash
set -e
#RootPath=$(cd $(dirname $0)/..; pwd)
#RootPath=$($(dirname $0); pwd)
RootPath=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
RootPath="$RootPath/conf"
GOPATH=/go
export DiskPath="$RootPath/conf/docker/docker_data"

MIN_DNDISK_AVAIL_SIZE_GB=10

help() {
    cat <<EOF

Usage: ./run_docker.sh [ -h | --help ] [ -d | --disk </disk/path> ] [ -l | --ltptest ]
    -h, --help              show help info
    -d, --disk </disk/path>     set CubeFS DataNode local disk path
    -b, --build             build CubeFS server and client
    -s, --server            start CubeFS servers docker image
    -c, --client            start CubeFS client docker image
    -m, --monitor           start monitor web ui
    -l, --ltptest           run ltp test
    -r, --run               run servers, client and monitor
    --clean                 cleanup old docker image
EOF
    exit 0
}


clean() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml down
}

# unit test
run_unit_test() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml run unit_test
}
# build
build() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml run build
}


# start server
start_servers() {
    isDiskAvailable $DiskPath
    mkdir -p ${DiskPath}/disk/{1..4}
    docker-compose -f ${RootPath}/docker/docker-compose.yml up -d servers
}

start_client() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml run client bash -c "/cfs/script/start_client.sh ; /bin/bash"
}



start_scenariotest() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml run client
}

start_ltptest() {
    docker-compose -f ${RootPath}/docker/docker-compose.yml run client bash -c "/cfs/script/start.sh -ltp"
}

run_scenariotest() {
    build
    start_servers
    start_scenariotest
    clean
}

run_ltptest() {
    build
    start_servers
    start_ltptest
    clean
}

run() {
    start_servers
    start_client
}

cmd="help"

ARGS=( "$@" )
for opt in ${ARGS[*]} ; do
    case "$opt" in
        -h|--help)
            help
            ;;
        -b|--build)
            cmd=build
            ;;
        -t|--test)
            cmd=run_test
            ;;
        -l|--ltptest)
            cmd=run_ltptest
            ;;
        -n|--scenariotest)
            cmd=run_scenariotest
            ;;
        -r|--run)
            cmd=run
            ;;
        -s|--server)
            cmd=run_servers
            ;;
        -c|--client)
            cmd=run_client
            ;;
        -m|--monitor)
            cmd=run_monitor
            ;;
        -clean|--clean)
            cmd=clean
            ;;
        *)
            ;;
    esac
done

function isDiskAvailable() {
    Disk=${1:-"need diskpath"}
    [[ -d $Disk ]] || mkdir -p $Disk
    if [[ ! -d $Disk ]] ; then
        echo "error: $DiskPath must be exist and at least 10GB free size"
        exit 1
    fi
    avail_kb=$(df -k $Disk | tail -1 | awk '{print $4}')
    avail_GB=$(( $avail_kb / 1000 / 1000 ))
    if (( $avail_GB < $MIN_DNDISK_AVAIL_SIZE_GB )) ; then
        echo "$Disk: avaible size $avail_GB GB < Min Disk avaible size $MIN_DNDISK_AVAIL_SIZE_GB GB" ;
        exit 1
    fi
}

for opt in ${ARGS[*]} ; do
    case "-$1" in
        --d|---disk)
            shift
            export DiskPath=${1:?"need disk dir path"}
            isDiskAvailable $DiskPath
            shift
            ;;
        -)
            break
            ;;
        *)
            shift
            ;;
    esac
done

case "-$cmd" in
    -help) help ;;
    -run) run ;;
    -build) build ;;
    -run_servers) start_servers ;;
    -run_client) start_client ;;
    -run_monitor) start_monitor ;;
    -run_ltptest) run_ltptest ;;
    -run_test) run_unit_test ;;
    -run_scenariotest) run_scenariotest ;;
    -clean) clean ;;
    *) help ;;
esac

