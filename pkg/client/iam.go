/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"
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

	groupData, userMapList, err := i.ScanGroup()
	if err != nil {
		return nil, err
	}

	if groupData != nil {
		result = append(result, groupData...)
	}

	userData, err := i.ScanUser(userMapList)
	if err != nil {
		return nil, err
	}

	if userData != nil {
		result = append(result, userData...)
	}

	roleData, err := i.ScanRole()
	if err != nil {
		return nil, err
	}

	if roleData != nil {
		result = append(result, roleData...)
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

// ScanGroup scans all IAM group
func (i IAMClient) ScanGroup() ([]resource.Resource, map[string][]string, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	logrus.Debug("Start scanning all IAM group list in the account")
	groupList, err := i.GetGroupList()
	if err != nil {
		return nil, nil, err
	}

	if len(groupList) == 0 {
		logrus.Debug("no IAM group found")
		return nil, nil, nil
	}

	userGroupMap := map[string][]string{}

	input := make(chan *resource.IAMGroupResource)
	output := make(chan []resource.Resource)
	defer close(output)

	go func(input chan *resource.IAMGroupResource, output chan []resource.Resource, wg *sync.WaitGroup) {
		var ret []resource.Resource
		for result := range input {
			if result != nil {
				ret = append(ret, result)
			}
			wg.Done()
		}

		output <- ret
	}(input, output, &wg)

	var mutex = new(sync.RWMutex)

	f := func(group *iam.Group, userGroupMap map[string][]string, ch chan *resource.IAMGroupResource) {
		logrus.Debugf("group found: %s", *group.GroupName)
		tmp := resource.IAMGroupResource{
			ResourceType: aws.String(constants.IAMGroupResourceName),
		}

		tmp.GroupName = group.GroupName

		policies, err := i.GetGroupPolicies(*group.GroupName)
		if err != nil {
			logrus.Error(err.Error())
			ch <- nil
			return
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
			ch <- nil
			return
		}

		tmp.UserCount = aws.Int(len(userListInGroup))
		var ul []string
		for _, u := range userListInGroup {
			mutex.Lock()

			cu := *u.UserName
			if _, ok := userGroupMap[cu]; !ok {
				userGroupMap[cu] = []string{}
			}

			userGroupMap[cu] = append(userGroupMap[cu], *group.GroupName)
			ul = append(ul, cu)

			mutex.Unlock()
		}

		if len(ul) == 0 {
			tmp.Users = aws.String(constants.EmptyString)
		} else {
			tmp.Users = aws.String(strings.Join(ul, constants.DefaultDelimiter))
		}

		ch <- &tmp
	}

	logrus.Debugf("Group found: %d", len(groupList))
	for _, group := range groupList {
		wg.Add(1)
		go f(group, userGroupMap, input)
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid IAM group data count: %d", len(result))

	return result, userGroupMap, nil
}

// ScanUser scans all IAM group
func (i IAMClient) ScanUser(userGroupMap map[string][]string) ([]resource.Resource, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	logrus.Debug("Start scanning all IAM user list in the account")
	userList, err := i.GetUserList()
	if err != nil {
		return nil, err
	}

	if len(userList) == 0 {
		logrus.Debug("no IAM user found")
		return nil, nil
	}

	input := make(chan *resource.IAMUserResource)
	output := make(chan []resource.Resource)
	defer close(output)

	go func(input chan *resource.IAMUserResource, output chan []resource.Resource, wg *sync.WaitGroup) {
		var ret []resource.Resource
		for result := range input {
			if result != nil {
				ret = append(ret, result)
			}
			wg.Done()
		}

		output <- ret
	}(input, output, &wg)

	f := func(user *iam.User, ch chan *resource.IAMUserResource) {
		tmp := resource.IAMUserResource{
			ResourceType: aws.String(constants.IAMUserResourceName),
		}

		tmp.UserName = user.UserName
		tmp.UserCreated = user.CreateDate

		accessKeys, err := i.GetAccessKeys(*user.UserName)
		if err != nil {
			logrus.Errorf(err.Error())
			ch <- nil
			return
		}

		devices, err := i.GetMFADevices(*user.UserName)
		if err != nil {
			logrus.Errorf(err.Error())
			ch <- nil
			return
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
				ch <- nil
				return
			}

			if lastUsed.LastUsedDate == nil {
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

		ch <- &tmp
	}

	logrus.Debugf("User found: %d", len(userList))
	for _, user := range userList {
		wg.Add(1)
		go f(user, input)
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid IAM user data count: %d", len(result))

	return result, nil
}

// ScanRole scans all IAM group
func (i IAMClient) ScanRole() ([]resource.Resource, error) {
	var wg sync.WaitGroup
	var result []resource.Resource

	logrus.Debug("Start scanning all IAM role list in the account")
	roleList, err := i.GetRoleList()
	if err != nil {
		return nil, err
	}

	if len(roleList) == 0 {
		logrus.Debug("no IAM role found")
		return nil, nil
	}

	input := make(chan *resource.IAMRoleResource)
	output := make(chan []resource.Resource)
	defer close(output)

	go func(input chan *resource.IAMRoleResource, output chan []resource.Resource, wg *sync.WaitGroup) {
		var ret []resource.Resource
		for result := range input {
			if result != nil {
				ret = append(ret, result)
			}
			wg.Done()
		}

		output <- ret
	}(input, output, &wg)

	f := func(role *iam.Role, ch chan *resource.IAMRoleResource) {
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
				logrus.Errorf("policy parsing error: %s", *role.RoleName)
				ch <- nil
				return
			}
		} else {
			stringPolicy = *role.AssumeRolePolicyDocument
		}

		err := json.Unmarshal([]byte(stringPolicy), &pd)
		if err != nil {
			logrus.Error(err.Error())
			ch <- nil
			return
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

		ch <- &tmp
	}

	logrus.Debugf("Role found: %d", len(roleList))
	for _, role := range roleList {
		wg.Add(1)
		go f(role, input)
	}

	wg.Wait()
	close(input)

	result = <-output
	logrus.Debugf("total valid IAM role data count: %d", len(result))

	return result, nil
}
