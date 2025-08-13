package db

import "github.com/JonJenson-MFIn/project-management-system-api/graph/model"

// RoleToDB converts GraphQL Role to database RoleDB
func RoleToDB(role model.Role) RoleDB {
	return RoleDB(role)
}

// RoleToModel converts database RoleDB to GraphQL Role
func RoleToModel(role RoleDB) model.Role {
	return model.Role(role)
}

// StatusToDB converts GraphQL Status to database StatusDB
func StatusToDB(status model.Status) StatusDB {
	return StatusDB(status)
}

// StatusToModel converts database StatusDB to GraphQL Status
func StatusToModel(status StatusDB) model.Status {
	return model.Status(status)
}
