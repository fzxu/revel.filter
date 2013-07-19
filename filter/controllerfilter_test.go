package filter

import (
	"fmt"
	"github.com/robfig/revel"
	"net/url"
	"reflect"
	"testing"
)

type TestController struct {
	*revel.Controller
}

func (t TestController) Show(id string, x int) revel.Result {
	fmt.Println("show sth")
	return t.Redirect("/")
}

func (t TestController) Edit(id string, x int) revel.Result {
	fmt.Println("edit sth")
	return t.Redirect("/")
}

func (t TestController) isOwner(id string, x int) revel.Result {
	fmt.Println("before, really works:" + id)
	fmt.Println(x)
	return nil //if it does not return nil, e.g. return t.Redirect , then it will stop the chain and redirect directly
}

func (t TestController) callAfter(id string, x int) revel.Result {
	fmt.Println("after, really works:" + id)
	fmt.Println(x)
	return nil // AFTER's return value does not really useful, unless nothing returned in the Action
}

func TestAddControllerFilter(t *testing.T) {
	AddControllerFilter(TestController.isOwner, revel.BEFORE, "Show", "Edit")
	AddControllerFilter(TestController.callAfter, revel.AFTER, "Show")

	c := revel.NewController(nil, nil)
	var controller TestController
	c.AppController = controller
	c.MethodName = "Show"
	c.Params = &revel.Params{Values: make(url.Values)}
	c.Params.Set("id", "cool")
	c.Params.Set("x", "5")

	methodArg := []*revel.MethodArg{
		&revel.MethodArg{Name: "id", Type: reflect.TypeOf("")},
		&revel.MethodArg{Name: "x", Type: reflect.TypeOf(1)}}
	c.MethodType = &revel.MethodType{Name: "Show", Args: methodArg}

	flt := func(c *revel.Controller, fc []revel.Filter) {}

	fc := []revel.Filter{flt}
	fmt.Println("Going to call Show:")
	ControllerFilter(c, fc)

	//test Edit method

	c.MethodName = "Edit"
	c.MethodType = &revel.MethodType{Name: "Edit", Args: methodArg}
	fmt.Println("Going to call Edit:")
	ControllerFilter(c, fc)
}
