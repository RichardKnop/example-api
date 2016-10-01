package accounts

import (
	"github.com/jinzhu/gorm"
)

// AccountPreload sets up Gorm preloads for an account object
func AccountPreload(db *gorm.DB) *gorm.DB {
	return AccountPreloadWithPrefix(db, "")
}

// AccountPreloadWithPrefix sets up Gorm preloads for an account object, and prefixes with prefix for nested objects
func AccountPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "OauthClient")
}

// UserPreload sets up Gorm preloads for a user object
func UserPreload(db *gorm.DB) *gorm.DB {
	return UserPreloadWithPrefix(db, "")
}

// UserPreloadWithPrefix sets up Gorm preloads for a user object,
// and prefixes with prefix for nested objects
func UserPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "Account.OauthClient").
		Preload(prefix + "OauthUser")
}

// ConfirmationPreload sets up Gorm preloads for a confirmation object
func ConfirmationPreload(db *gorm.DB) *gorm.DB {
	return ConfirmationPreloadWithPrefix(db, "")
}

// ConfirmationPreloadWithPrefix sets up Gorm preloads for a confirmation object,
// and prefixes with prefix for nested objects
func ConfirmationPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "User.OauthUser")
}

// InvitationPreload sets up Gorm preloads for an invitation object
func InvitationPreload(db *gorm.DB) *gorm.DB {
	return InvitationPreloadWithPrefix(db, "")
}

// InvitationPreloadWithPrefix sets up Gorm preloads for an invitation object,
// and prefixes with prefix for nested objects
func InvitationPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "InvitedUser.OauthUser").
		Preload(prefix + "InvitedByUser.OauthUser")
}

// PasswordResetPreload sets up Gorm preloads for a password reset object
func PasswordResetPreload(db *gorm.DB) *gorm.DB {
	return PasswordResetPreloadWithPrefix(db, "")
}

// PasswordResetPreloadWithPrefix sets up Gorm preloads for a password reset object,
// and prefixes with prefix for nested objects
func PasswordResetPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "User.OauthUser")
}
