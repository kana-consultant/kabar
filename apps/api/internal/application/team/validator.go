// internal/application/team/validator.go
package team

import (
	"context"
	"fmt"

	"seo-backend/internal/domain/team"
)

type ValidatorImpl struct {
	memberRepo team.MemberRepository
}

func NewValidator(memberRepo team.MemberRepository) team.Validator {
	return &ValidatorImpl{
		memberRepo: memberRepo,
	}
}

func (v *ValidatorImpl) CheckMemberLimit(ctx context.Context, teamID string) error {
	currentCount, err := v.memberRepo.GetCount(ctx, teamID)
	if err != nil {
		return err
	}

	maxMembers, err := v.memberRepo.GetMaxMembers(ctx, teamID)
	if err != nil {
		return err
	}

	if currentCount >= maxMembers {
		return fmt.Errorf("team has reached maximum member limit of %d", maxMembers)
	}

	return nil
}

func (v *ValidatorImpl) CheckDeleteTeam(ctx context.Context, teamID string) error {
	memberCount, err := v.memberRepo.GetCount(ctx, teamID)
	if err != nil {
		return err
	}

	if memberCount > 0 {
		return fmt.Errorf("cannot delete team with active members")
	}

	return nil
}
