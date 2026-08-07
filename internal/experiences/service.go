package experiences

import (
	"context"
	"fmt"
	"math"
	"time"

	"backend-porto/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

// JobInput is the input for a job within an experience.
type JobInput struct {
	Position       string   `json:"position"`
	EmploymentType string   `json:"employment_type"`
	Description    string   `json:"description"`
	Competency     []string `json:"competency,omitempty"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date,omitempty"`
	IsCurrent      bool     `json:"is_current"`
}

// CreateInput is the input for creating an experience with jobs.
type CreateInput struct {
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Tags     []string `json:"tags"`
	Jobs     []JobInput `json:"jobs"`
}

// UpdateInput is the input for updating an experience with jobs.
type UpdateInput struct {
	ID       string   `json:"id"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Tags     []string `json:"tags"`
	Jobs     []JobInput `json:"jobs"`
}

// ExperienceWithJobs is an experience with its nested jobs.
type ExperienceWithJobs struct {
	ID            pgtype.UUID      `json:"id"`
	Company       string           `json:"company"`
	Location      string           `json:"location"`
	TotalDuration string           `json:"total_duration"`
	Tags          []string         `json:"tags"`
	IsCurrent     pgtype.Bool      `json:"is_current"`
	Jobs          []ExperienceJob  `json:"jobs"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
}

// PaginatedResponse represents a paginated list of experiences with jobs.
type PaginatedResponse struct {
	Experiences []ExperienceWithJobs `json:"experiences"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	TotalPages  int                  `json:"total_pages"`
}

// Service handles experience business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new experience Service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns paginated experiences with nested jobs.
func (s *Service) List(ctx context.Context, page, limit int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	rows, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, fmt.Errorf("listing experiences: %w", err)
	}

	total := int64(0)
	expIDs := make([]pgtype.UUID, 0, len(rows))

	for _, row := range rows {
		expIDs = append(expIDs, row.ID)
		total = row.TotalCount
	}

	// Fetch all jobs for these experiences in one query
	jobs, err := s.repo.GetJobsByExperienceIDs(ctx, expIDs)
	if err != nil {
		return nil, fmt.Errorf("getting jobs: %w", err)
	}

	jobsByExpID := make(map[string][]ExperienceJob)
	for _, j := range jobs {
		idStr := uuidToString(j.ExperienceID)
		jobsByExpID[idStr] = append(jobsByExpID[idStr], j)
	}

	experiences := make([]ExperienceWithJobs, 0, len(rows))
	for _, row := range rows {
		idStr := uuidToString(row.ID)
		expJobs := jobsByExpID[idStr]
		if expJobs == nil {
			expJobs = []ExperienceJob{}
		}

		experiences = append(experiences, ExperienceWithJobs{
			ID:            row.ID,
			Company:       row.Company,
			Location:      row.Location,
			TotalDuration: row.TotalDuration,
			Tags:          row.Tags,
			IsCurrent:     row.IsCurrent,
			Jobs:          expJobs,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &PaginatedResponse{
		Experiences: experiences,
		Total:       total,
		Page:        page,
		TotalPages:  totalPages,
	}, nil
}

// GetByID returns an experience with its nested jobs.
func (s *Service) GetByID(ctx context.Context, id string) (*ExperienceWithJobs, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	exp, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("getting experience: %w", err)
	}

	jobs, err := s.repo.GetJobsByExperienceID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("getting jobs: %w", err)
	}
	if jobs == nil {
		jobs = []ExperienceJob{}
	}

	return &ExperienceWithJobs{
		ID:            exp.ID,
		Company:       exp.Company,
		Location:      exp.Location,
		TotalDuration: exp.TotalDuration,
		Tags:          exp.Tags,
		IsCurrent:     exp.IsCurrent,
		Jobs:          jobs,
		CreatedAt:     exp.CreatedAt,
		UpdatedAt:     exp.UpdatedAt,
	}, nil
}

// Create creates a new experience with jobs and calculates total_duration.
func (s *Service) Create(ctx context.Context, input CreateInput) (*ExperienceWithJobs, error) {
	totalDuration := calculateTotalDuration(input.Jobs)

	// Determine if any job is current
	isCurrent := false
	for _, j := range input.Jobs {
		if j.IsCurrent {
			isCurrent = true
			break
		}
	}

	exp, err := s.repo.CreateExperience(ctx, CreateExperienceParams{
		Company:       input.Company,
		Location:      input.Location,
		TotalDuration: totalDuration,
		Tags:          input.Tags,
		IsCurrent:     pgtype.Bool{Bool: isCurrent, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("creating experience: %w", err)
	}

	jobs := make([]ExperienceJob, 0, len(input.Jobs))
	for _, ji := range input.Jobs {
		job, err := s.repo.CreateJob(ctx, CreateJobParams{
			ExperienceID:   exp.ID,
			Position:       ji.Position,
			EmploymentType: ji.EmploymentType,
			Description:    ji.Description,
			Competency:     ji.Competency,
			StartDate:      ji.StartDate,
			EndDate:        toPgText(ji.EndDate),
			IsCurrent:      pgtype.Bool{Bool: ji.IsCurrent, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("creating job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return &ExperienceWithJobs{
		ID:            exp.ID,
		Company:       exp.Company,
		Location:      exp.Location,
		TotalDuration: exp.TotalDuration,
		Tags:          exp.Tags,
		IsCurrent:     exp.IsCurrent,
		Jobs:          jobs,
		CreatedAt:     exp.CreatedAt,
		UpdatedAt:     exp.UpdatedAt,
	}, nil
}

// Update modifies an existing experience: deletes old jobs, updates experience, inserts new jobs.
func (s *Service) Update(ctx context.Context, input UpdateInput) (*ExperienceWithJobs, error) {
	uid, err := parseUUID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	// Delete old jobs
	if err := s.repo.DeleteJobsByExperienceID(ctx, uid); err != nil {
		return nil, fmt.Errorf("deleting old jobs: %w", err)
	}

	totalDuration := calculateTotalDuration(input.Jobs)
	isCurrent := false
	for _, j := range input.Jobs {
		if j.IsCurrent {
			isCurrent = true
			break
		}
	}

	exp, err := s.repo.UpdateExperience(ctx, UpdateExperienceParams{
		ID:            uid,
		Company:       input.Company,
		Location:      input.Location,
		TotalDuration: totalDuration,
		Tags:          input.Tags,
		IsCurrent:     pgtype.Bool{Bool: isCurrent, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("updating experience: %w", err)
	}

	jobs := make([]ExperienceJob, 0, len(input.Jobs))
	for _, ji := range input.Jobs {
		job, err := s.repo.CreateJob(ctx, CreateJobParams{
			ExperienceID:   uid,
			Position:       ji.Position,
			EmploymentType: ji.EmploymentType,
			Description:    ji.Description,
			Competency:     ji.Competency,
			StartDate:      ji.StartDate,
			EndDate:        toPgText(ji.EndDate),
			IsCurrent:      pgtype.Bool{Bool: ji.IsCurrent, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("creating job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return &ExperienceWithJobs{
		ID:            exp.ID,
		Company:       exp.Company,
		Location:      exp.Location,
		TotalDuration: exp.TotalDuration,
		Tags:          exp.Tags,
		IsCurrent:     exp.IsCurrent,
		Jobs:          jobs,
		CreatedAt:     exp.CreatedAt,
		UpdatedAt:     exp.UpdatedAt,
	}, nil
}

// Delete removes an experience and its jobs (CASCADE).
func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteExperience(ctx, uid)
}

// Count returns the total number of experiences.
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func calculateTotalDuration(jobs []JobInput) string {
	pairs := make([][2]time.Time, 0, len(jobs))
	for _, j := range jobs {
		start, err := time.Parse("2006-01-02", j.StartDate)
		if err != nil {
			continue
		}

		var end time.Time
		if j.IsCurrent || j.EndDate == "" {
			end = time.Now()
		} else {
			end, err = time.Parse("2006-01-02", j.EndDate)
			if err != nil {
				continue
			}
		}

		pairs = append(pairs, [2]time.Time{start, end})
	}

	return utils.CalculateTotalDuration(pairs...)
}
