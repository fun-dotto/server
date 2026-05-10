package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/fun-dotto/server/internal/shared/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&model.Announcement{},
		&model.CancelledClass{},
		&model.CourseRegistration{},
		&model.FCMToken{},
		&model.Faculty{},
		&model.FacultyRoom{},
		&model.MakeupClass{},
		&model.Notification{},
		&model.NotificationTargetUser{},
		&model.Room{},
		&model.RoomChange{},
		&model.Subject{},
		&model.SubjectEligibleAttribute{},
		&model.SubjectFaculty{},
		&model.SubjectRequirement{},
		&model.Syllabus{},
		&model.TimetableItem{},
		&model.TimetableItemRoom{},
		&model.User{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
