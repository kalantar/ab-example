
FROM node:19-buster
WORKDIR /usr/app

# # Install protoc
# RUN apt update && apt install -y protobuf-compiler

# Get source (including requirements)
COPY samples/abn-sample/frontend/node .
COPY samples/abn-sample/frontend/abn.proto .

RUN npm install -g grpc-tools
RUN grpc_tools_node_protoc --js_out=import_style=commonjs,binary:. --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` abn.proto

# Install requirements
RUN npm install

# # Generate interface to Iter8A/B(/n) service
# RUN protoc --js_out=import_style=commonjs,binary:. abn.proto

CMD ["node", "main.js"]