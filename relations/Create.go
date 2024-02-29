package relations

import (
	"context"
	"fmt"
	"log"
	"os"

	mgrpc "github.com/rgrente/mtp-golibs/grpc"
	me "github.com/rgrente/mtp-golibs/merror"
	pb "github.com/rgrente/mtp-golibs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CreateRelationWithSubjectSet(params *pb.CreatePermissionWithSubjectSetParams) error {
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		return me.ProcessError(&me.MError{
			Code:           422,
			Message:        "Failed Dependency",
			Description:    "retry in a short moment or contact the developer team",
			DevCode:        424,
			DevMessage:     "Failed Dependency",
			DevDescription: "service Project not reachable",
			Trace:          "",
		})
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call CreatePermission Route from permission Microservice
	var trailer metadata.MD
	_, err = grpc_client.CreatePermissionWithSubjectSet(ctx, params, grpc.Trailer(&trailer))
	if err != nil {
		return me.ToMError(trailer)
	}
	return nil
}

func CreateRelation(permissionParams *pb.CreatePermissionParams) error {

	log.Println("Function CreatePermission() -> Dialing : ", fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call CreatePermission Route from permission Microservice
	_, err = grpc_client.CreatePermission(ctx, permissionParams)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
