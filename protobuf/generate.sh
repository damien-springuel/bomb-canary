SCRIPT_PATH=$(dirname "$(realpath -s "$0")")
PROJECT=github.com/damien-springuel/bomb-canary
OUT_PACKAGE=server/generated
OUT=$SCRIPT_PATH/../$OUT_PACKAGE
rm -rf $OUT
mkdir $OUT
protoc \
  --proto_path=$SCRIPT_PATH \
  --go_out=$OUT \
  --go-grpc_out=$OUT \
  --go_opt=module=$PROJECT/$OUT_PACKAGE \
  --go-grpc_opt=module=$PROJECT/$OUT_PACKAGE \
  $SCRIPT_PATH/api.proto
