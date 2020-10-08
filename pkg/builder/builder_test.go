package builder

import (
	"testing"

	"github.com/DevopsArtFactory/redhawk/pkg/client"
	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

func TestRecentRegion(t *testing.T) {
	allAWSREgionsFromAPI, err := client.GetAllRegions()
	if err != nil {
		t.Error(err.Error())
	}

	for _, r := range allAWSREgionsFromAPI {
		if !tools.IsStringInArray(r, constants.AllAWSRegions) {
			t.Errorf("region is not added: %s", r)
		}
	}
}
