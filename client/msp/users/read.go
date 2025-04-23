package users

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

func ReadCreatedUsersInTenant(ctx context.Context, client http.Client, readInput MspUsersInput) (*[]ComputedUserDetails, error) {
	client.Logger.Printf("Reading users in tenant %s\n", readInput.TenantUid)

	// create a map of the users that were created
	// find the list of deleted users by removing from the list every time a user is found in the response
	readUserDetailsMap := map[string]UserDetails{}
	for _, createdUser := range readInput.Users {
		client.Logger.Printf("Adding user %s to readUserDetailsMap\n", createdUser.Username)
		readUserDetailsMap[createdUser.Username] = createdUser
	}

	// get human users
	readHumanUserDetailsMap, err := fetchUsersInTenant(ctx, client, readInput, false)
	if err != nil {
		return nil, err
	}
	readApiOnlyUserDetailsMap, err := fetchUsersInTenant(ctx, client, readInput, true)
	if err != nil {
		return nil, err
	}

	return mergeUserDetails(readHumanUserDetailsMap, readApiOnlyUserDetailsMap), nil
}

func mergeUserDetails(map1, map2 *map[string]ComputedUserDetails) *[]ComputedUserDetails {
	mergedUserDetails := make([]ComputedUserDetails, 0)

	// Add all entries from map1
	for _, value := range *map1 {
		mergedUserDetails = append(mergedUserDetails, value)
	}

	// Add all entries from map2 (overwriting if key exists)
	for _, value := range *map2 {
		mergedUserDetails = append(mergedUserDetails, value)
	}

	return &mergedUserDetails
}

func fetchUsersInTenant(ctx context.Context, client http.Client, readInput MspUsersInput, shouldFetchApiOnlyUsers bool) (*map[string]ComputedUserDetails, error) {
	limit := 200
	offset := 0
	count := 1
	var userPage UserPage
	readUserDetails := map[string]ComputedUserDetails{}
	expectedUsernames := buildExpectedUsernamesFrom(readInput)
	client.Logger.Printf("Expected usernames: %v\n", expectedUsernames)

	for count > offset {
		readUrl := getUrl(client.BaseUrl(), readInput.TenantUid, limit, offset, shouldFetchApiOnlyUsers)
		req := client.NewGet(ctx, readUrl)
		if err := req.Send(&userPage); err != nil {
			return nil, err
		}

		for _, userDetails := range userPage.Items {
			matchingUserDetails := getMatchingUser(&userDetails, expectedUsernames)
			if matchingUserDetails != nil {
				client.Logger.Printf("Found user %s in tenant %s\n", matchingUserDetails.Username, readInput.TenantUid)
				readUserDetails[matchingUserDetails.Username] = *matchingUserDetails
			} else {
				client.Logger.Printf("User %s not found in expected users\n", userDetails.Username)
			}
		}

		offset += limit
		count = userPage.Count
	}

	client.Logger.Printf("Got %d users in tenant %s\n", count, readInput.TenantUid)
	client.Logger.Printf("Found %d created users in tenant %s\n", len(readUserDetails), readInput.TenantUid)

	return &readUserDetails, nil
}

func getMatchingUser(userDetails *UserDetails, expectedUsers mapset.Set[string]) *ComputedUserDetails {
	var strippedUsername string
	if userDetails.ApiOnlyUser {
		strippedUsername = getApiOnlyUsernameWithoutTenant(userDetails.Username)
	} else {
		strippedUsername = userDetails.Username
	}
	if expectedUsers.Contains(strippedUsername) {
		return &ComputedUserDetails{
			Uid:                          userDetails.Uid,
			Username:                     strippedUsername,
			UsernameInSccFirewallManager: userDetails.Username,
			Roles:                        userDetails.Roles,
			ApiOnlyUser:                  userDetails.ApiOnlyUser,
		}
	}
	return nil
}

func getUrl(baseUrl string, tenantUid string, limit int, offset int, shouldFetchApiOnlyUsers bool) string {
	if shouldFetchApiOnlyUsers {
		return url.GetApiOnlyUsersInMspManagedTenant(baseUrl, tenantUid, limit, offset)
	}
	return url.GetUsersInMspManagedTenant(baseUrl, tenantUid, limit, offset)
}

func buildExpectedUsernamesFrom(readInput MspUsersInput) mapset.Set[string] {
	usernames := mapset.NewSet[string]()
	for _, user := range readInput.Users {
		usernames.Add(user.Username)
	}
	return usernames
}

// an API-only username looks like username@tenantName. Strip off the tenant name
func getApiOnlyUsernameWithoutTenant(username string) string {
	parts := strings.SplitN(username, "@", 2)
	return parts[0]
}
