package main

import "time"

type ExerciseSet struct {
	Name        string
	PerformedAt time.Time
	Weight      float32
	Unit        string
	Reps        int
	SetCount    int
	Notes       string
	SplitDay    string
	Program     string
	Tags        string
}
