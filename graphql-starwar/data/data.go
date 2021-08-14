package data

import (
	"graphql/graphql-starwar/model"

	"github.com/golang/protobuf/proto"
)

var (
	Humans    map[string]*model.Human
	Droids    map[string]*model.Droid
	Starships map[string]*model.Starship
	Reviews   map[model.Episode][]*model.Review
)

func init() {
	Luke := &model.Human{
		ID:        "1000",
		Name:      "Luke Skywalker",
		AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		Height:    1.72,
		Mass:      proto.Float64(77),
	}
	Darth := &model.Human{
		ID:        "1001",
		Name:      "Darth Vader",
		AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		Height:    2.02,
		Mass:      proto.Float64(136),
	}
	Han := &model.Human{
		ID:        "1002",
		Name:      "Han Solo",
		AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		Height:    1.8,
		Mass:      proto.Float64(80),
	}
	Leia := &model.Human{
		ID:        "1003",
		Name:      "Leia Organa",
		AppearsIn: []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		Height:    1.5,
		Mass:      proto.Float64(49),
	}
	Wilhuff := &model.Human{
		ID:        "1004",
		Name:      "Wilhuff Tarkin",
		AppearsIn: []model.Episode{model.EpisodeNewhope},
		Height:    1.8,
		Mass:      proto.Float64(0),
	}
	C_3PO := &model.Droid{
		ID:              "2000",
		Name:            "C-3PO",
		AppearsIn:       []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		PrimaryFunction: proto.String("Protocol"),
	}
	R2_D2 := &model.Droid{
		ID:              "2001",
		Name:            "R2-D2",
		AppearsIn:       []model.Episode{model.EpisodeNewhope, model.EpisodeEmpire, model.EpisodeJedi},
		PrimaryFunction: proto.String("Astromech"),
	}

	Millennium := &model.Starship{
		ID:   "3000",
		Name: "Millennium Falcon",
		History: [][]int{
			{1, 2},
			{4, 5},
			{1, 2},
			{3, 2},
		},
		Length: 34.37,
	}
	X_Wing := &model.Starship{
		ID:   "3001",
		Name: "X-Wing",
		History: [][]int{
			{6, 4},
			{3, 2},
			{2, 3},
			{5, 1},
		},
		Length: 12.5,
	}
	TIE := &model.Starship{
		ID:   "3002",
		Name: "TIE Advanced x1",
		History: [][]int{
			{3, 2},
			{7, 2},
			{6, 4},
			{3, 2},
		},
		Length: 9.2,
	}
	Imperial := &model.Starship{
		ID:   "3003",
		Name: "Imperial shuttle",
		History: [][]int{
			{1, 7},
			{3, 5},
			{5, 3},
			{7, 1},
		},
		Length: 20,
	}

	Luke.Friends = []model.Character{Han, Leia, C_3PO, R2_D2}
	Darth.Friends = []model.Character{Wilhuff}
	Han.Friends = []model.Character{Luke, Leia, R2_D2}
	Leia.Friends = []model.Character{Luke, Han, C_3PO, R2_D2}
	Wilhuff.Friends = []model.Character{Darth}
	C_3PO.Friends = []model.Character{Luke, Han, Leia, R2_D2}
	R2_D2.Friends = []model.Character{Luke, Han, Leia}

	Luke.Starships = []*model.Starship{X_Wing, Imperial}
	Darth.Starships = []*model.Starship{TIE}
	Han.Starships = []*model.Starship{Millennium, Imperial}

	Humans = map[string]*model.Human{
		Luke.ID:    Luke,
		Darth.ID:   Darth,
		Han.ID:     Han,
		Leia.ID:    Leia,
		Wilhuff.ID: Wilhuff,
	}

	Droids = map[string]*model.Droid{
		C_3PO.ID: C_3PO,
		R2_D2.ID: R2_D2,
	}

	Starships = map[string]*model.Starship{
		Millennium.ID: Millennium,
		X_Wing.ID:     X_Wing,
		TIE.ID:        TIE,
		Imperial.ID:   Imperial,
	}

	Reviews = map[model.Episode][]*model.Review{}

}
