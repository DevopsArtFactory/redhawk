package client

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
	"github.com/DevopsArtFactory/redhawk/pkg/resource"
	"github.com/DevopsArtFactory/redhawk/pkg/tools"
)

type IAMClient struct {
	Resource string
	Client   *iam.IAM
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Principal *Principal
	Resource  interface{}
}

type Principal struct {
	Service   interface{}
	Federated interface{}
}

// GetResourceName returns resource name of client
func (i IAMClient) GetResourceName() string {
	return i.Resource
}

// NewIAMClient creates IAMClient
func NewIAMClient(helper Helper) (Client, error) {
	session := GetAwsSession()
	return &IAMClient{
		Resource: constants.IAMResourceName,
		Client:   GetIAMClientFn(session, helper.Region, helper.Credentials),
	}, nil
}

// GetIAMClientFn creates iam client
func GetIAMClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *iam.IAM {
	if creds == nil {
		return iam.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return iam.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// Scan scans all data
func (i IAMClient) Scan() ([]resource.Resource, error) {
	var result []resource.Resource

	logrus.Debug("Start scanning all IAM group list in the account")
	groupList, err := i.GetGroupList()
	if err != nil {
		return nil, err
	}

	userGroupMap := map[string][]string{}

	logrus.Debugf("Group found: %d", len(groupList))
	for _, group := range groupList {
		tmp := resource.IAMGroupResource{
			ResourceType: aws.String(constants.IAMGroupResourceName),
		}

		tmp.GroupName = group.GroupName

		policies, err := i.GetGroupPolicies(*group.GroupName)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		var gpList []string
		for _, policy := range policies {
			gpList = append(gpList, *policy.PolicyName)
		}

		if len(gpList) == 0 {
			tmp.GroupPolicies = aws.String(constants.EmptyString)
		} else {
			tmp.GroupPolicies = aws.String(strings.Join(gpList, constants.DefaultDelimiter))
		}

		userListInGroup, err := i.GetUserListInGroup(*group.GroupName)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		tmp.UserCount = aws.Int(len(userListInGroup))
		var ul []string
		for _, u := range userListInGroup {
			cu := *u.UserName
			if _, ok := userGroupMap[cu]; !ok {
				userGroupMap[cu] = []string{}
			}

			userGroupMap[cu] = append(userGroupMap[cu], *group.GroupName)
			ul = append(ul, cu)
		}

		if len(ul) == 0 {
			tmp.Users = aws.String(constants.EmptyString)
		} else {
			tmp.Users = aws.String(strings.Join(ul, constants.DefaultDelimiter))
		}

		result = append(result, tmp)
	}

	logrus.Debug("Start scanning all IAM user list in the account")
	userList, err := i.GetUserList()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("User found: %d", len(userList))
	for _, user := range userList {
		tmp := resource.IAMUserResource{
			ResourceType: aws.String(constants.IAMUserResourceName),
		}

		tmp.UserName = user.UserName
		tmp.UserCreated = user.CreateDate

		accessKeys, err := i.GetAccessKeys(*user.UserName)
		if err != nil {
			logrus.Errorf(err.Error())
			continue
		}

		devices, err := i.GetMFADevices(*user.UserName)
		if err != nil {
			logrus.Errorf(err.Error())
			continue
		}

		var mfaDevices []string
		for _, device := range devices {
			mfaDevices = append(mfaDevices, *device.SerialNumber)
		}

		if len(mfaDevices) == 0 {
			tmp.MFA = aws.String(constants.EmptyString)
		} else {
			tmp.MFA = aws.String(strings.Join(mfaDevices, constants.DefaultDelimiter))
		}

		var lastlyUsed *time.Time
		for _, ak := range accessKeys {
			lastUsed, err := i.GetLastAccessKeyUsed(ak.AccessKeyId)
			if err != nil {
				logrus.Errorf(err.Error())
				continue
			}

			if lastlyUsed == nil || lastUsed.LastUsedDate.Sub(*lastlyUsed) > 0 {
				lastlyUsed = lastUsed.LastUsedDate
			}
		}
		tmp.AccessKeyLastUsed = lastlyUsed

		if _, ok := userGroupMap[*user.UserName]; !ok {
			tmp.GroupCount = aws.Int(0)
		} else {
			tmp.GroupCount = aws.Int(len(userGroupMap[*user.UserName]))
		}

		result = append(result, tmp)
	}

	logrus.Debug("Start scanning all IAM role list in the account")
	roleList, err := i.GetRoleList()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Role found: %d", len(roleList))
	for _, role := range roleList {
		tmp := resource.IAMRoleResource{
			ResourceType: aws.String(constants.IAMRoleResourceName),
		}

		tmp.RoleName = role.RoleName
		if role.RoleLastUsed != nil {
			tmp.RoleLastActivity = role.RoleLastUsed.LastUsedDate
		}

		var pd PolicyDocument

		var stringPolicy string
		if strings.HasPrefix(*role.AssumeRolePolicyDocument, "%") {
			stringPolicy, err = tools.DecodeURLEncodedString(*role.AssumeRolePolicyDocument)
			if err != nil {
				continue
			}
		} else {
			stringPolicy = *role.AssumeRolePolicyDocument
		}
		err := json.Unmarshal([]byte(stringPolicy), &pd)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		var trustedEntities []string
		for _, statement := range pd.Statement {
			if statement.Principal != nil {
				if statement.Principal.Service != nil {
					sps := reflect.TypeOf(statement.Principal.Service).Name()
					if sps != "string" {
						// Service is a slice of string
						serviceSlice := statement.Principal.Service.([]interface{})

						for _, ss := range serviceSlice {
							trustedEntities = append(trustedEntities, ss.(string))
						}
					} else {
						trustedEntities = append(trustedEntities, statement.Principal.Service.(string))
					}
				}

				if statement.Principal.Federated != nil {
					spf := reflect.TypeOf(statement.Principal.Federated).Name()
					if spf != "string" {
						// Service is a slice of string
						federatedSlice := statement.Principal.Federated.([]interface{})

						for _, fs := range federatedSlice {
							trustedEntities = append(trustedEntities, fs.(string))
						}
					} else {
						trustedEntities = append(trustedEntities, statement.Principal.Federated.(string))
					}
				}
			}
		}

		tmp.TrustedEntities = aws.String(strings.Join(trustedEntities, constants.DefaultDelimiter))

		result = append(result, tmp)
	}

	return result, nil
}

// GetUserList returns all IAM User list
func (i IAMClient) GetUserList() ([]*iam.User, error) {
	input := &iam.ListUsersInput{}

	result, err := i.Client.ListUsers(input)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// GetGroupList returns all IAM group list
func (i IAMClient) GetGroupList() ([]*iam.Group, error) {
	input := &iam.ListGroupsInput{}

	result, err := i.Client.ListGroups(input)
	if err != nil {
		return nil, err
	}

	return result.Groups, nil
}

// GetRoleList returns all IAM role list
func (i IAMClient) GetRoleList() ([]*iam.Role, error) {
	input := &iam.ListRolesInput{}

	result, err := i.Client.ListRoles(input)
	if err != nil {
		return nil, err
	}

	return result.Roles, nil
}

// GetAccessKeys returns all access keys of user
func (i IAMClient) GetAccessKeys(user string) ([]*iam.AccessKeyMetadata, error) {
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(user),
	}

	result, err := i.Client.ListAccessKeys(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKeyMetadata, nil
}

// GetLastAccessKeyUsed returns lastly used date of access key
func (i IAMClient) GetLastAccessKeyUsed(accessKey *string) (*iam.AccessKeyLastUsed, error) {
	input := &iam.GetAccessKeyLastUsedInput{
		AccessKeyId: accessKey,
	}

	result, err := i.Client.GetAccessKeyLastUsed(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKeyLastUsed, nil
}

// GetMFADevices returns all MFA devices
func (i IAMClient) GetMFADevices(user string) ([]*iam.MFADevice, error) {
	input := &iam.ListMFADevicesInput{
		UserName: aws.String(user),
	}

	result, err := i.Client.ListMFADevices(input)
	if err != nil {
		return nil, err
	}

	return result.MFADevices, nil
}

// GetUserListInGroup returns user list of group
func (i IAMClient) GetUserListInGroup(group string) ([]*iam.User, error) {
	input := &iam.GetGroupInput{
		GroupName: aws.String(group),
	}

	result, err := i.Client.GetGroup(input)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// GetGroupPolicies returns policies of group
func (i IAMClient) GetGroupPolicies(group string) ([]*iam.AttachedPolicy, error) {
	input := &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(group),
	}

	result, err := i.Client.ListAttachedGroupPolicies(input)
	if err != nil {
		return nil, err
	}

	return result.AttachedPolicies, nil
}
