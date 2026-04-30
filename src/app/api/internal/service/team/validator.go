package team

import "fmt"

type Validator struct {
	memberRepo *MemberRepository
}

func NewValidator(memberRepo *MemberRepository) *Validator {
	return &Validator{
		memberRepo: memberRepo,
	}
}

func (v *Validator) CheckMemberLimit(teamID string) error {
	currentCount, err := v.memberRepo.GetCount(teamID)
	if err != nil {
		return err
	}

	maxMembers, err := v.memberRepo.GetMaxMembers(teamID)
	if err != nil {
		return err
	}

	if currentCount >= maxMembers {
		return fmt.Errorf("team has reached maximum member limit of %d", maxMembers)
	}

	return nil
}

func (v *Validator) CheckDeleteTeam(teamID string) error {
	memberCount, err := v.memberRepo.GetCount(teamID)
	if err != nil {
		return err
	}

	if memberCount > 0 {
		return fmt.Errorf("cannot delete team with active members")
	}

	return nil
}
