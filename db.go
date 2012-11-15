package main

import (
	"github.com/simonz05/godis/redis"
)

func opmlFromUUID(uuid string) (opml []byte, err error) {
	r := redis.New(REDIS_ADDR, OPMLFEED_REDIS_DB, REDIS_PASS)
	opml, err = r.Get("OF_" + uuid)
        return
}

func uuidFromShort(shortid string) (uuid []byte, err error) {
	r := redis.New(REDIS_ADDR, OPMLFEED_REDIS_DB, REDIS_PASS)
	uuid, err = r.Get("OF_" + shortid)
        return
}

func shortFromUUID(uuid string) (shortid []byte, err error) {
	r := redis.New(REDIS_ADDR, OPMLFEED_REDIS_DB, REDIS_PASS)
        shortid, err = r.Get("OF_id_" + uuid)
        return
}

func storeClientData(uuid string, clientData []byte) (err error) {
	r := redis.New(REDIS_ADDR, OPMLFEED_REDIS_DB, REDIS_PASS)
        err = r.Set("OF_" + uuid, clientData)
        return
}

func associateUUIDandShortid(shortid string, uuid string) (err error) {
	r := redis.New(REDIS_ADDR, OPMLFEED_REDIS_DB, REDIS_PASS)
        err = r.Set("OF_"+shortid, uuid)
        if err == nil {
                err = r.Set("OF_id_"+uuid, shortid)
        }
        return
}
