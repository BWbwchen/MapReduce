{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go # go
    protobuf
    protoc-gen-go
    protoc-gen-go-grpc
  ];
}
