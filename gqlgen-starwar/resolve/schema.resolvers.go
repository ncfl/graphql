package resolve

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"graphql/gqlgen-starwar/generated"
	"graphql/gqlgen-starwar/model"
	"strconv"
	"strings"
	"time"
)

type droidResolver struct {
	*Resolver
}

func (r *droidResolver) Friends(ctx context.Context, obj *model.Droid) ([]model.Character, error) {
	return r.resolveCharacters(ctx, obj.Friends)
}

func (r *droidResolver) FriendsConnection(ctx context.Context, obj *model.Droid, first *int, after *string) (*model.FriendsConnection, error) {
	return r.resolveFriendConnection(ctx, obj.Friends, first, after)
}

type friendsConnectionResolver struct {
	*Resolver
}

func (r *friendsConnectionResolver) Edges(ctx context.Context, obj *model.FriendsConnection) ([]*model.FriendsEdge, error) {
	friends, err := r.resolveCharacters(ctx, obj.Friends)
	if err != nil {
		return nil, err
	}

	to, _ := decodeCursor(obj.PageInfo.EndCursor)
	from, _ := decodeCursor(obj.PageInfo.StartCursor)

	edges := make([]*model.FriendsEdge, to-from)
	for i := range edges {
		edges[i] = &model.FriendsEdge{
			Cursor: encodeCursor(from),
			Node:   friends[i],
		}
	}
	return edges, nil
}

func (r *friendsConnectionResolver) Friends(ctx context.Context, obj *model.FriendsConnection) ([]model.Character, error) {
	return r.resolveCharacters(ctx, obj.Friends)
}

type humanResolver struct {
	*Resolver
}

func (r *humanResolver) Height(ctx context.Context, obj *model.Human, unit *model.LengthUnit) (float64, error) {
	if unit == nil || *unit == model.LengthUnitMeter {
		return obj.Height, nil
	}

	if *unit == model.LengthUnitFoot {
		return obj.Height * 3.28084, nil
	}

	return obj.Height, nil
}

func (r *humanResolver) Friends(ctx context.Context, obj *model.Human) ([]model.Character, error) {
	return r.resolveCharacters(ctx, obj.Friends)
}

func (r *humanResolver) FriendsConnection(ctx context.Context, obj *model.Human, first *int, after *string) (*model.FriendsConnection, error) {
	return r.resolveFriendConnection(ctx, obj.Friends, first, after)
}

func (r *humanResolver) Starships(ctx context.Context, obj *model.Human) ([]*model.Starship, error) {
	var result []*model.Starship
	for _, id := range obj.Starships {
		char, err := r.Query().Starship(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		if char != nil {
			result = append(result, char)
		}
	}
	return result, nil
}

func (r *mutationResolver) CreateReview(ctx context.Context, episode model.Episode, review model.ReviewInput) (*model.Review, error) {
	now := time.Now()
	time.Sleep(1 * time.Second)

	reviewRes := model.Review{}
	reviewRes.Commentary = review.Commentary
	reviewRes.Stars = review.Stars
	reviewRes.Time = &now
	r.reviews[episode] = append(r.reviews[episode], &reviewRes)
	return &reviewRes, nil
}

func (r *queryResolver) Hero(ctx context.Context, episode *model.Episode) (model.Character, error) {
	if *episode == model.EpisodeEmpire {
		return r.humans["1000"], nil
	}
	return r.droid["2001"], nil
}

func (r *queryResolver) Reviews(ctx context.Context, episode model.Episode, since *time.Time) ([]*model.Review, error) {
	if since == nil {
		return r.reviews[episode], nil
	}

	var filtered []*model.Review
	for _, rev := range r.reviews[episode] {
		if rev.Time.After(*since) {
			filtered = append(filtered, rev)
		}
	}
	return filtered, nil
}

func (r *queryResolver) Search(ctx context.Context, text string) ([]model.SearchResult, error) {
	var l []model.SearchResult
	for _, h := range r.humans {
		if strings.Contains(h.Name, text) {
			l = append(l, h)
		}
	}
	for _, d := range r.droid {
		if strings.Contains(d.Name, text) {
			l = append(l, d)
		}
	}
	for _, s := range r.starships {
		if strings.Contains(s.Name, text) {
			l = append(l, s)
		}
	}
	return l, nil
}

func (r *queryResolver) Character(ctx context.Context, id string) (model.Character, error) {
	if h, ok := r.humans[id]; ok {
		return &h, nil
	}
	if d, ok := r.droid[id]; ok {
		return &d, nil
	}
	return nil, nil
}

func (r *queryResolver) Droid(ctx context.Context, id string) (*model.Droid, error) {
	if d, ok := r.droid[id]; ok {
		return &d, nil
	}
	return nil, nil
}

func (r *queryResolver) Human(ctx context.Context, id string) (*model.Human, error) {
	if h, ok := r.humans[id]; ok {
		return &h, nil
	}
	return nil, nil
}

func (r *queryResolver) Starship(ctx context.Context, id string) (*model.Starship, error) {
	if s, ok := r.starships[id]; ok {
		return &s, nil
	}
	return nil, nil
}

func (r *starshipResolver) Length(ctx context.Context, obj *model.Starship, unit *model.LengthUnit) (float64, error) {
	switch *unit {
	case model.LengthUnitMeter, "":
		return obj.Length, nil
	case model.LengthUnitFoot:
		return obj.Length * 3.28084, nil
	default:
		return 0, errors.New("invalid unit")
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Human returns generated.HumanResolver implementation.
func (r *Resolver) Human() generated.HumanResolver { return &humanResolver{r} }

// Droid returns generated.DroidResolver implementation.
func (r *Resolver) Droid() generated.DroidResolver { return &droidResolver{r} }

// FriendsConnection returns generated.FriendsConnectionResolver implementation.
func (r *Resolver) FriendsConnection() generated.FriendsConnectionResolver {
	return &friendsConnectionResolver{r}
}

// Starship returns generated.StarshipResolver implementation.
func (r *Resolver) Starship() generated.StarshipResolver { return &starshipResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

type starshipResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *Resolver) resolveCharacters(ctx context.Context, ids []model.Character) ([]model.Character, error) {
	result := make([]model.Character, len(ids))
	for i, id := range ids {
		var realId string
		if human, ok := id.(model.Human); ok {
			realId = human.ID
		}
		if droid, ok := id.(model.Droid); ok {
			realId = droid.ID
		}
		char, err := r.Query().Character(ctx, realId)
		if err != nil {
			return nil, err
		}
		result[i] = char
	}
	return result, nil
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
func (r *Resolver) resolveFriendConnection(_ context.Context, ids []model.Character, limit *int, begin *string) (*model.FriendsConnection, error) {
	from := 0
	if begin != nil {
		cursor, _ := decodeCursor(*begin)
		from = cursor
	}

	if from >= len(ids) {
		return &model.FriendsConnection{
			TotalCount: len(ids),
			PageInfo:   &model.PageInfo{HasNextPage: false},
		}, nil
	}

	to := len(ids)
	if limit != nil {
		to = from + *limit
		if to > len(ids) {
			to = len(ids)
		}
	}

	return &model.FriendsConnection{
		Friends:    ids[from:to],
		TotalCount: len(ids),
		PageInfo: &model.PageInfo{
			StartCursor: encodeCursor(from + 1),
			EndCursor:   encodeCursor(to + 1),
			HasNextPage: to >= len(ids),
		},
	}, nil
}
