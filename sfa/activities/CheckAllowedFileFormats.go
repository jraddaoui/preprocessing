package activities

import (
	"context"
	"fmt"

	"github.com/jraddaoui/preprocessing/sfa/fformat"
	"github.com/jraddaoui/preprocessing/sfa/sip"
)

const AllowedFileFormatsName = "allowed-file-formats"

type AllowedFileFormatsActivity struct{}

func NewAllowedFileFormatsActivity() *AllowedFileFormatsActivity {
	return &AllowedFileFormatsActivity{}
}

type AllowedFileFormatsParams struct {
	SipPath string
}

type AllowedFileFormatsResult struct {
	Ok         bool
	Formats    map[string]string
	NotAllowed []string
}

func (md *AllowedFileFormatsActivity) Execute(ctx context.Context, params *AllowedFileFormatsParams) (*AllowedFileFormatsResult, error) {
	res := &AllowedFileFormatsResult{}
	sf := fformat.NewSiegfriedEmbed()
	// TODO(daniel): make allowed list configurable.
	allowed := map[string]struct{}{
		"fmt/276": {},
		"fmt/95":  {},
	}

	s, err := sip.NewSFASip(params.SipPath)
	if err != nil {
		return nil, err
	}

	res.Formats = make(map[string]string)
	for _, path := range s.Files {
		ff, err := sf.Identify(path)
		if err != nil {
			return nil, err
		}
		res.Formats[path] = ff.ID
	}

	for path, formatID := range res.Formats {
		if _, exists := allowed[formatID]; !exists {
			msg := fmt.Sprintf("File format not allowed: %s", path)
			res.NotAllowed = append(res.NotAllowed, msg)
		}
	}

	if len(res.NotAllowed) == 0 {
		res.Ok = true
	}
	return res, nil
}
