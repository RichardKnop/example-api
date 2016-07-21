package accounts

import (
	"errors"
	"time"

	"github.com/RichardKnop/example-api/accounts/roles"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrInvitationNotFound ...
	ErrInvitationNotFound = errors.New("Invitation not found")
)

// FindInvitationByID looks up an invitation by ID and returns it
func (s *Service) FindInvitationByID(id uint) (*Invitation, error) {
	// Fetch from the database
	invitation := new(Invitation)
	notFound := s.db.
		Preload("InvitedUser.OauthUser").
		First(invitation, id).
		RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrInvitationNotFound
	}

	return invitation, nil
}

// FindInvitationByReference looks up an invitation by a reference and returns it
func (s *Service) FindInvitationByReference(reference string) (*Invitation, error) {
	// Fetch the invitation from the database
	invitation := new(Invitation)
	notFound := s.db.Where("reference = ?", reference).
		Preload("InvitedUser.OauthUser").First(invitation).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrInvitationNotFound
	}

	return invitation, nil
}

// InviteUser invites a new user and sends an invitation email
func (s *Service) InviteUser(invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error) {
	// Begin a transaction
	tx := s.db.Begin()

	invitation, err := s.inviteUserCommon(tx, invitedByUser, invitationRequest)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return invitation, nil
}

// InviteUserTx invites a new user and sends an invitation email in a transaction
func (s *Service) InviteUserTx(tx *gorm.DB, invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error) {
	return s.inviteUserCommon(tx, invitedByUser, invitationRequest)
}

// ConfirmInvitation sets password on the oauth user object and deletes the invitation
func (s *Service) ConfirmInvitation(invitation *Invitation, password string) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Set the new password
	err := s.oauthService.SetPasswordTx(tx, invitation.InvitedUser.OauthUser, password)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Soft delete the invitation
	if err := tx.Delete(invitation).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

func (s *Service) inviteUserCommon(db *gorm.DB, invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error) {
	// Check if oauth user exists
	if s.GetOauthService().UserExists(invitationRequest.Email) {
		return nil, oauth.ErrUsernameTaken
	}

	// Fetch the user from the database
	role, err := s.FindRoleByID(roles.User)
	if err != nil {
		return nil, err
	}

	// Create a new oauth user without a password
	oauthUser, err := s.GetOauthService().CreateUserTx(
		db,
		invitationRequest.Email,
		"", // password
	)
	if err != nil {
		return nil, err
	}

	// Create a new user account
	invitedUser := NewUser(
		invitedByUser.Account,
		oauthUser,
		role,
		"", // facebook ID
		invitationRequest.FirstName,
		invitationRequest.LastName,
		"",    // picture
		false, // confirmed
	)

	// Save the user to the database
	if err := db.Create(invitedUser).Error; err != nil {
		return nil, err
	}

	// Assign related objects
	invitedUser.Account = invitedByUser.Account
	invitedUser.OauthUser = oauthUser
	invitedUser.Role = role

	// Update the meta user ID field
	err = db.Model(oauthUser).UpdateColumn(oauth.User{MetaUserID: invitedUser.ID}).Error
	if err != nil {
		return nil, err
	}

	// Create a new invitation
	invitation := NewInvitation(invitedUser, invitedByUser)
	if err := db.Create(invitation).Error; err != nil {
		return nil, err
	}

	// Assign related objects
	invitation.InvitedUser = invitedUser
	invitation.InvitedByUser = invitedByUser

	// Send invitation email
	go func() {
		invitationEmail, err := s.emailFactory.NewInvitationEmail(invitation)
		if err != nil {
			logger.Errorf("New invitation email error: %s", err)
			return
		}

		// Try to send the invitation email
		if err := s.emailService.Send(invitationEmail); err != nil {
			logger.Errorf("Send email error: %s", err)
			return
		}

		// If the email was sent successfully, update the email_sent flag
		now := time.Now()
		s.db.Model(invitation).UpdateColumns(Invitation{
			EmailSent:   true,
			EmailSentAt: util.TimeOrNull(&now),
			Model:       gorm.Model{UpdatedAt: time.Now()},
		})
	}()

	return invitation, nil
}
