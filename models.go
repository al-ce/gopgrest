package main

import "time"

type ExerciseSet struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PerformedAt time.Time `json:"performed_at"`
	Weight      float32   `json:"weight"`
	Unit        string    `json:"unit"`
	Reps        int       `json:"reps"`
	SetCount    int       `json:"set_count"`
	Notes       string    `json:"notes"`
	SplitDay    string    `json:"split_day"`
	Program     string    `json:"program"`
	Tags        string    `json:"tags"`
}
