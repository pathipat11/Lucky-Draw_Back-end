package migrations

import "app/app/model"

func Models() []any {
	return []any{
		// (*model.Permission)(nil),
		// (*model.RolePermission)(nil),
		// (*model.Role)(nil),
		// (*model.User)(nil),
		(*model.User)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
