package bookmark_endpoint

import (
	"fmt"
	"net/http"
	"testing"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	"github.com/bookmark-project-learn/bookmark-service/internal/api"
	"github.com/bookmark-project-learn/bookmark-service/internal/config"
	"github.com/bookmark-project-learn/bookmark-service/internal/connection"
	"github.com/bookmark-project-learn/bookmark-service/internal/test/data/fixture"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testBookmarkUserID = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func buildBookmarkIntegrationConfig() *config.Config {
	return &config.Config{
		AppPort:     "8080",
		ServiceName: "app_service",
		InstanceID:  "instance_01",
	}
}

func BuildBookmarkApiEngine(t *testing.T, fix fixture.Fixture, jwtMock *jwt_pkg.MockJwt) api.Engine {
	t.Helper()
	connectorMock, err := connection.InitDBConnectorMock(t, fix)
	if err != nil {
		t.Fatal(err)
	}

	engine := api.NewEngine(&api.EnginOpt{
		App:          gin.New(),
		Cfg:          buildBookmarkIntegrationConfig(),
		Connector:    connectorMock,
		JwtGenerator: jwtMock.JwtGenarate,
		JwtValidator: jwtMock.JwtValidate,
	})

	return engine
}

var dummyUserId = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func RouterSetupAuthorization(eng api.Engine, jwtMock *jwt_pkg.MockJwt, w http.ResponseWriter, req *http.Request) {
	claims := jwt.MapClaims{
		"sub": dummyUserId,
	}
	token, err := jwtMock.JwtGenarate.GenerateToken(claims)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	eng.ServeHTTP(w, req)
}

func buildBookmarkURL(path string) string {
	return fmt.Sprintf("/v1/bookmarks%s", path)
}
