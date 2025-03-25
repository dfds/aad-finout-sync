package finout

import (
	"strings"
	"time"
)

type ListVirtualTagResponseTag struct {
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	ID        string `json:"id"`
}

type CreateVirtualTagResponse struct {
	AccountID string                         `json:"accountId"`
	Name      string                         `json:"name"`
	Rules     []CreateVirtualTagResponseRule `json:"rules"`
	Category  string                         `json:"category"`
	CreatedBy string                         `json:"createdBy"`
	UpdatedBy string                         `json:"updatedBy"`
	Default   struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"default"`
	Endpoints []interface{} `json:"endpoints"`
	AnomalyID string        `json:"anomalyId"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
	ID        string        `json:"id"`
}

type CreateVirtualTagResponseRule struct {
	To      string                             `json:"to"`
	Filters CreateVirtualTagResponseRuleFilter `json:"filters"`
}

type CreateVirtualTagResponseRuleFilter struct {
	CostCenter string   `json:"costCenter"`
	Type       string   `json:"type"`
	Operator   string   `json:"operator"`
	Value      []string `json:"value"`
	Name       string   `json:"name"`
}

type CreateVirtualTagRequest struct {
	Default CreateVirtualTagRequestDefault `json:"default"`
	Rules   []CreateVirtualTagRequestRule  `json:"rules"`
	//Category  string                         `json:"category"`
	Name      string `json:"name"`
	UpdatedBy string `json:"updatedBy"`
	CreatedBy string `json:"createdBy"`
}

type CreateVirtualTagRequestDefault struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CreateVirtualTagRequestRule struct {
	To      string                            `json:"to"`
	Filters CreateVirtualTagRequestRuleFilter `json:"filters"`
}

type CreateVirtualTagRequestRuleFilter struct {
	CostCenter string   `json:"costCenter"`
	Key        string   `json:"key"`
	Type       string   `json:"type"`
	Operator   string   `json:"operator"`
	Value      []string `json:"value"`
}

type UpdateVirtualTagRequest struct {
	Rules []UpdateVirtualTagRequestRule `json:"rules"`
	//Category  string                        `json:"category"`
	Endpoints []string `json:"endpoints"`
	Name      string   `json:"name"`
	//UpdatedBy string                         `json:"updatedBy"`
	Default CreateVirtualTagRequestDefault `json:"default"`
}

type UpdateVirtualTagRequestRule struct {
	To            string                            `json:"to"`
	DisableRename *bool                             `json:"disableRename,omitempty"`
	Filters       UpdateVirtualTagRequestRuleFilter `json:"filters"`
}

type UpdateVirtualTagRequestRuleFilter struct {
	CostCenter string   `json:"costCenter"`
	Key        string   `json:"key"`
	Type       string   `json:"type"`
	Operator   string   `json:"operator"`
	Value      []string `json:"value"`
	Path       string   `json:"path,omitempty"`
}

type UpdateVirtualTagResponse struct {
	AccountID string                         `json:"accountId"`
	Name      string                         `json:"name"`
	Rules     []UpdateVirtualTagResponseRule `json:"rules"`
	Category  string                         `json:"category"`
	CreatedBy string                         `json:"createdBy"`
	UpdatedBy string                         `json:"updatedBy"`
	Default   struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"default"`
	Endpoints []interface{} `json:"endpoints"`
	AnomalyID string        `json:"anomalyId"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
	ID        string        `json:"id"`
}

type UpdateVirtualTagResponseRule struct {
	To      string                             `json:"to"`
	Filters UpdateVirtualTagResponseRuleFilter `json:"filters,omitempty"`
}

type UpdateVirtualTagResponseRuleFilter struct {
	CostCenter string   `json:"costCenter"`
	Type       string   `json:"type"`
	Operator   string   `json:"operator"`
	Value      []string `json:"value"`
	Path       string   `json:"path"`
	Name       string   `json:"name"`
}

type ListViewsResponse struct {
	Data      []ListViewsResponseData `json:"data"`
	RequestID string                  `json:"requestId"`
}

type ListViewsResponseData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (l *ListViewsResponse) GetByName(value string) *ListViewsResponseData {
	for _, view := range l.Data {
		if strings.ToLower(view.Name) == strings.ToLower(value) {
			return &view
		}
	}

	return nil
}

type QueryByViewResponse struct {
	Data      []QueryByViewResponseData  `json:"data"`
	Request   QueryByViewResponseRequest `json:"request"`
	RequestID string                     `json:"requestId"`
}

type QueryByViewResponseData struct {
	Name string                        `json:"name"`
	Data []QueryByViewResponseDataData `json:"data"`
}

type QueryByViewResponseDataData struct {
	Time int64   `json:"time"`
	Cost float64 `json:"cost"`
}

