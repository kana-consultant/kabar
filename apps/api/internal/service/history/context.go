package history

// SimpleUserContext implements UserContext interface
type SimpleUserContext struct {
	UserID string
	TeamID string
	Role   string
}

func (c *SimpleUserContext) GetUserID() string {
	return c.UserID
}

func (c *SimpleUserContext) GetTeamID() string {
	return c.TeamID
}

func (c *SimpleUserContext) GetRole() string {
	return c.Role
}

func (c *SimpleUserContext) IsAdmin() bool {
	return c.Role == "admin" || c.Role == "super_admin"
}
