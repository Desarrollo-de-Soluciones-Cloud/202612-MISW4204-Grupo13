package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

type deletePeriodRepoStub struct {
	period      *domain.Period
	findErr     error
	deleteErr   error
	deletedWith uint
}

func (r *deletePeriodRepoStub) Create(period *domain.Period) error { return nil }

func (r *deletePeriodRepoStub) FindByID(id uint) (*domain.Period, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.period == nil || r.period.ID != id {
		return nil, domain.ErrPeriodNotFound
	}
	return r.period, nil
}

func (r *deletePeriodRepoStub) FindByName(name string) (*domain.Period, error) {
	return nil, domain.ErrPeriodNotFound
}

func (r *deletePeriodRepoStub) FindAll() ([]domain.Period, error) { return nil, nil }

func (r *deletePeriodRepoStub) FindAllByState(state domain.PeriodState) ([]domain.Period, error) {
	return nil, nil
}

func (r *deletePeriodRepoStub) Update(period *domain.Period) error { return nil }

func (r *deletePeriodRepoStub) Delete(id uint) error {
	r.deletedWith = id
	return r.deleteErr
}

func TestDeletePeriodSuccess(t *testing.T) {
	repo := &deletePeriodRepoStub{
		period: &domain.Period{ID: 9},
	}

	err := applicationpkg.NewDeletePeriod(repo).Execute(applicationpkg.DeletePeriodInput{ID: 9})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.deletedWith != 9 {
		t.Fatalf("expected delete to be called with 9, got %d", repo.deletedWith)
	}
}

func TestDeletePeriodPropagatesFindError(t *testing.T) {
	repo := &deletePeriodRepoStub{findErr: domain.ErrPeriodNotFound}

	err := applicationpkg.NewDeletePeriod(repo).Execute(applicationpkg.DeletePeriodInput{ID: 4})
	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Fatalf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestDeletePeriodPropagatesDeleteError(t *testing.T) {
	repo := &deletePeriodRepoStub{
		period:    &domain.Period{ID: 11},
		deleteErr: domain.ErrPeriodNotFound,
	}

	err := applicationpkg.NewDeletePeriod(repo).Execute(applicationpkg.DeletePeriodInput{ID: 11})
	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Fatalf("expected delete error, got %v", err)
	}
}
