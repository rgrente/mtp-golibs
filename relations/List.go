package relations

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	mgrpc "github.com/rgrente/mtp-golibs/grpc"
	me "github.com/rgrente/mtp-golibs/merror"
	pb "github.com/rgrente/mtp-golibs/pb"
)

func ListObjectRelations(params *pb.CreatePermissionParams) ([]string, error) {
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call ListObjectRelations Route from permission Microservice
	relations, err := grpc_client.ListObjectRelations(ctx, params)
	if err != nil {
		return nil, me.ProcessError(err)
	}

	return relations.Relation, nil
}

func ListSubjectRelations(params *pb.CreatePermissionParams) ([]string, error) {

	log.Println("Function ListSubjectRelations() -> Dialing : ", fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call ListSubjectRelations Route from permission Microservice
	relations, err := grpc_client.ListSubjectRelations(ctx, params)
	if err != nil {
		return nil, err
	}

	return relations.Relation, nil
}

func ListObjectRelationsWithSubjectSet(c *gin.Context, params *pb.CreatePermissionWithSubjectSetParams) ([]string, error) {
	conn, err := mgrpc.InitGRPCClient(fmt.Sprintf("%s:%s", os.Getenv("PERMISSION_ADDR"), os.Getenv("PERMISSION_PORT")))
	if err != nil {
		me.RenderError(c, me.ProcessError(&me.MError{
			Code:           404,
			Message:        "Unprocessable Content",
			Description:    "something went wrong, make sure to provide the right project id or contact the developer team",
			DevCode:        404,
			DevMessage:     "Not Found",
			DevDescription: "No relation found in Keto corresponding to parameters",
			Trace:          "",
		}))
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	grpc_client := pb.NewPermissionClient(conn)

	// Call ListObjectRelationsWithSubjectSet Route from permission Microservice
	var trailer metadata.MD
	relations, err := grpc_client.ListObjectRelationsWithSubjectSet(ctx, params, grpc.Trailer(&trailer))
	if err != nil {
		return nil, me.ToMError(trailer)
	}

	return relations.Relation, nil
}
