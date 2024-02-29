package keycloak

import (
	"context"
	"os"

	"github.com/Clarilab/gocloaksession"
	gocloak "github.com/Nerzal/gocloak/v13"

	me "github.com/rgrente/mtp-golibs/merror"
)

func GetUsersIdByEmailList(emailList []string) ([]string, error) {
	ctx := context.Background()

	// Create new Keycloak session
	gc := gocloak.NewClient(
		os.Getenv("KEYCLOAK_BASE_URL"),
		gocloak.SetAuthAdminRealms("admin/realms"),
		gocloak.SetAuthRealms("realms"))
	session, e := gocloaksession.NewSession(
		os.Getenv("KEYCLOAK_ADMIN_CLIENT_ID"),
		os.Getenv("KEYCLOAK_ADMIN_CLIENT_SECRET"),
		"master", os.Getenv("KEYCLOAK_BASE_URL"),
		gocloaksession.SetGocloak(gc))
	if e != nil {
		return nil, me.ProcessError(&me.MError{
			Code:           424,
			Message:        "Failed Dependency",
			Description:    "retry in a short moment or contact the developer team",
			DevCode:        424,
			DevMessage:     "Failed Dependency",
			DevDescription: "service Keycloak not reachable",
			Trace:          "",
		})
	}

	// Getting session token
	token, _ := session.GetKeycloakAuthToken()

	// Get usersId
	var usersId []string
	for _, email := range emailList {
		gParams := gocloak.GetUsersParams{
			Email: &email,
		}
		users, err := session.GetGoCloakInstance().GetUsers(ctx, token.AccessToken, os.Getenv("KEYCLOAK_REALM"), gParams)
		if err != nil {
			return nil, me.ProcessError(&me.MError{
				Code:           422,
				Message:        "Unprocessable Content",
				Description:    "retry in a short moment or contact the developer team",
				DevCode:        422,
				DevMessage:     "Unprocessable Content",
				DevDescription: "can't get user id from Keycloak, reponse message: " + err.Error(),
				Trace:          "",
			})
		} else if len(users) != 0 {
			usersId = append(usersId, *users[0].ID)
		}
	}

	return usersId, nil
}
