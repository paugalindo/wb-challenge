package bootstrap

import (
	"sync"
	"wb-challenge/bus"
	"wb-challenge/internal"
	"wb-challenge/internal/commands"
)

func InitCommandBus(vehicleRepository internal.VehicleRepository, vehicleView internal.VehicleView,
	groupRepository internal.GroupRepository, groupView internal.GroupView,
	publisher internal.EventsPublisher,
) bus.CommandBus {
	var mutex sync.Mutex
	commandBus := bus.NewCommandBus()

	loadVehiclesCmdHandler := commands.NewLoadVehiclesHandler(&mutex, vehicleRepository, groupRepository, publisher)
	commandBus.RegisterHandler(commands.LoadVehiclesType, &loadVehiclesCmdHandler)

	createGroupCmdHandler := commands.NewCreateGroupHandler(groupRepository, publisher)
	commandBus.RegisterHandler(commands.CreateGroupType, &createGroupCmdHandler)

	assignVehicleToGroupCmdHandler := commands.NewAssignVehicleToGroupHandler(&mutex, groupRepository, vehicleView, publisher)
	commandBus.RegisterHandler(commands.AssignVehicleToGroupType, &assignVehicleToGroupCmdHandler)

	occupyVehicleCmdHandler := commands.NewOccupyVehicleHandler(&mutex, groupView, vehicleRepository, publisher)
	commandBus.RegisterHandler(commands.OccupyVehicleType, &occupyVehicleCmdHandler)

	dropoffGroupCmdHandler := commands.NewDropOffGroupHandler(&mutex, groupRepository, publisher)
	commandBus.RegisterHandler(commands.DropOffGroupType, &dropoffGroupCmdHandler)

	releaseVehicleCMHandler := commands.NewReleaseVehicleHandler(&mutex, groupView, vehicleRepository, publisher)
	commandBus.RegisterHandler(commands.ReleaseVehicleType, &releaseVehicleCMHandler)

	assignVehiclesCmdHandler := commands.NewAssignVehiclesHandler(&mutex, groupRepository, vehicleView, publisher)
	commandBus.RegisterHandler(commands.AssignVehiclesType, &assignVehiclesCmdHandler)

	return commandBus
}
