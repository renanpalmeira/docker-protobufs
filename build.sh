docker build -t protobuf:latest $(cat .versions | sed 's@^@--build-arg @g ') .
