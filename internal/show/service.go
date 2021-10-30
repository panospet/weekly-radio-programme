package show

import (
	"context"
	"fmt"
	"weekly-radio-programme/common"
)

type Service struct {
	repo Repository
}

func NewService(repo *PostgresRepo) *Service {
	return &Service{repo: repo}
}

type ErrTimeslotConflict struct {
	Message string
}

func (e ErrTimeslotConflict) Error() string {
	return e.Message
}

func (o *Service) Add(ctx context.Context, show Show) (int, error) {
	return o.repo.Add(ctx, show)
}

func (o *Service) GetAll(ctx context.Context) ([]Show, error) {
	return o.repo.GetAll(ctx)
}

func (o *Service) Update(ctx context.Context, show Show) error {
	return o.repo.Update(ctx, show)
}

func (o *Service) Get(ctx context.Context, id int) (Show, error) {
	return o.repo.Get(ctx, id)
}

func (o *Service) Delete(ctx context.Context, id int) error {
	return o.repo.Delete(ctx, id)
}

func (o *Service) CheckForConflicts(ctx context.Context, show Show) error {
	showsSameDay, err := o.repo.GetShowsWithSameWeekday(ctx, show)
	if err != nil {
		return err
	}
	for _, s := range showsSameDay {
		conflict, err := common.HasConflict(show.Timeslot, s.Timeslot)
		if err != nil {
			return fmt.Errorf("cannot calculate conflicts: %s", err.Error())
		}
		if conflict && show.Id != s.Id {
			return ErrTimeslotConflict{
				Message: fmt.Sprintf(`conflict with show '%s' which takes place every %s %s`, s.Title, s.Weekday,
					s.Timeslot)}
		}
	}
	return nil
}
