package exec

import (
	"encoding/base64"
	"fmt"
	"graphql/graphql-starwar/data"
	"graphql/graphql-starwar/model"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/graphql-go/graphql"
)

var StarWarsSchema graphql.Schema

var (
	episodeEnum    *graphql.Enum
	lengthUnitEnum *graphql.Enum

	characterInterface *graphql.Interface

	reviewInputType *graphql.InputObject

	searchResultUnion *graphql.Union

	friendsConnectionType *graphql.Object
	friendsEdgeType       *graphql.Object
	pageInfoType          *graphql.Object
	starShipType          *graphql.Object
	humanType             *graphql.Object
	droidType             *graphql.Object
	reviewType            *graphql.Object
)

func init() {
	episodeEnum = graphql.NewEnum(graphql.EnumConfig{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: graphql.EnumValueConfigMap{
			"NEWHOPE": &graphql.EnumValueConfig{
				Value:       model.EpisodeNewhope,
				Description: "Released in 1977.",
			},
			"EMPIRE": &graphql.EnumValueConfig{
				Value:       model.EpisodeEmpire,
				Description: "Released in 1980.",
			},
			"JEDI": &graphql.EnumValueConfig{
				Value:       model.EpisodeJedi,
				Description: "Released in 1983.",
			},
		},
	})

	lengthUnitEnum = graphql.NewEnum(graphql.EnumConfig{
		Name:        "LengthUnit",
		Description: "Units of height",
		Values: graphql.EnumValueConfigMap{
			"METER": &graphql.EnumValueConfig{
				Value:       model.LengthUnitMeter,
				Description: "The standard unit around the world",
			},
			"FOOT": &graphql.EnumValueConfig{
				Value:       model.LengthUnitFoot,
				Description: "Primarily used in the United States",
			},
		},
	})

	characterInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "Character",
		Description: "A character in the Star Wars Trilogy",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The id of the character.",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the character.",
			},
			"appearsIn": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(episodeEnum))),
				Description: "Which movies they appear in.",
			},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if character, ok := p.Value.(*model.Human); ok {
				if _, ok := data.Humans[character.ID]; ok {
					return humanType
				}
			}
			if character, ok := p.Value.(*model.Droid); ok {
				if _, ok := data.Droids[character.ID]; ok {
					return droidType
				}
			}
			return nil
		},
	})
	characterInterface.AddFieldConfig("friends", &graphql.Field{
		Type:        graphql.NewList(graphql.NewNonNull(characterInterface)),
		Description: "The friends of the character, or an empty list if they have none.",
	})

	friendsEdgeType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "FriendsEdge",
		Description: "An edge object for a character's friends",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "A cursor used for pagination",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.FriendsEdge); ok {
						return human.Cursor, nil
					}
					return nil, nil
				},
			},
			"node": &graphql.Field{
				Type:        characterInterface,
				Description: "The character represented by this friendship edge",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.FriendsEdge); ok {
						return human.Node, nil
					}
					return nil, nil
				},
			},
		},
	})

	pageInfoType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "PageInfo",
		Description: "Information for paginating this connection",
		Fields: graphql.Fields{
			"startCursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "start cursor",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if pageInfo, ok := p.Source.(*model.PageInfo); ok {
						return pageInfo.StartCursor, nil
					}
					return nil, nil
				},
			},
			"endCursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "end cursor",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if pageInfo, ok := p.Source.(*model.PageInfo); ok {
						return pageInfo.EndCursor, nil
					}
					return nil, nil
				},
			},
			"hasNextPage": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "has next page",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if pageInfo, ok := p.Source.(*model.PageInfo); ok {
						return pageInfo.HasNextPage, nil
					}
					return nil, nil
				},
			},
		},
	})

	friendsConnectionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "FriendsConnection",
		Description: "A connection object for a character's friends",
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The total number of friends",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if fc, ok := p.Source.(*model.FriendsConnection); ok {
						return fc.TotalCount, nil
					}
					return 0, nil
				},
			},
			"edges": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(friendsEdgeType)),
				Description: "An edge object for a character's friends",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if fc, ok := p.Source.(*model.FriendsConnection); ok {
						return fc.Edges, nil
					}
					return nil, nil
				},
			},
			"friends": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(characterInterface)),
				Description: "A list of the friends, as a convenience when edges are not needed.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if fc, ok := p.Source.(*model.FriendsConnection); ok {
						return fc.Friends, nil
					}
					return nil, nil
				},
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(pageInfoType),
				Description: "Information for paginating this connection",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if fc, ok := p.Source.(*model.FriendsConnection); ok {
						return fc.PageInfo, nil
					}
					return nil, nil
				},
			},
		},
	})

	characterInterface.AddFieldConfig("friendsConnection", &graphql.Field{
		Type:        graphql.NewNonNull(friendsConnectionType),
		Description: "The friends of the human exposed as a connection with edges",
		Args: graphql.FieldConfigArgument{
			"first": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Height in the preferred unit, default is meters",
			},
			"after": &graphql.ArgumentConfig{
				Type:        graphql.ID,
				Description: "Height in the preferred unit, default is meters",
			},
		},
	})

	starShipType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Starship",
		Description: "star ship",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The ID of the starship",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Starship); ok {
						return review.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the starship",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Starship); ok {
						return review.Name, nil
					}
					return nil, nil
				},
			},
			"length": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Float),
				Description: "Length of the starship, along the longest axis",
				Args: graphql.FieldConfigArgument{
					"unit": &graphql.ArgumentConfig{
						Type:        lengthUnitEnum,
						Description: "Height in the preferred unit, default is meters",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if starship, ok := p.Source.(*model.Starship); ok {
						if unit, ok := p.Args["unit"].(model.LengthUnit); ok {
							if unit == model.LengthUnitFoot {
								return starship.Length * 3.28084, nil
							}
						}
						return starship.Length, nil
					}
					return 0, nil
				},
			},
			"history": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int))))),
				Description: "coordinates tracking this ship",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Starship); ok {
						return review.History, nil
					}
					return nil, nil
				},
			},
		},
	})

	humanType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Human",
		Description: "A humanoid creature in the Star Wars universe.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The id of the human.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						return human.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the human.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						return human.Name, nil
					}
					return nil, nil
				},
			},
			"height": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Float),
				Description: "Height in the preferred unit, default is meters",
				Args: graphql.FieldConfigArgument{
					"unit": &graphql.ArgumentConfig{
						Type:        lengthUnitEnum,
						Description: "Height in the preferred unit, default is meters",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						if unit, ok := p.Args["unit"].(model.LengthUnit); ok {
							if unit == model.LengthUnitFoot {
								return human.Height * 3.28084, nil
							}
						}
						return human.Height, nil
					}
					return nil, nil
				},
			},
			"mass": &graphql.Field{
				Type:        graphql.Float,
				Description: "Mass in kilograms, or null if unknown",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); !ok {
						return human.Mass, nil
					}
					return nil, nil
				},
			},
			"friends": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(characterInterface)),
				Description: "The friends of the human, or an empty list if they have none.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						return human.Friends, nil
					}
					return []interface{}{}, nil
				},
			},
			"friendsConnection": &graphql.Field{
				Type:        graphql.NewNonNull(friendsConnectionType),
				Description: "The friends of the human exposed as a connection with edges",
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Height in the preferred unit, default is meters",
					},
					"after": &graphql.ArgumentConfig{
						Type:        graphql.ID,
						Description: "Height in the preferred unit, default is meters",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveFriendConnection(p)
				},
			},
			"appearsIn": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(episodeEnum))),
				Description: "Which movies they appear in.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						return human.AppearsIn, nil
					}
					return nil, nil
				},
			},
			"starships": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(starShipType)),
				Description: "A list of starships this person has piloted, or an empty list if none",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if human, ok := p.Source.(*model.Human); ok {
						return human.Starships, nil
					}
					return nil, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			characterInterface,
		},
	})

	droidType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Droid",
		Description: "A mechanical creature in the Star Wars universe.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.ID),
				Description: "The id of the droid.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if droid, ok := p.Source.(*model.Droid); ok {
						return droid.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the droid.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if droid, ok := p.Source.(*model.Droid); ok {
						return droid.Name, nil
					}
					return nil, nil
				},
			},
			"friends": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(characterInterface)),
				Description: "The friends of the droid, or an empty list if they have none.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if droid, ok := p.Source.(*model.Droid); ok {
						return droid.Friends, nil
					}
					return []interface{}{}, nil
				},
			},
			"friendsConnection": &graphql.Field{
				Type:        graphql.NewNonNull(friendsConnectionType),
				Description: "The friends of the human exposed as a connection with edges",
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Height in the preferred unit, default is meters",
					},
					"after": &graphql.ArgumentConfig{
						Type:        graphql.ID,
						Description: "Height in the preferred unit, default is meters",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolveFriendConnection(p)
				},
			},
			"appearsIn": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(episodeEnum))),
				Description: "Which movies they appear in.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if droid, ok := p.Source.(*model.Droid); ok {
						return droid.AppearsIn, nil
					}
					return nil, nil
				},
			},
			"primaryFunction": &graphql.Field{
				Type:        graphql.String,
				Description: "The primary function of the droid.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if droid, ok := p.Source.(*model.Droid); ok {
						return droid.PrimaryFunction, nil
					}
					return nil, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			characterInterface,
		},
	})

	searchResultUnion = graphql.NewUnion(graphql.UnionConfig{
		Name:        "SearchResult",
		Description: "search result",
		Types: []*graphql.Object{
			humanType,
			droidType,
			starShipType,
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if _, ok := p.Value.(*model.Human); ok {
				return humanType
			}
			if _, ok := p.Value.(*model.Droid); ok {
				return droidType
			}
			if _, ok := p.Value.(*model.Starship); ok {
				return starShipType
			}
			return nil
		},
	})

	reviewType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Review",
		Description: "Represents a review for a movie",
		Fields: graphql.Fields{
			"stars": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The number of stars this review gave, 1-5",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Review); ok {
						return review.Stars, nil
					}
					return 0, nil
				},
			},
			"commentary": &graphql.Field{
				Type:        graphql.String,
				Description: "Comment about the movie",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Review); ok {
						return review.Commentary, nil
					}
					return nil, nil
				},
			},
			"time": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "when the review was posted",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if review, ok := p.Source.(*model.Review); ok {
						return review.Time, nil
					}
					return nil, nil
				},
			},
		},
	})

	reviewInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "ReviewInput",
		Description: "The input object sent when someone is creating a new review",
		Fields: graphql.InputObjectConfigFieldMap{
			"stars": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The number of stars this review gave, 1-5",
			},
			"commentary": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Comment about the movie",
			},
			"time": &graphql.InputObjectFieldConfig{
				Type:        graphql.DateTime,
				Description: "when the review was posted",
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: characterInterface,
				Args: graphql.FieldConfigArgument{
					"episode": &graphql.ArgumentConfig{
						Description: "If omitted, returns the hero of the whole saga. If " +
							"provided, returns the hero of that particular episode.",
						Type: episodeEnum,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if p.Args["episode"] != nil {
						if episode, ok := p.Args["episode"].(model.Episode); ok {
							if episode == model.EpisodeEmpire {
								return data.Humans["1000"], nil
							}
						}
					}
					return data.Droids["2001"], nil
				},
			},
			"reviews": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(reviewType))),
				Args: graphql.FieldConfigArgument{
					"episode": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(episodeEnum),
						Description: "The episodes in the Star Wars trilogy",
					},
					"since": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.DateTime),
						Description: "when the review was posted",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					episode, ok := p.Args["episode"].(model.Episode)
					if !ok {
						return nil, nil
					}
					since, ok := p.Args["since"].(time.Time)
					if !ok {
						return nil, nil
					}
					review, ok := data.Reviews[episode]
					if !ok {
						return nil, nil
					}
					var filtered []*model.Review
					for _, r := range review {
						if r.Time.After(since) {
							filtered = append(filtered, r)
						}
					}
					return filtered, nil
				},
			},
			"search": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(searchResultUnion))),
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "text of search",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var l []model.SearchResult
					text := p.Args["text"].(string)
					for _, human := range data.Humans {
						if strings.Contains(human.Name, text) {
							l = append(l, human)
						}
					}
					for _, droid := range data.Droids {
						if strings.Contains(droid.Name, text) {
							l = append(l, droid)
						}
					}
					for _, starship := range data.Starships {
						if strings.Contains(starship.Name, text) {
							l = append(l, starship)
						}
					}
					return l, nil
				},
			},
			"character": &graphql.Field{
				Type: characterInterface,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "the id of character",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(string); ok {
						if h, ok := data.Humans[id]; ok {
							return h, nil
						}
						if d, ok := data.Droids[id]; ok {
							return d, nil
						}
					}
					return nil, nil
				},
			},
			"starship": &graphql.Field{
				Type: starShipType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "the id of star ship",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(string); ok {
						if s, ok := data.Starships[id]; ok {
							return s, nil
						}
					}
					return nil, nil
				},
			},
			"human": &graphql.Field{
				Type: humanType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the human",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(string); ok {
						return data.Humans[id], nil
					}
					return nil, nil
				},
			},
			"droid": &graphql.Field{
				Type: droidType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the droid",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if id, ok := p.Args["id"].(string); ok {
						return data.Droids[id], nil
					}
					return nil, nil
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createReview": &graphql.Field{
				Type: reviewType,
				Args: graphql.FieldConfigArgument{
					"episode": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(episodeEnum),
					},
					"review": &graphql.ArgumentConfig{
						Type: reviewInputType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					episode, ok := p.Args["episode"].(model.Episode)
					if !ok {
						return nil, nil
					}
					input, ok := p.Args["review"].(map[string]interface{})
					if !ok {
						return nil, nil
					}
					review := &model.Review{
						Stars: input["stars"].(int),
					}
					if input["commentary"] != nil {
						review.Commentary = proto.String(input["commentary"].(string))
					}
					if input["time"] != nil {
						if t, ok := input["time"].(time.Time); ok {
							review.Time = &t
						}
					}
					data.Reviews[episode] = append(data.Reviews[episode], review)
					return review, nil
				},
			},
		},
	})

	StarWarsSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})

}

