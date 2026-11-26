package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated"
	"github.com/levisantosp/atm-participa/apps/api/ent/generated/user"
	"github.com/levisantosp/atm-participa/apps/api/internal/db"
	"github.com/levisantosp/atm-participa/pkg/infra/redis"
	"github.com/taskforcesh/bullmq-go"
	"github.com/taskforcesh/bullmq-go/redisbackend"
)

const name = "unban"

type unbanJob struct {
	UserID int64 `json:"userId"`
}

func ScheduleUnban(
	ctx context.Context,
	id int64,
	endsAt time.Time,
) (err error) {
	backend := redisbackend.New(redisbackend.Options{Client: redis.Client})
	queue, err := bullmq.NewQueue(name, &bullmq.QueueOptions{Backend: backend})
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, queue.Close())
	}()

	now := time.Now().UnixMilli()
	delay := max(
		0,
		endsAt.Add(time.Millisecond-time.Nanosecond).UnixMilli()-now,
	)
	_, err = queue.Add(ctx, name, unbanJob{UserID: id}, &bullmq.JobOptions{
		JobID:     fmt.Sprintf("unban-%d-%d", id, endsAt.UnixMicro()),
		Timestamp: now,
		Delay:     &delay,
		Attempts:  new(int64(10)),
		Backoff: &bullmq.Backoff{
			Type:  bullmq.BackoffExponential,
			Delay: (5 * time.Second).Milliseconds(),
		},
		RemoveOnComplete: bullmq.RemoveAll(),
	})

	return err
}

func Unban(ctx context.Context) (err error) {
	backend := redisbackend.New(redisbackend.Options{Client: redis.Client})
	queue, err := bullmq.NewQueue(name, &bullmq.QueueOptions{Backend: backend})
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, queue.Close())
	}()

	if err := queue.SetGlobalConcurrency(ctx, 1); err != nil {
		return err
	}

	if _, err := queue.RemoveJobScheduler(ctx, name); err != nil {
		return err
	}

	users, err := db.Client.User.Query().
		Where(user.BannedEQ(true), user.BanEndsAtNotNil()).
		Select(user.FieldID, user.FieldBanEndsAt).
		All(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		if err := ScheduleUnban(ctx, u.ID, *u.BanEndsAt); err != nil {
			return err
		}
	}

	worker, err := bullmq.NewWorker(
		name,
		func(ctx context.Context, job *bullmq.Job) (any, error) {
			var payload unbanJob
			if err := job.DecodeData(&payload); err != nil {
				return nil, err
			}

			return db.WithTx(ctx, func(tx *generated.Tx) (*struct{}, error) {
				_, err := tx.User.Query().Unique(false).
					Where(user.IDEQ(payload.UserID), func(s *sql.Selector) {
						s.ForUpdate()
					}).
					OnlyID(ctx)
				if generated.IsNotFound(err) {
					return nil, nil
				}

				if err != nil {
					return nil, err
				}

				return nil, tx.User.Update().
					Where(user.IDEQ(payload.UserID), user.BannedEQ(true), user.BanEndsAtLTE(time.Now())).
					ClearBanned().
					ClearBannedAt().
					ClearBanEndsAt().
					Exec(ctx)
			})
		},
		&bullmq.WorkerOptions{
			Backend: backend,
		},
	)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, worker.Close())
	}()

	return worker.Run(ctx)
}
