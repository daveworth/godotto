package floatingips

import (
	"fmt"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/context"
)

func applyAction(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := actionSvc{
		ctx: ctx,
		svc: client.FloatingIPs().Actions(),
	}
	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"assign", svc.assign},
		{"unassign", svc.unassign},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type actionSvc struct {
	ctx context.Context
	svc floatingips.ActionClient
}

func (svc *actionSvc) assign(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	ip := argFloatingIP(vm, all.Argument(0))
	dropletID := argDropletID(vm, all.Argument(1))
	err := svc.svc.Assign(svc.ctx, ip, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) unassign(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argFloatingIP(vm, all.Argument(0))
	err := svc.svc.Unassign(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}
