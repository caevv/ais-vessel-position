package system_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caevv/ais-vessel-position/env"
	"github.com/caevv/ais-vessel-position/api"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
)

func TestGetDistance(t *testing.T) {
	apitest.New().
		Handler(api.New(&env.Config{
			RepositoryJsonPath: "data/",
			Files: []string{"202007291231.json", "202007291931.json", "202007292331.json"},
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
		Handler(api.New(&env.Config{Files: []string{"non-existent-file"}}).Router).
		Get("/vessel/position/1000710").
		Expect(t).
		Assert(jsonpath.Contains(`$.error`, "no such file or directory")).
		Status(http.StatusBadRequest).
		End()
}
