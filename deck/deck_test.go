package deck

import (
	"reflect"
	"testing"
)

func TestCard_String(t *testing.T) {
	type fields struct {
		Suit  uint8
		Value uint8
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"Spade:A", fields{Suit: Spade, Value: 1}},
		{"Diamond:J", fields{Suit: Diamond, Value: 11}},
		{"Heart:K", fields{Suit: Heart, Value: 13}},
		{"Joker", fields{Suit: Joker, Value: 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Card{
				Suit:  tt.fields.Suit,
				Value: tt.fields.Value,
			}
			if got := c.String(); got != tt.name {
				t.Errorf("String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts Opts
	}
	tests := []struct {
		name    string
		args    args
		want    []Card
		wantErr bool
	}{
		{
			"Default",
			args{Opts{}},
			[]Card{
				{Spade, 1},
				{Spade, 2},
				{Spade, 3},
				{Spade, 4},
				{Spade, 5},
				{Spade, 6},
				{Spade, 7},
				{Spade, 8},
				{Spade, 9},
				{Spade, 10},
				{Spade, 11},
				{Spade, 12},
				{Spade, 13},
				{Diamond, 1},
				{Diamond, 2},
				{Diamond, 3},
				{Diamond, 4},
				{Diamond, 5},
				{Diamond, 6},
				{Diamond, 7},
				{Diamond, 8},
				{Diamond, 9},
				{Diamond, 10},
				{Diamond, 11},
				{Diamond, 12},
				{Diamond, 13},
				{Heart, 1},
				{Heart, 2},
				{Heart, 3},
				{Heart, 4},
				{Heart, 5},
				{Heart, 6},
				{Heart, 7},
				{Heart, 8},
				{Heart, 9},
				{Heart, 10},
				{Heart, 11},
				{Heart, 12},
				{Heart, 13},
				{Club, 1},
				{Club, 2},
				{Club, 3},
				{Club, 4},
				{Club, 5},
				{Club, 6},
				{Club, 7},
				{Club, 8},
				{Club, 9},
				{Club, 10},
				{Club, 11},
				{Club, 12},
				{Club, 13},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildFilterMap(t *testing.T) {
	type args struct {
		filter []Card
	}
	tests := []struct {
		name    string
		args    args
		want    map[Card]bool
		wantErr bool
	}{
		{
			"One Card",
			args{[]Card{{Spade, 1}}},
			map[Card]bool{Card{Spade, 1}: true},
			false,
		},
		{
			"Suit Wildcard",
			args{[]Card{{Wildcard, 2}}},
			map[Card]bool{
				Card{Spade, 2}:   true,
				Card{Diamond, 2}: true,
				Card{Heart, 2}:   true,
				Card{Club, 2}:    true,
			},
			false,
		},
		{
			"Value Wildcard",
			args{[]Card{{Heart, Wildcard}}},
			map[Card]bool{
				Card{Heart, 1}:  true,
				Card{Heart, 2}:  true,
				Card{Heart, 3}:  true,
				Card{Heart, 4}:  true,
				Card{Heart, 5}:  true,
				Card{Heart, 6}:  true,
				Card{Heart, 7}:  true,
				Card{Heart, 8}:  true,
				Card{Heart, 9}:  true,
				Card{Heart, 10}: true,
				Card{Heart, 11}: true,
				Card{Heart, 12}: true,
				Card{Heart, 13}: true,
			},
			false,
		},
		{
			"Double Wildcard",
			args{[]Card{{Wildcard, Wildcard}}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildFilterMap(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildFilterMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildFilterMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
