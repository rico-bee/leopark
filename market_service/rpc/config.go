package rpc

type Config struct {
	// Secret keys
	// WARNING! These defaults are insecure, and should be changed for deployment
	secret          string
	aesKey          string
	batchPrivateKey string
}

var RpcConfig = &Config{
	// Secret keys
	// WARNING! These defaults are insecure, and should be changed for deployment
	secret:          "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",                             // any string
	aesKey:          "ffffffffffffffffffffffffffffffff",                                 // 32 character hex string
	batchPrivateKey: "1111111111111111111111111111111111111111111111111111111111111111", // 64 character hex string
}
