package mail

import (
	"../external/net/html"
	"../model"
	"io"
	"strings"
	"time"
	"errors"
)

func ParseHtml(r io.Reader) (model.Voicemail, error) {
	const (
		Caller int = iota
		Called
		Duration
		Time
		Date
		None
	)

	var caller string = ""
	var called string = ""
	var durationStr string = ""
	var timeStr string = ""
	var dateStr string = ""
	var err error

	var nextTokenDataType = None
	d := html.NewTokenizer(r)
	for {
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := d.Token()
		switch tokenType {
		case html.TextToken: // text between start and end tag
			tokens := strings.TrimSpace(token.String())
			if len(tokens) == 0 {
				continue
			}

			switch nextTokenDataType {
			case Caller:
				nextTokenDataType = None
				caller = tokens
			case Called:
				nextTokenDataType = None
				called = tokens
			case Duration:
				nextTokenDataType = None
				durationStr = tokens
			case Time:
				nextTokenDataType = None
				timeStr = strings.TrimSuffix(tokens, " Uhr")
			case Date:
				nextTokenDataType = None
				dateStr = tokens
			case None:
				switch strings.ToLower(tokens) {
				case "anruf von:":
					nextTokenDataType = Caller
				case "für die rufnummer:":
					nextTokenDataType = Called
				case "uhrzeit:":
					nextTokenDataType = Time
				case "datum:":
					nextTokenDataType = Date
				case "aufnahmelänge:":
					nextTokenDataType = Duration
				}
			}
		}
	}

	if timeStr == "" || dateStr == "" {
		return model.Voicemail{}, errors.New("unable to find time and/or date in message")
	}

	var date time.Time
	date, err = time.ParseInLocation("2.01.2006 15:04", dateStr+" "+timeStr, time.Local)
	if err != nil {
		return model.Voicemail{}, err
	}

	if durationStr == "" {
		return model.Voicemail{}, errors.New("unable to find message duration")
	}

	var duration time.Duration
	duration, err = time.ParseDuration(
		strings.Replace(durationStr, ":", "m", 1) + "s")
	if err != nil {
		return model.Voicemail{}, err
	}

	if caller == "" || called == "" {
		return model.Voicemail{}, errors.New("unable to find caller and/or called in message")
	}

	return model.Voicemail{
		Called:   called,
		Caller:   caller,
		Duration: duration,
		Date:     date,
	}, nil

}
