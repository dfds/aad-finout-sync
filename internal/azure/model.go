package azure

import (
	"errors"
	"strings"
	"time"
)

type GroupsListResponse struct {
	OdataContext  string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value         []struct {
		ID                            string        `json:"id"`
		DeletedDateTime               interface{}   `json:"deletedDateTime"`
		Classification                interface{}   `json:"classification"`
		CreatedDateTime               time.Time     `json:"createdDateTime"`
		CreationOptions               []interface{} `json:"creationOptions"`
		Description                   string        `json:"description"`
		DisplayName                   string        `json:"displayName"`
		ExpirationDateTime            interface{}   `json:"expirationDateTime"`
		GroupTypes                    []interface{} `json:"groupTypes"`
		IsAssignableToRole            interface{}   `json:"isAssignableToRole"`
		Mail                          interface{}   `json:"mail"`
		MailEnabled                   bool          `json:"mailEnabled"`
		MailNickname                  string        `json:"mailNickname"`
		MembershipRule                interface{}   `json:"membershipRule"`
		MembershipRuleProcessingState interface{}   `json:"membershipRuleProcessingState"`
		OnPremisesDomainName          string        `json:"onPremisesDomainName"`
		OnPremisesLastSyncDateTime    time.Time     `json:"onPremisesLastSyncDateTime"`
		OnPremisesNetBiosName         string        `json:"onPremisesNetBiosName"`
		OnPremisesSamAccountName      string        `json:"onPremisesSamAccountName"`
		OnPremisesSecurityIdentifier  string        `json:"onPremisesSecurityIdentifier"`
		OnPremisesSyncEnabled         bool          `json:"onPremisesSyncEnabled"`
		PreferredDataLocation         interface{}   `json:"preferredDataLocation"`
		PreferredLanguage             interface{}   `json:"preferredLanguage"`
		ProxyAddresses                []interface{} `json:"proxyAddresses"`
		RenewedDateTime               time.Time     `json:"renewedDateTime"`
		ResourceBehaviorOptions       []interface{} `json:"resourceBehaviorOptions"`
		ResourceProvisioningOptions   []interface{} `json:"resourceProvisioningOptions"`
		SecurityEnabled               bool          `json:"securityEnabled"`
		SecurityIdentifier            string        `json:"securityIdentifier"`
		Theme                         interface{}   `json:"theme"`
		Visibility                    interface{}   `json:"visibility"`
		OnPremisesProvisioningErrors  []interface{} `json:"onPremisesProvisioningErrors"`
	} `json:"value"`
}

type GroupMembers struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		OdataType         string        `json:"@odata.type"`
		ID                string        `json:"id"`
		BusinessPhones    []interface{} `json:"businessPhones"`
		DisplayName       string        `json:"displayName"`
		GivenName         string        `json:"givenName"`
		JobTitle          string        `json:"jobTitle"`
		Mail              string        `json:"mail"`
		MobilePhone       string        `json:"mobilePhone"`
		OfficeLocation    interface{}   `json:"officeLocation"`
		PreferredLanguage interface{}   `json:"preferredLanguage"`
		Surname           string        `json:"surname"`
		UserPrincipalName string        `json:"userPrincipalName"`
		Department        string        `json:"department"`
	} `json:"value"`
}

type GetAdministrativeUnitsResponse struct {
	OdataContext string                                `json:"@odata.context"`
	Value        []*GetAdministrativeUnitsResponseUnit `json:"value"`
}

type CreateAdministrativeUnitGroupRequest struct {
	OdataType       string        `json:"@odata.type"`
	Description     string        `json:"description"`
	DisplayName     string        `json:"displayName"`
	MailNickname    string        `json:"mailNickname"`
	GroupTypes      []interface{} `json:"groupTypes"`
	MailEnabled     bool          `json:"mailEnabled"`
	SecurityEnabled bool          `json:"securityEnabled"`

	ParentAdministrativeUnitId string `json:"-"`
}

type AddGroupMemberRequest struct {
	OdataId string `json:"@odata.id"`
}

type GetUserViaUPNResponse struct {
	OdataContext      string        `json:"@odata.context"`
	BusinessPhones    []interface{} `json:"businessPhones"`
	DisplayName       string        `json:"displayName"`
	GivenName         string        `json:"givenName"`
	JobTitle          string        `json:"jobTitle"`
	Mail              string        `json:"mail"`
	MobilePhone       interface{}   `json:"mobilePhone"`
	OfficeLocation    interface{}   `json:"officeLocation"`
	PreferredLanguage interface{}   `json:"preferredLanguage"`
	Surname           string        `json:"surname"`
	UserPrincipalName string        `json:"userPrincipalName"`
	ID                string        `json:"id"`
}

