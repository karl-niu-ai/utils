package main

import (
	"fmt"
	"reflect"
)

func setupCatchalls(helper *provisionerTestHelper) {
	setupCatchallsFor(helper.mockServiceDB.EXPECT())
	setupCatchallsFor(helper.mockAWSClient.EXPECT())
	setupCatchallsFor(helper.mockKubeClient.EXPECT())
	setupCatchallsFor(helper.mockDB.EXPECT())
}

func makeHandler(receiverName string, methodName string, arity int) reflect.Value {
	switch arity {
	case 0:
		return reflect.ValueOf(func() {
			panic(fmt.Sprintf("UNEXPECTED METHOD CALL: %s.%s()", receiverName, methodName))
		})
	case 1:
		return reflect.ValueOf(func(arg interface{}) {
			panic(fmt.Sprintf("UNEXPECTED METHOD CALL: %s.%s(%+v)", receiverName, methodName, arg))
		})
	case 2:
		return reflect.ValueOf(func(arg0 interface{}, arg1 interface{}) {
			panic(
				fmt.Sprintf(
					"UNEXPECTED METHOD CALL: %s.%s(%+v, %+v)",
					receiverName,
					methodName,
					arg0,
					arg1,
				),
			)
		})
	case 3:
		return reflect.ValueOf(func(arg0 interface{}, arg1 interface{}, arg2 interface{}) {
			panic(
				fmt.Sprintf(
					"UNEXPECTED METHOD CALL: %s.%s(%+v, %+v, %+v)",
					receiverName,
					methodName,
					arg0,
					arg1,
					arg2,
				),
			)
		})
	case 4:
		return reflect.ValueOf(
			func(arg0 interface{}, arg1 interface{}, arg2 interface{}, arg3 interface{}) {
				panic(
					fmt.Sprintf(
						"UNEXPECTED METHOD CALL: %s.%s(%+v, %+v, %+v, %+v)",
						receiverName,
						methodName,
						arg0,
						arg1,
						arg2,
						arg3,
					),
				)
			},
		)
	case 6:
		return reflect.ValueOf(
			func(arg0 interface{}, arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}) {
				panic(
					fmt.Sprintf(
						"UNEXPECTED METHOD CALL: %s.%s(%+v, %+v, %+v, %+v, %+v, %+v)",
						receiverName,
						methodName,
						arg0,
						arg1,
						arg2,
						arg3,
						arg4,
						arg5,
					),
				)
			},
		)
	case 8:
		return reflect.ValueOf(
			func(arg0 interface{}, arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}, arg6 interface{}, arg7 interface{}) {
				panic(
					fmt.Sprintf(
						"UNEXPECTED METHOD CALL: %s.%s(%+v, %+v, %+v, %+v, %+v, %+v, %+v, %+v)",
						receiverName,
						methodName,
						arg0,
						arg1,
						arg2,
						arg3,
						arg4,
						arg5,
						arg6,
						arg7,
					),
				)
			},
		)
	default:
		panic(fmt.Sprintf("UNEXPECTED %d-ary METHOD: %s.%s", arity, receiverName, methodName))
	}
}

func setupCatchallsFor(any interface{}) {
	val := reflect.ValueOf(any)
	typ := val.Type()
	receiverName := val.Elem().Type().Name()
	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)

		var args []reflect.Value
		for j := 0; j < method.Type().NumIn(); j++ {
			args = append(args, reflect.ValueOf(gomock.Any()))
		}

		methodName := typ.Method(i).Name
		fmt.Printf("Calling %d-ary method '%s.%s'...\n", len(args), receiverName, methodName)
		method.
			Call(args)[0].
			MethodByName("AnyTimes").
			Call([]reflect.Value{})[0].
			MethodByName("Do").
			Call([]reflect.Value{makeHandler(receiverName, methodName, len(args))})
	}
}
