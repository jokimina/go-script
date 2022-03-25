package main

import (
	"fmt"
	"log"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
	"github.com/OpenNebula/one/src/oca/go/src/goca/parameters"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/shared"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/vm/keys"
)

func ListVms(controller *goca.Controller) {
	// Get short informations of the VMs
	vms, err := controller.VMs().Info(parameters.PoolWhoAll)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(vms.VMs); i++ {
		// This Info method, per VM instance, give us detailed informations on the instance
		// Check xsd files to see the difference
		vm, err := controller.VM(vms.VMs[i].ID).Info(false)
		if err != nil {
			log.Fatal(err)
		}

		//Do some others stuffs on vm
		fmt.Printf("%+v\n", vm.Name)
	}
}

func ListTemplates(controller *goca.Controller) {
	tpls, err := controller.Templates().Info(parameters.PoolWhoAll)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(tpls.Templates); i++ {
		// tpl, err := controller.Template(tpls.Templates[i].ID).Info(false)
		tpl, err := controller.Template(tpls.Templates[i].ID).Info(false, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", tpl.Name)
	}
}

func createVm(controller *goca.Controller) {
	tplCtr := controller.Template(5)
	// tplCtr := controller.Template(17)
	tpl, err := tplCtr.Info(false, false)
	if err != nil {
		log.Fatal(err)
	}

	tpl.Template.Add(keys.Name, "test-nebula-api")
	tpl.Template.CPU(2)
	tpl.Template.VCPU(3)
	// M
	tpl.Template.Memory(2048)
	disk := tpl.Template.GetVectors(string(shared.DiskVec))[0]
	disk.AddPair(string(shared.Size), "53248")
	fmt.Println(tpl.Template.String())
	vmID, err := controller.VMs().Create(tpl.Template.String(), false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("vmID: %v\n", vmID)
}

func main() {
	client := goca.NewDefaultClient(
		goca.NewConfig("x", "x", "http://127.0.0.1:2633/RPC2"),
	)
	controller := goca.NewController(client)
	// // ListVms(controller)
	// // ListTemplates(controller)
	createVm(controller)
}
