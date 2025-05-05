package db

type JoinableRole struct {
	GuildID  int64  `gorm:"primaryKey;column:guild_id"`
	RoleID   int64  `gorm:"primaryKey;column:role_id"`
	Nickname string `gorm:"column:nickname"`
}

func GetJoinableRoles(guildID string) ([]JoinableRole, error) {
	var roles []JoinableRole
	err := db.Where("guild_id = ?", guildID).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
