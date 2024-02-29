package relations

import (
	"context"
	"fmt"
	"log"
	"os"

	mgrpc "github.com/rgrente/mtp-golibs/grpc"
	pb "github.com/rgrente/mtp-golibs/pb"
)

func CheckRelationExist(permissionParams *pb.CreatePermissionParams) (bool, error) {

	log.Println("Function CheckPermission() -> Dialing : ", fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call CheckPermission Route from permission Microservice
	permission, err := grpc_client.CheckPermission(ctx, permissionParams)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return permission.IsAllowed, nil
}
