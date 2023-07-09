package users

import (
	"testing"

	"github.com/JakeStrang1/escapehatch/db"
	"gotest.tools/assert"
)

func Test_getSecondDegreeConnections(t *testing.T) {
	tests := []struct {
		name   string
		caller User
		users  []User
		want   []User
	}{
		{
			name: "results are filtered and ranked",
			caller: User{
				DefaultModel: db.DefaultModel{
					IDField: db.IDField{
						ID: "1",
					},
				},
				Followers: []Follower{
					{FollowerUserID: "a"},
					{FollowerUserID: "b"},
					{FollowerUserID: "c"},
					{FollowerUserID: "d"},
				},
			},
			users: []User{
				{
					Following: []Follower{
						{TargetUserID: "g"},
						{TargetUserID: "h"},
						{TargetUserID: "c"},
						{TargetUserID: "d"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "k"},
						{TargetUserID: "l"},
						{TargetUserID: "m"},
						{TargetUserID: "n"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "k"},
						{TargetUserID: "1"},
						{TargetUserID: "n"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "a"},
						{TargetUserID: "b"},
						{TargetUserID: "e"},
						{TargetUserID: "d"},
						{TargetUserID: "f"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "i"},
						{TargetUserID: "b"},
						{TargetUserID: "j"},
					},
				},
			},
			want: []User{
				{
					Following: []Follower{
						{TargetUserID: "k"},
						{TargetUserID: "1"},
						{TargetUserID: "n"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "a"},
						{TargetUserID: "b"},
						{TargetUserID: "e"},
						{TargetUserID: "d"},
						{TargetUserID: "f"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "g"},
						{TargetUserID: "h"},
						{TargetUserID: "c"},
						{TargetUserID: "d"},
					},
				},
				{
					Following: []Follower{
						{TargetUserID: "i"},
						{TargetUserID: "b"},
						{TargetUserID: "j"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSecondDegreeConnections(tt.caller, tt.users)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}
