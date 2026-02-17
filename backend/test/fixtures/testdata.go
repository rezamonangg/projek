package fixtures

import "strings"

type E2ERegisterCommunityInput struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

func E2ENewCommunityInput() E2ERegisterCommunityInput {
	return E2ERegisterCommunityInput{
		Name:          "Test Community",
		Slug:          "test-community",
		AdminEmail:    "admin@test.com",
		AdminPassword: "password123",
		FirstName:     "Admin",
		LastName:      "User",
	}
}

func E2ENewCommunityInputWithEmail(email string) E2ERegisterCommunityInput {
	input := E2ENewCommunityInput()
	input.AdminEmail = email
	input.Slug = email[:strings.Index(email, "@")]
	return input
}

type E2ELoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func E2ENewLoginInput(email, password string) E2ELoginInput {
	return E2ELoginInput{
		Email:    email,
		Password: password,
	}
}

type E2ECreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

func E2ENewProjectInput() E2ECreateProjectInput {
	return E2ECreateProjectInput{
		Name:        "Test Project",
		Description: "A test project for E2E testing",
		Key:         "TEST",
	}
}

func E2ENewProjectInputWithKey(name, key string) E2ECreateProjectInput {
	return E2ECreateProjectInput{
		Name:        name,
		Description: "Project: " + name,
		Key:         key,
	}
}

type E2ECreateBoardInput struct {
	Name string `json:"name"`
}

func E2ENewBoardInput() E2ECreateBoardInput {
	return E2ECreateBoardInput{
		Name: "Test Board",
	}
}

type E2ECreateTaskInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	ReporterID  string  `json:"reporter_id"`
	EpicID      *string `json:"epic_id,omitempty"`
	AssigneeID  *string `json:"assignee_id,omitempty"`
}

func E2ENewTaskInput(reporterID string) E2ECreateTaskInput {
	return E2ECreateTaskInput{
		Title:       "Test Task",
		Description: "A test task",
		Status:      "backlog",
		ReporterID:  reporterID,
	}
}

func E2ENewTaskInputWithStatus(reporterID, title, status string) E2ECreateTaskInput {
	return E2ECreateTaskInput{
		Title:       title,
		Description: "Task: " + title,
		Status:      status,
		ReporterID:  reporterID,
	}
}

type E2EInviteMemberInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func E2ENewInviteInput(email string) E2EInviteMemberInput {
	return E2EInviteMemberInput{
		Email: email,
		Role:  "member",
	}
}

func E2ENewAdminInviteInput(email string) E2EInviteMemberInput {
	return E2EInviteMemberInput{
		Email: email,
		Role:  "admin",
	}
}

type E2ECreateLabelInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func E2ENewLabelInput() E2ECreateLabelInput {
	return E2ECreateLabelInput{
		Name:  "Bug",
		Color: "#FF0000",
	}
}

func E2ENewLabelInputWithValues(name, color string) E2ECreateLabelInput {
	return E2ECreateLabelInput{
		Name:  name,
		Color: color,
	}
}

type E2ECreateEpicInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

func E2ENewEpicInput() E2ECreateEpicInput {
	return E2ECreateEpicInput{
		Name:        "Test Epic",
		Description: "A test epic",
	}
}

func E2ENewEpicInputWithName(name string) E2ECreateEpicInput {
	return E2ECreateEpicInput{
		Name:        name,
		Description: "Epic: " + name,
	}
}

type E2EUpdateProfileInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func E2ENewUpdateProfileInput(firstName, lastName string) E2EUpdateProfileInput {
	return E2EUpdateProfileInput{
		FirstName: firstName,
		LastName:  lastName,
	}
}

type E2EMoveTaskInput struct {
	Status   string `json:"status"`
	Position int    `json:"position"`
}

func E2ENewMoveTaskInput(status string, position int) E2EMoveTaskInput {
	return E2EMoveTaskInput{
		Status:   status,
		Position: position,
	}
}

type E2EAssignLabelInput struct {
	LabelID string `json:"label_id"`
}

func E2ENewAssignLabelInput(labelID string) E2EAssignLabelInput {
	return E2EAssignLabelInput{
		LabelID: labelID,
	}
}

type E2ECreateWikiPageInput struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id,omitempty"`
}

func E2ENewWikiPageInput() E2ECreateWikiPageInput {
	return E2ECreateWikiPageInput{
		Title:   "Test Wiki Page",
		Content: "This is test wiki content",
	}
}

type E2EUpdateSettingsInput struct {
	AllowMemberRegistration  bool `json:"allow_member_registration"`
	RequireEmailVerification bool `json:"require_email_verification"`
}

func E2ENewUpdateSettingsInput(allowReg, requireVerify bool) E2EUpdateSettingsInput {
	return E2EUpdateSettingsInput{
		AllowMemberRegistration:  allowReg,
		RequireEmailVerification: requireVerify,
	}
}
