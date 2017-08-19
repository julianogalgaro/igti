package mappings

import "errors"

func GetIndexMapping(indexName string) ([]byte, error) {
	mappings := make(map[string][]byte)

	mappings["twitter"] = []byte(`
								  {
									"settings": {
								    "number_of_shards": 1
								  },
								  "mappings": {
								    "tweet": {
								      "_all": {
								        "enabled": false
								      },
								        "properties": {
								            "text": {
								                "type": "string"
								            },
								            "classification":{
								                "type": "string",
								                "index": "not_analyzed"  
								            },
								            "classificationDateStr":{
								                "type": "date",
								                "format": "yyyy-MM-dd HH:mm:ss"
								            },
								            "classificationPredict":{
								                "type": "string",
								                "index": "not_analyzed"
								            },
								            "classificationPredictDateStr":{
								                "type": "date",
								                "format": "yyyy-MM-dd HH:mm:ss"
								            },
								            "classificationPredictRate":{
								                "type": "double"
								            },
								            "coordinates": {
								                "properties":{
								                    "coordinates": {
								                        "type": "geo_point"
								                    },
								                    "type":{
								                        "type": "string"
								                    }
								                }
								            },
								            "createdat": {
								                "format": "EEE MMM dd HH:mm:ss Z YYYY",
								                "type": "date"
								            },
								            "retweetcount":{
								                "type": "long"
								            },
								            "retweeted":{
								                "type": "boolean"
								            },
								            "id_str": {
								                "type": "string",
								                "index": "not_analyzed"
								            },
								            "lang": {
								                "type": "string"
								            },
								            "entities": {
								                "properties": {
								                    "hashtags": {
								                        "properties": {
								                            "text": {
								                                "type": "string",
								                                "fields":{
								                                    "raw":{
								                                        "type": "string",
								                                        "index": "not_analyzed"
								                                    }
								                                }
								                            }
								                        }
								                    },
								                    "urls": {
								                        "properties": {
								                            "display_url": {
								                                "type": "string"
								                            },
								                            "expanded_url": {
								                                "index": "not_analyzed",
								                                "type": "string"
								                            },
								                            "url": {
								                                "type": "string"
								                            }
								                        }
								                    },
								                    "usermentions": {
								                        "properties": {
								                            "indices": {
								                                "type": "long"
								                            },
								                            "screenname": {
								                                "type": "string",
								                                "fields":{
								                                    "raw":{
								                                        "type": "string",
								                                        "index": "not_analyzed"
								                                    }
								                                }
								                            },
								                            "name": {
								                                "type": "string"
								                            }  
								                        }
								                    }
								                }
								            },
								            "place": {
								                "properties": {
								                    "bounding_box": {
								                        "type": "geo_shape",
								                        "tree": "quadtree",
								                        "precision": "1m",
								                        "coerce": true
								                    },
								                    "country": {
								                        "index": "not_analyzed",
								                        "type": "string"
								                    },
								                    "country_code": {
								                        "type": "string"
								                    },
								                    "full_name": {
								                        "index": "not_analyzed",
								                        "type": "string"
								                    },
								                    "id": {
								                        "type": "string"
								                    },
								                    "name": {
								                        "index": "not_analyzed",
								                        "type": "string"
								                    },
								                    "place_type": {
								                        "type": "string"
								                    },
								                    "url": {
								                        "type": "string"
								                    }
								                }
								            },
								            "user": {
								                "properties": {
								                    "createdat": {
								                        "format": "EEE MMM dd HH:mm:ss Z YYYY",
								                        "type": "date"
								                    },
								                    "description": {
								                        "type": "string"
								                    },
								                    "favouritescount": {
								                        "type": "long"
								                    },
								                    "followerscount": {
								                        "type": "long"
								                    },
								                    "friendscount": {
								                        "type": "long"
								                    },
								                    "id": {
								                        "type": "long"
								                    },
								                    "idstr": {
								                        "type": "string"
								                    },
								                    "lang": {
								                        "type": "string"
								                    },
								                    "location": {
								                        "index": "not_analyzed",
								                        "type": "string"
								                    },
								                    "name": {
								                        "index": "not_analyzed",
								                        "type": "string"
								                    },
								                    "profileimageurl": {
								                        "type": "string"
								                    },
								                    "screenname": {
								                        "type": "string",
								                        "fields":{
						                                    "raw":{
						                                        "type": "string",
						                                        "index": "not_analyzed"
						                                    }
						                                }
								                    },
								                    "statusescount": {
								                        "type": "long"
								                    },
								                    "timezone": {
								                        "type": "string"
								                    },
								                    "url": {
								                        "type": "string"
								                    }
								                }
								            }
								        }
								    }
								}
								}`,
	)

	mapping, exists := mappings[indexName]
	if !exists {
		return nil, errors.New("Mapping not found to index [" + indexName + "] ")
	}
	return mapping, nil
}
