package main

import (
	"fmt"
	"strconv"

	"github.com/khareyash05/redis-packet-decoder/models"
)

func ParseRedisPackets(buf string) ([]models.RedisBodyType, error) {
	prs := NewParser(buf)
	vals, err := prs.ParseAll()
	if err != nil {
		return nil, err
	}
	bodies := make([]models.RedisBodyType, len(vals))
	for i, v := range vals {
		b, err := toRedisBody(v)
		if err != nil {
			return nil, err
		}
		bodies[i] = b
	}
	return bodies, nil
}

func toRedisBody(v RespValue) (models.RedisBodyType, error) {
	switch t := v.(type) {
	case []interface{}:
		elems := make([]models.RedisBodyType, len(t))
		for i, e := range t {
			b, err := toRedisBody(e)
			if err != nil {
				return models.RedisBodyType{}, err
			}
			elems[i] = b
		}
		return models.RedisBodyType{Type: "array", Size: len(elems), Data: elems}, nil

	case map[string]interface{}:
		entries := make([]models.RedisMapBody, 0, len(t))
		for k, v2 := range t {
			keyElem := models.RedisElement{Length: len(k), Value: k}
			switch v3 := v2.(type) {
			case string:
				valElem := models.RedisElement{Length: len(v3), Value: v3}
				entries = append(entries, models.RedisMapBody{Key: keyElem, Value: valElem})
			case int64:
				s := strconv.FormatInt(v3, 10)
				valElem := models.RedisElement{Length: len(s), Value: v3}
				entries = append(entries, models.RedisMapBody{Key: keyElem, Value: valElem})
			default:
				nested, err := toRedisBody(v3)
				if err != nil {
					return models.RedisBodyType{}, err
				}
				valElem := models.RedisElement{Length: 0, Value: nested}
				entries = append(entries, models.RedisMapBody{Key: keyElem, Value: valElem})
			}
		}
		return models.RedisBodyType{Type: "map", Size: len(entries), Data: entries}, nil

	case string:
		return models.RedisBodyType{Type: "string", Size: len(t), Data: t}, nil

	case int64:
		s := strconv.FormatInt(t, 10)
		return models.RedisBodyType{Type: "integer", Size: len(s), Data: t}, nil

	default:
		return models.RedisBodyType{}, fmt.Errorf("unsupported type %T", v)
	}
}
