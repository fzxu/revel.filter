package filter

import (
	"code.google.com/p/go.net/websocket"
	"github.com/robfig/revel"
	"reflect"
)

var (
	controllerFilters map[reflect.Type][]*RegisteredMethod = make(map[reflect.Type][]*RegisteredMethod)
)

type RegisteredMethod struct {
	When         revel.When
	TargetMethod interface{}
	Methods      []string //registered methods
}

func AddControllerFilter(target interface{}, when revel.When, methods ...string) {
	receiverType := reflect.TypeOf(target).In(0) // the receiver type is actually the controller static type
	controllerFilters[receiverType] = append(controllerFilters[receiverType], &RegisteredMethod{When: when, TargetMethod: target, Methods: methods})
}

func ControllerFilter(c *revel.Controller, fc []revel.Filter) {

	// Collect the values for the method's arguments.
	var methodArgs []reflect.Value
	methodArgs = append(methodArgs, reflect.ValueOf(c.AppController).Elem()) // The receiver of the filter function

	// Bind the funciton signature
	for _, arg := range c.MethodType.Args {
		var boundArg reflect.Value
		// Ignore websocket for now
		if arg.Type != reflect.TypeOf((*websocket.Conn)(nil)) {
			boundArg = revel.Bind(c.Params, arg.Name, arg.Type)
		}

		methodArgs = append(methodArgs, boundArg)
	}

	var resultValue reflect.Value

	// Call before
	for _, registeredMethod := range controllerFilters[c.Type.Type] {

		if registeredMethod.When == revel.BEFORE {
			for _, method := range registeredMethod.Methods {
				if method == c.MethodName {
					targetMethod := reflect.ValueOf(registeredMethod.TargetMethod)
					resultValue = targetMethod.Call(methodArgs)[0]
				}
			}
		}
	}

	// The filter chain only continue when the result Value is nil
	// Call after
	if !resultValue.IsValid() || resultValue.IsNil() {
		fc[0](c, fc[1:])

		for _, registeredMethod := range controllerFilters[c.Type.Type] {
			if registeredMethod.When == revel.AFTER {
				for _, method := range registeredMethod.Methods {
					if method == c.MethodName {
						targetMethod := reflect.ValueOf(registeredMethod.TargetMethod)
						resultValue = targetMethod.Call(methodArgs)[0]
					}
				}
			}
		}
	}

	if resultValue.Kind() == reflect.Interface && !resultValue.IsNil() {
		c.Result = resultValue.Interface().(revel.Result)
	}

}
