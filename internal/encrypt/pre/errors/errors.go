package errors

var DesKeyError = "key size should be 8"
var DesIvError = "IV size should be 8"
var TripleDesKeyError  = "key size should be 24"
var AesKeyError = "key size should be 16, 24 or 32"
var AesIvError = "IV size should be 16, 24 or 32"
var RsatransError = "error occur when trans to *rsa.Publickey"
//var RsaNilError = "error occur when decrypt"
var EcckeyError = "key size should be 224 256 384 521"