func (q *QueryByViewResponseDataData) TimeFormatted() string {
	return time.UnixMilli(q.Time).Format(time.DateOnly)
}

type QueryByViewResponseRequest struct {
	ViewID string `json:"viewId"`
}

type QueryByViewRequest struct {
	ViewId string `json:"viewId"`
}

type ListGroupsResponse struct {
	Groups []ListGroupsResponseGroup `json:"groups"`
}

type ListGroupsResponseGroup struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Color       interface{}   `json:"color"`
	Description interface{}   `json:"description"`
	Metadata    interface{}   `json:"metadata"`
	Roles       []interface{} `json:"roles"`
	Users       []struct {
		ID                 string    `json:"id"`
		Name               string    `json:"name"`
		ProfilePictureURL  string    `json:"profilePictureUrl"`
		Email              string    `json:"email"`
		CreatedAt          time.Time `json:"createdAt"`
		ActivatedForTenant bool      `json:"activatedForTenant"`
	} `json:"users"`
	ManagedBy string `json:"managedBy"`
}

type CreateGroupResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Color       interface{}   `json:"color"`
	Description string        `json:"description"`
	Metadata    string        `json:"metadata"`
	Roles       []interface{} `json:"roles"`
	Users       []interface{} `json:"users"`
	ManagedBy   string        `json:"managedBy"`
}

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`
}

type AddUsersToGroupRequest struct {
	UserIds []string `json:"userIds"`
}

type AddRolesToGroupRequest struct {
	RoleIds []string `json:"roleIds"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Metadata    string `json:"metadata"`
}

type UpdateAccountDataAccessForGroupsRequestGroupConfig struct {
	DataAccessEnabled bool `json:"dataAccessEnabled"`
	Filters           struct {
		CostCenter string   `json:"costCenter"`
		Key        string   `json:"key"`
		Operator   string   `json:"operator"`
		Value      []string `json:"value"`
		Type       string   `json:"type"`
	} `json:"filters"`
}

type UpdateAccountDataAccessForGroupsRequest struct {
	GroupsConfig map[string]UpdateAccountDataAccessForGroupsRequestGroupConfig `json:"groupsConfig"`
}

type UpdateAccountDataAccessForGroupsResponse struct {
	AccountId string `json:"accountId"`
	Name      string `json:"name"`
	PayerId   string `json:"payerId"`
	Plan      struct {
		Type             string    `json:"type"`
		EndDate          time.Time `json:"endDate"`
		UpgradeRequested bool      `json:"upgradeRequested"`
	} `json:"plan"`
	Configuration struct {
		IsCostGuardEnabled bool `json:"isCostGuardEnabled"`
	} `json:"configuration"`
	CreatedAt                   string `json:"createdAt"`
	UpdatedAt                   string `json:"updatedAt"`
	DefaultContextId            string `json:"defaultContextId"`
	LatestCompletedRunTimestamp int64  `json:"latestCompletedRunTimestamp"`
	DataLocation                struct {
		Bucket string `json:"bucket"`
		Path   string `json:"path"`
	} `json:"dataLocation"`
	SignTermsAndConditions struct {
		ApprovedAt string `json:"approvedAt"`
		ApprovedBy string `json:"approvedBy"`
	} `json:"signTermsAndConditions"`
	SignCostOptimizerTermsAndConditions struct {
		ApprovedAt time.Time `json:"approvedAt"`
		ApprovedBy string    `json:"approvedBy"`
	} `json:"signCostOptimizerTermsAndConditions"`
	GroupsConfig map[string]UpdateAccountDataAccessForGroupsRequestGroupConfig `json:"groupsConfig"`
	Id           string                                                        `json:"id"`
}

// auth

type VerifyMfaResponse struct {
	MfaRequired  bool   `json:"mfaRequired"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	Expires      string `json:"expires"`
	UserID       string `json:"userId"`
}

type VerifyMfaRequest struct {
	MfaToken       string `json:"mfaToken"`
	Value          string `json:"value"`
	RememberDevice bool   `json:"rememberDevice"`
}

type UserLoginResponse struct {
	MfaRequired bool   `json:"mfaRequired"`
	MfaToken    string `json:"mfaToken"`
	MfaDevices  struct {
		Webauthn       []interface{} `json:"webauthn"`
		Phones         []interface{} `json:"phones"`
		Emails         []interface{} `json:"emails"`
		Authenticators []struct {
			ID string `json:"id"`
		} `json:"authenticators"`
	} `json:"mfaDevices"`
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
	Expires       string `json:"expires"`
	ExpiresIn     int    `json:"expiresIn"`
	UserID        string `json:"userId"`
	UserEmail     string `json:"userEmail"`
	EmailVerified bool   `json:"emailVerified"`
}

type UserLoginRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	InvitationToken string `json:"invitationToken"`
}
