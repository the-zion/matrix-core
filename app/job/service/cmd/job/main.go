package main

func main() {
	log, err := logInit()
	if err != nil {
		return
	}
	defer log.client.Close(60000)

	config, err := getConfig()
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	news, err := getNews(config)
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	dbClient, err := newDB(config)
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	redisClient, err := newRedis(config)
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	cosCli, err := newCosServiceClient(config)
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	es, err := newElasticsearch(config)
	if err != nil {
		log.SendLog(err.Error())
		return
	}

	setNews(news, cosCli, es, dbClient, redisClient, config, log)
}

func setNews(news []*News, cosCli *CosClient, es *ElasticSearch, db *DB, redis *Redis, config map[string]interface{}, log *Log) {
	for _, n := range news {
		item := convertToMap(n)

		err := getContent(config, item)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}

		err = getImages(config, item, cosCli)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}

		err = cosCli.setNews(item)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}

		err = db.setNews(item)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}

		err = es.setNews(item)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}

		err = redis.setNews(item)
		if err != nil {
			log.SendLog(err.Error())
			continue
		}
		return
	}
}

func convertToMap(news *News) map[string]string {
	return map[string]string{
		"id":        news.Id,
		"cover":     news.Cover,
		"update":    news.Update,
		"title":     news.Title,
		"author":    news.Author,
		"tags":      news.Tags,
		"url":       news.Url,
		"introduce": "",
		"content":   "",
	}
}