func resolveFriendConnection(p graphql.ResolveParams) (interface{}, error) {
	limit := 0
	if p.Args["first"] != nil {
		if t, ok := p.Args["first"].(int); ok {
			limit = t
		}
	}

	from := 0
	if p.Args["after"] != nil {
		if t, ok := p.Args["after"].(string); ok {
			cursor, _ := decodeCursor(t)
			from = cursor
		}
	}

	var ids []model.Character
	if human, ok := p.Source.(*model.Human); ok {
		ids = append(ids, human.Friends...)
	}
	if droid, ok := p.Source.(*model.Droid); ok {
		ids = append(ids, droid.Friends...)
	}

	if from >= len(ids) {
		return &model.FriendsConnection{
			TotalCount: len(ids),
			PageInfo:   &model.PageInfo{HasNextPage: false},
		}, nil
	}

	to := len(ids)
	if limit != 0 {
		to = from + limit
		if to > len(ids) {
			to = len(ids)
		}
	}
	var edges []*model.FriendsEdge
	for i := from; i < to && i < len(ids); i++ {
		edge := &model.FriendsEdge{
			Cursor: encodeCursor(i + 1),
			Node:   ids[i],
		}
		edges = append(edges, edge)
	}
	friendsConnection := &model.FriendsConnection{
		TotalCount: len(ids),
		Friends:    ids[from:to],
		PageInfo: &model.PageInfo{
			StartCursor: encodeCursor(from + 1),
			EndCursor:   encodeCursor(to + 1),
			HasNextPage: to >= len(ids),
		},
		Edges: edges,
	}
	return friendsConnection, nil
}

func decodeCursor(s string) (int, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
	if err != nil {
		return 0, err
	}
	return i, err
}

func encodeCursor(i int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}
