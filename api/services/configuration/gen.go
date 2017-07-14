package configuration

//go:generate protoc -I.:../../../vendor:../../../vendor/github.com/gogo/protobuf --gogo_out=plugins=grpc,import_path=github.com/ehazlett/interlock/api/services/configuration,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:. configuration.proto
