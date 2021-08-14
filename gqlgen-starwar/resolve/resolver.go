package resolve

import (
	"graphql/gqlgen-starwar/generated"
	"graphql/gqlgen-starwar/model"

	"github.com/golang/protobuf/proto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	humans    map[string]model.Human
	droid     map[string]model.Droid
	starships map[string]model.Starship
	reviews   map[model.Episode][]*model.Review
}

func NewResolver() generated.Config {
	r := Resolver{}
	r.humans = map[string]model.Human{
		"1000": {
			ID:        "1000",
			Name:      "Luke Skywalker",
			Friends:   []model.Character{model.Human{ID: "1002"}, model.Human{ID: "1003"}, model.Droid{ID: "2000"}, model.Droid{ID: "2001"}},
			AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			Height:    1.72,
			Mass:      proto.Float64(77),
			Starships: []*model.Starship{{ID: "3001"}, {ID: "3003"}},
		},
		"1001": {
			ID:        "1001",
			Name:      "Darth Vader",
			Friends:   []model.Character{model.Human{ID: "1004"}},
			AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			Height:    2.02,
			Mass:      proto.Float64(136),
			Starships: []*model.Starship{{ID: "3002"}},
		},
		"1002": {
			ID:        "1002",
			Name:      "Han Solo",
			Friends:   []model.Character{model.Human{ID: "1000"}, model.Human{ID: "1003"}, model.Droid{ID: "2001"}},
			AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			Height:    1.8,
			Mass:      proto.Float64(80),
			Starships: []*model.Starship{{ID: "3000"}, {ID: "3003"}},
		},
		"1003": {
			ID:        "1003",
			Name:      "Leia Organa",
			Friends:   []model.Character{model.Human{ID: "1000"}, model.Human{ID: "1002"}, model.Droid{ID: "2000"}, model.Droid{ID: "2001"}},
			AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			Height:    1.5,
			Mass:      proto.Float64(49),
		},
		"1004": {
			ID:        "1004",
			Name:      "Wilhuff Tarkin",
			Friends:   []model.Character{model.Human{ID: "1001"}},
			AppearsIn: []model.Episode{model.EpisodeNewhope},
			Height:    1.8,
			Mass:      proto.Float64(0),
		},
	}

	r.droid = map[string]model.Droid{
		"2000": {
			ID:              "2000",
			Name:            "C-3PO",
			Friends:         []model.Character{model.Human{ID: "1000"}, model.Human{ID: "1002"}, model.Human{ID: "1003"}, model.Droid{ID: "2001"}},
			AppearsIn:       []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			PrimaryFunction: proto.String("Protocol"),
		},
		"2001": {
			ID:              "2001",
			Name:            "R2-D2",
			Friends:         []model.Character{model.Human{ID: "1000"}, model.Human{ID: "1002"}, model.Human{ID: "1003"}},
			AppearsIn:       []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
			PrimaryFunction: proto.String("Astromech"),
		},
	}

	r.starships = map[string]model.Starship{
		"3000": {
			ID:   "3000",
			Name: "Millennium Falcon",
			History: [][]int{
				{1, 2},
				{4, 5},
				{1, 2},
				{3, 2},
			},
			Length: 34.37,
		},
		"3001": {
			ID:   "3001",
			Name: "X-Wing",
			History: [][]int{
				{6, 4},
				{3, 2},
				{2, 3},
				{5, 1},
			},
			Length: 12.5,
		},
		"3002": {
			ID:   "3002",
			Name: "TIE Advanced x1",
			History: [][]int{
				{3, 2},
				{7, 2},
				{6, 4},
				{3, 2},
			},
			Length: 9.2,
		},
		"3003": {
			ID:   "3003",
			Name: "Imperial shuttle",
			History: [][]int{
				{1, 7},
				{3, 5},
				{5, 3},
				{7, 1},
			},
			Length: 20,
		},
	}

	r.reviews = map[model.Episode][]*model.Review{}

	return generated.Config{
		Resolvers: &r,
	}
}
