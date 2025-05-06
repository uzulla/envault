


FILE_NAME=".env"
NUM_VARS=5

if [ ! -z "$1" ]; then
  FILE_NAME=$1
fi

if [ ! -z "$2" ]; then
  NUM_VARS=$2
fi

echo "# テスト用の.envファイル - $(date)" > $FILE_NAME

for i in $(seq 1 $NUM_VARS); do
  VAR_NAME="TEST_VAR_$i"
  
  VALUE=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
  
  echo "$VAR_NAME=$VALUE" >> $FILE_NAME
done

echo "テスト用の.envファイルを作成しました: $FILE_NAME"
echo "環境変数の数: $NUM_VARS"
