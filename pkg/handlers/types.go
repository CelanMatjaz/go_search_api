package handlers

import "github.com/CelanMatjaz/job_application_tracker_api/pkg/types"

type ResPreWithTags = types.RecordWithTags[types.ResumePreset]
type ResSecWithTags = types.RecordWithTags[types.ResumeSection]

type AppPreWithTags = types.RecordWithTags[types.ApplicationPreset]
type AppSecWithTags = types.RecordWithTags[types.ApplicationSection]
