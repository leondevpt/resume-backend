package main

/***
wire.go 第一行 // +build wireinject 这个 build tag 确保在常规编译时忽略wire.go 文件（因为常规编译时不会指定 wireinject 标签）
与之相对的是 wire_gen.go 中的 //+build !wireinject 两组对立的build tag保证在任意情况下,wire.go 与 wire_gen.go 只有一个文件生效
避免了方法被重复定义的编译错误
*/
