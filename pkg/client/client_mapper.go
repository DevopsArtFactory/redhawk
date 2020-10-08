package client

var (
	clientMapper = map[string]func(Helper) (Client, error){
		"ec2": NewEC2Client,
	}
)
