package client

import "github.com/DevopsArtFactory/redhawk/pkg/constants"

var (
	clientMapper = map[string]func(Helper) (Client, error){
		constants.EC2ResourceName: NewEC2Client,
		constants.SGResourceName:  NewSGClient,
		constants.Route53ResourceName: NewRoute53Client,
	}
)
