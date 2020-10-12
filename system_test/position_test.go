package system_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caevv/ais-vessel-position/env"
	"github.com/caevv/ais-vessel-position/service"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
)

func TestGetDistance(t *testing.T) {
	apitest.New().
		Handler(service.New(&env.Config{
			RepositoryJsonPath: "../repository/data/",
		}).Router).
		Get("/vessel/position/1000710").
		Expect(t).
		Body(
			fmt.Sprintf(
				`{"kilometers": %g, "nauticalMiles": %g, "statuteMiles": %g}`,
				48.723096179857535,
				26.29092147023152,
				30.275128362772367,
			)).
		Status(http.StatusOK).
		End()
}

func TestError(t *testing.T) {
	apitest.New().
		Handler(service.New(&env.Config{}).Router).
		Get("/vessel/position/1000710").
		Expect(t).
		Assert(jsonpath.Contains(`$.error`, "such file or directory")).
		Status(http.StatusOK).
		End()
}
