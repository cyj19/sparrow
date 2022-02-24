/**
 * @Author: cyj19
 * @Date: 2022/2/24 14:47
 */

package server

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// 服务的方法
type methodType struct {
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
	numCall   int
}

type service struct {
	name      string                 // 服务名称
	refVal    reflect.Value          // 服务实例
	refType   reflect.Type           // 服务类型
	methodMap map[string]*methodType // 服务方法
}

func newService(v interface{}, serviceName string, useName bool) (*service, error) {
	s := &service{
		refVal:  reflect.ValueOf(v),
		refType: reflect.TypeOf(v),
	}

	sName := reflect.Indirect(s.refVal).Type().Name()
	// 判断服务是不是公开的
	if !ast.IsExported(sName) {
		return nil, errors.New(fmt.Sprintf("rpc server: %s is not valid service name ", sName))
	}

	if useName {
		if serviceName == "" {
			return nil, errors.New("service name is null")
		}
		sName = serviceName
	}
	s.name = sName
	methodMap, err := registerMethods(s.refType)
	if err != nil {
		return nil, err
	}
	s.methodMap = methodMap

	return s, nil
}

func registerMethods(refType reflect.Type) (map[string]*methodType, error) {
	methodMap := make(map[string]*methodType)
	for i := 0; i < refType.NumMethod(); i++ {
		method := refType.Method(i)
		mType := method.Type
		mName := method.Name
		if !ast.IsExported(mName) {
			return nil, errors.New(fmt.Sprintf("method %s is not public", mName))
		}
		// 校验函数的输入参数是否符合规则 func(*server.MethodTest, *arg, *reply)
		if mType.NumIn() != 3 {
			continue
		}
		// 检验输入参数，必须是指针类型
		argType := mType.In(1)
		replyType := mType.In(2)
		if argType.Kind() != reflect.Ptr || replyType.Kind() != reflect.Ptr {
			continue
		}
		// 校验函数的返回参数，必须是error
		if mType.NumOut() != 1 || mType.Out(0) != typeOfError {
			continue
		}

		methodMap[mName] = &methodType{
			method:    method,
			argType:   argType,
			replyType: replyType,
		}

	}

	if len(methodMap) == 0 {
		return nil, errors.New("the service does not provide a public method")
	}

	return methodMap, nil
}
