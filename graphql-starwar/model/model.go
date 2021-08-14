package model

import (
	"time"
)

type Character interface {
	isCharacter()
}

type SearchResult interface {
	IsSearchResult()
}

type Droid struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Friends           []Character        `json:"friends"`
	FriendsConnection *FriendsConnection `json:"friendsConnection"`
	AppearsIn         []Episode          `json:"appearsIn"`
	PrimaryFunction   *string            `json:"primaryFunction"`
}

func (Droid) isCharacter()    {}
func (Droid) IsSearchResult() {}

type FriendsConnection struct {
	TotalCount int            `json:"totalCount"`
	Edges      []*FriendsEdge `json:"edges"`
	Friends    []Character    `json:"friends"`
	PageInfo   *PageInfo      `json:"pageInfo"`
}

type FriendsEdge struct {
	Cursor string    `json:"cursor"`
	Node   Character `json:"node"`
}

type Human struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Height            float64            `json:"height"`
	Mass              *float64           `json:"mass"`
	Friends           []Character        `json:"friends"`
	FriendsConnection *FriendsConnection `json:"friendsConnection"`
	AppearsIn         []Episode          `json:"appearsIn"`
	Starships         []*Starship        `json:"starships"`
}

func (Human) isCharacter()    {}
func (Human) IsSearchResult() {}

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type Review struct {
	Stars      int        `json:"stars"`
	Commentary *string    `json:"commentary"`
	Time       *time.Time `json:"time"`
}

type ReviewInput struct {
	Stars      int        `json:"stars"`
	Commentary *string    `json:"commentary"`
	Time       *time.Time `json:"time"`
}

type Starship struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Length  float64 `json:"length"`
	History [][]int `json:"history"`
}

func (Starship) IsSearchResult() {}

type Episode string

const (
	EpisodeNewhope Episode = "NEWHOPE"
	EpisodeEmpire  Episode = "EMPIRE"
	EpisodeJedi    Episode = "JEDI"
)

var AllEpisode = []Episode{
	EpisodeNewhope,
	EpisodeEmpire,
	EpisodeJedi,
}

type LengthUnit string

const (
	LengthUnitMeter LengthUnit = "METER"
	LengthUnitFoot  LengthUnit = "FOOT"
)

var AllLengthUnit = []LengthUnit{
	LengthUnitMeter,
	LengthUnitFoot,
}