type CreateAdministrativeUnitGroupResponse struct {
	OdataContext                  string        `json:"@odata.context"`
	OdataType                     string        `json:"@odata.type"`
	ID                            string        `json:"id"`
	DeletedDateTime               interface{}   `json:"deletedDateTime"`
	Classification                interface{}   `json:"classification"`
	CreatedDateTime               time.Time     `json:"createdDateTime"`
	CreationOptions               []interface{} `json:"creationOptions"`
	Description                   string        `json:"description"`
	DisplayName                   string        `json:"displayName"`
	ExpirationDateTime            interface{}   `json:"expirationDateTime"`
	GroupTypes                    []interface{} `json:"groupTypes"`
	IsAssignableToRole            interface{}   `json:"isAssignableToRole"`
	Mail                          interface{}   `json:"mail"`
	MailEnabled                   bool          `json:"mailEnabled"`
	MailNickname                  string        `json:"mailNickname"`
	MembershipRule                interface{}   `json:"membershipRule"`
	MembershipRuleProcessingState interface{}   `json:"membershipRuleProcessingState"`
	OnPremisesDomainName          interface{}   `json:"onPremisesDomainName"`
	OnPremisesLastSyncDateTime    interface{}   `json:"onPremisesLastSyncDateTime"`
	OnPremisesNetBiosName         interface{}   `json:"onPremisesNetBiosName"`
	OnPremisesSamAccountName      interface{}   `json:"onPremisesSamAccountName"`
	OnPremisesSecurityIdentifier  interface{}   `json:"onPremisesSecurityIdentifier"`
	OnPremisesSyncEnabled         interface{}   `json:"onPremisesSyncEnabled"`
	PreferredDataLocation         interface{}   `json:"preferredDataLocation"`
	PreferredLanguage             interface{}   `json:"preferredLanguage"`
	ProxyAddresses                []interface{} `json:"proxyAddresses"`
	RenewedDateTime               time.Time     `json:"renewedDateTime"`
	ResourceBehaviorOptions       []interface{} `json:"resourceBehaviorOptions"`
	ResourceProvisioningOptions   []interface{} `json:"resourceProvisioningOptions"`
	SecurityEnabled               bool          `json:"securityEnabled"`
	SecurityIdentifier            string        `json:"securityIdentifier"`
	Theme                         interface{}   `json:"theme"`
	Visibility                    interface{}   `json:"visibility"`
	OnPremisesProvisioningErrors  []interface{} `json:"onPremisesProvisioningErrors"`
}

type GetAdministrativeUnitsResponseUnit struct {
	ID                            string      `json:"id"`
	DeletedDateTime               interface{} `json:"deletedDateTime"`
	DisplayName                   string      `json:"displayName"`
	Description                   string      `json:"description"`
	MembershipRule                interface{} `json:"membershipRule"`
	MembershipType                interface{} `json:"membershipType"`
	MembershipRuleProcessingState interface{} `json:"membershipRuleProcessingState"`
	Visibility                    interface{} `json:"visibility"`
}

func (g *GetAdministrativeUnitsResponse) GetUnit(name string) *GetAdministrativeUnitsResponseUnit {
	for _, aUnit := range g.Value {
		if aUnit.DisplayName == name {
			return aUnit
		}
	}

	return nil
}

type GetAdministrativeUnitMembersResponse struct {
	OdataContext  string                                     `json:"@odata.context"`
	OdataNextLink string                                     `json:"@odata.nextLink,omitempty"`
	Value         []GetAdministrativeUnitMembersResponseUnit `json:"value"`
}

