package redis_repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/dtos"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type redisUsersRepoImpl struct {
	redisClient *redis.Client
	logger      *zerolog.Logger
}

type RedisUsersRepository interface {
	GetUsersCtx(ctx context.Context, usersid int) (*dtos.ResDataUserSingle, error)
	SetUsersCtx(ctx context.Context, data *dtos.ResDataUserSingle) error
	DeleteUsersCtx(ctx context.Context, usersid int) error
}

func NewRedisRepoUsers(redisClient *redis.Client, logger *zerolog.Logger) RedisUsersRepository {
	return &redisUsersRepoImpl{
		redisClient: redisClient,
		logger:      logger,
	}
}

func keyUsersGenerator(userid int) string {
	var key = "users:get:"
	return fmt.Sprint(key, userid)
}

func (roi *redisUsersRepoImpl) GetUsersCtx(ctx context.Context, userid int) (r *dtos.ResDataUserSingle, err error) {

	// Get data from Redis
	realKey := keyUsersGenerator(userid)
	log.Info().Msg(fmt.Sprintf("Get keys %s from redis\n", realKey))
	result, err := roi.redisClient.Get(ctx, realKey).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.Wrap(err, "usersRedisRepo.GetUsersCtx.redisClient.Get")
	}

	if err == redis.Nil {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(result), &r); err != nil {
		return nil, errors.Wrap(err, "usersRedisRepo.GetUsersCtx.json.Marshal")
	}

	log.Info().Msg(fmt.Sprintf("Succedd get keys %s from redis\n", realKey))
	return r, nil
}

func (roi *redisUsersRepoImpl) SetUsersCtx(ctx context.Context, data *dtos.ResDataUserSingle) error {
	newBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "usersRedisRepo.SetUsersCtx.json.Marshal")
	}

	log.Info().Msg("set key")
	realKey := keyUsersGenerator(int(data.ID))
	if err = roi.redisClient.Set(ctx, realKey, newBytes, RedisTTL).Err(); err != nil {
		return errors.Wrap(err, "usersRedisRepo.SetUsersCtx.redisClient.set")
	}

	log.Info().Msg(fmt.Sprintf("Set keys %s to redis with TTL %d \n", realKey, RedisTTL))
	return nil
}

func (roi *redisUsersRepoImpl) DeleteUsersCtx(ctx context.Context, userid int) error {
	realKey := keyUsersGenerator(userid)
	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis\n", realKey))
	if err := roi.redisClient.Del(ctx, realKey).Err(); err != nil {
		return errors.Wrap(err, "usersRedisRepo.DeleteUsersCtx.redisClient.Del")
	}

	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis succeed\n", realKey))
	return nil
}