type GetAdministrativeUnitMembersResponseUnit struct {
	OdataType                     string        `json:"@odata.type"`
	ID                            string        `json:"id"`
	DeletedDateTime               interface{}   `json:"deletedDateTime"`
	Classification                interface{}   `json:"classification"`
	CreatedDateTime               time.Time     `json:"createdDateTime"`
	CreationOptions               []interface{} `json:"creationOptions"`
	Description                   interface{}   `json:"description"`
	DisplayName                   string        `json:"displayName"`
	ExpirationDateTime            interface{}   `json:"expirationDateTime"`
	GroupTypes                    []interface{} `json:"groupTypes"`
	IsAssignableToRole            interface{}   `json:"isAssignableToRole"`
	Mail                          interface{}   `json:"mail"`
	MailEnabled                   bool          `json:"mailEnabled"`
	MailNickname                  string        `json:"mailNickname"`
	MembershipRule                interface{}   `json:"membershipRule"`
	MembershipRuleProcessingState interface{}   `json:"membershipRuleProcessingState"`
	OnPremisesDomainName          interface{}   `json:"onPremisesDomainName"`
	OnPremisesLastSyncDateTime    interface{}   `json:"onPremisesLastSyncDateTime"`
	OnPremisesNetBiosName         interface{}   `json:"onPremisesNetBiosName"`
	OnPremisesSamAccountName      interface{}   `json:"onPremisesSamAccountName"`
	OnPremisesSecurityIdentifier  interface{}   `json:"onPremisesSecurityIdentifier"`
	OnPremisesSyncEnabled         interface{}   `json:"onPremisesSyncEnabled"`
	PreferredDataLocation         interface{}   `json:"preferredDataLocation"`
	PreferredLanguage             interface{}   `json:"preferredLanguage"`
	ProxyAddresses                []interface{} `json:"proxyAddresses"`
	RenewedDateTime               time.Time     `json:"renewedDateTime"`
	ResourceBehaviorOptions       []interface{} `json:"resourceBehaviorOptions"`
	ResourceProvisioningOptions   []interface{} `json:"resourceProvisioningOptions"`
	SecurityEnabled               bool          `json:"securityEnabled"`
	SecurityIdentifier            string        `json:"securityIdentifier"`
	Theme                         interface{}   `json:"theme"`
	Visibility                    interface{}   `json:"visibility"`
	OnPremisesProvisioningErrors  []interface{} `json:"onPremisesProvisioningErrors"`
}

type GetApplicationRolesResponse struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		DisplayName string `json:"displayName"`
		AppID       string `json:"appId"`
		AppRoles    []struct {
			AllowedMemberTypes []string    `json:"allowedMemberTypes"`
			Description        string      `json:"description"`
			DisplayName        string      `json:"displayName"`
			ID                 string      `json:"id"`
			IsEnabled          bool        `json:"isEnabled"`
			Origin             string      `json:"origin"`
			Value              interface{} `json:"value"`
		} `json:"appRoles"`
	} `json:"value"`
}

func (g *GetApplicationRolesResponse) GetRoleId(name string) (string, error) {
	if len(g.Value) != 1 {
		return "", errors.New("too many application entries returned, only expected one")
	}
	app := g.Value[0]
	for _, role := range app.AppRoles {
		if role.DisplayName == name {
			return role.ID, nil
		}
	}

	return "", errors.New("application role not found")
}

type AssignGroupToApplicationRequest struct {
	PrincipalID string `json:"principalId"`
	ResourceID  string `json:"resourceId"`
	AppRoleID   string `json:"appRoleId"`
}

type AssignGroupToApplicationResponse struct {
	OdataContext         string      `json:"@odata.context"`
	ID                   string      `json:"id"`
	DeletedDateTime      interface{} `json:"deletedDateTime"`
	AppRoleID            string      `json:"appRoleId"`
	CreatedDateTime      time.Time   `json:"createdDateTime"`
	PrincipalDisplayName string      `json:"principalDisplayName"`
	PrincipalID          string      `json:"principalId"`
	PrincipalType        string      `json:"principalType"`
	ResourceDisplayName  string      `json:"resourceDisplayName"`
	ResourceID           string      `json:"resourceId"`
}

type GetAssignmentsForApplicationResponse struct {
	OdataContext  string                                            `json:"@odata.context"`
	OdataNextLink string                                            `json:"@odata.nextLink,omitempty"`
	Value         []*GetAssignmentsForApplicationResponseAssignment `json:"value"`
}

type GetAssignmentsForApplicationResponseAssignment struct {
	ID                   string    `json:"id"`
	CreationTimestamp    time.Time `json:"creationTimestamp"`
	AppRoleID            string    `json:"appRoleId"`
	PrincipalDisplayName string    `json:"principalDisplayName"`
	PrincipalID          string    `json:"principalId"`
	PrincipalType        string    `json:"principalType"`
	ResourceDisplayName  string    `json:"resourceDisplayName"`
	ResourceID           string    `json:"resourceId"`
}

func (g *GetAssignmentsForApplicationResponse) ContainsGroup(name string) bool {
	for _, assignment := range g.Value {
		if assignment.PrincipalDisplayName == name {
			return true
		}
	}

	return false
}

func (g *GetAssignmentsForApplicationResponse) GetAssignmentByGroupName(name string) *GetAssignmentsForApplicationResponseAssignment {
	for _, assignment := range g.Value {
		if assignment.PrincipalDisplayName == name {
			return assignment
		}
	}

	return nil
}

type Group struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Members     []*Member `json:"members"`
}

func (g *Group) HasMember(email string) bool {
	for _, member := range g.Members {
		if strings.ToLower(member.UserPrincipalName) == strings.ToLower(email) {
			return true
		}
	}

	return false
}

type Member struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	UserPrincipalName string `json:"userPrincipalName"`
}
