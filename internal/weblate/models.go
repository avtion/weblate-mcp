package weblate

import (
	"time"
)

// Language represents a Weblate language
type Language struct {
	URL           string `json:"url"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Direction     string `json:"direction"`
	Plural        Plural `json:"plural"`
	WebURL        string `json:"web_url"`
}

// Plural represents pluralization information for a language
type Plural struct {
	ID       int    `json:"id"`
	Source   int    `json:"source"`
	Number   int    `json:"number"`
	Formula  string `json:"formula"`
	Type     int    `json:"type"`
}

// Project represents a Weblate project
type Project struct {
	URL                    string    `json:"url"`
	Name                   string    `json:"name"`
	Slug                   string    `json:"slug"`
	Web                    string    `json:"web,omitempty"`
	SourceLanguage         Language  `json:"source_language"`
	ComponentsListURL      string    `json:"components_list_url"`
	LanguagesURL           string    `json:"languages_url"`
	RepositoryURL          string    `json:"repository_url"`
	StatisticsURL          string    `json:"statistics_url"`
	LockURL                string    `json:"lock_url"`
	ChangesListURL         string    `json:"changes_list_url"`
	TranslationReview      bool      `json:"translation_review"`
	SourceReview           bool      `json:"source_review"`
	SetLanguageTeam        bool      `json:"set_language_team"`
	EnableHooks            bool      `json:"enable_hooks"`
	Instructions           string    `json:"instructions"`
	UseSharedTM            bool      `json:"use_shared_tm"`
	ContributeSharedTM     bool      `json:"contribute_shared_tm"`
	AccessControl          int       `json:"access_control"`
	TranslationStart       bool      `json:"translation_start"`
	SourceLanguageCode     string    `json:"source_language_code"`
	LanguageAliases        string    `json:"language_aliases"`
	WebURL                 string    `json:"web_url,omitempty"`
}

// Component represents a Weblate component
type Component struct {
	URL                   string    `json:"url"`
	Name                  string    `json:"name"`
	Slug                  string    `json:"slug"`
	Project               Project   `json:"project"`
	VCS                   string    `json:"vcs"`
	Repo                  string    `json:"repo"`
	GitExport             string    `json:"git_export"`
	Branch                string    `json:"branch"`
	PushBranch            string    `json:"push_branch,omitempty"`
	Filemask              string    `json:"filemask"`
	Template              string    `json:"template"`
	EditTemplate          bool      `json:"edit_template"`
	NewLang               string    `json:"new_lang"`
	LanguageCode          string    `json:"language_code"`
	LanguageCodeStyle     string    `json:"language_code_style"`
	SourceLanguage        Language  `json:"source_language"`
	CheckFlags            string    `json:"check_flags"`
	Priority              int       `json:"priority"`
	EnforceCheck          bool      `json:"enforce_check"`
	CommitMessage         string    `json:"commit_message"`
	AddMessage            string    `json:"add_message"`
	DeleteMessage         string    `json:"delete_message"`
	MergeMessage          string    `json:"merge_message"`
	AddonMessage          string    `json:"addon_message"`
	PullMessage           string    `json:"pull_message"`
	AllowTranslationPropagation bool `json:"allow_translation_propagation"`
	EnableSuggestions     bool      `json:"enable_suggestions"`
	SuggestionVoting      bool      `json:"suggestion_voting"`
	SuggestionAutoaccept  int       `json:"suggestion_autoaccept"`
	PushOnCommit          bool      `json:"push_on_commit"`
	CommitPendingAge      int       `json:"commit_pending_age"`
	AutoLockError         bool      `json:"auto_lock_error"`
	RepositoryURL         string    `json:"repository_url"`
	TranslationsURL       string    `json:"translations_url"`
	StatisticsURL         string    `json:"statistics_url"`
	LockURL               string    `json:"lock_url"`
	LinksURL              string    `json:"links_url"`
	WebURL                string    `json:"web_url"`
}

// Translation represents a Weblate translation
type Translation struct {
	URL                          string    `json:"url"`
	Language                     Language  `json:"language"`
	Component                    Component `json:"component"`
	LanguageCode                 string    `json:"language_code"`
	ID                          int       `json:"id"`
	Filename                     string    `json:"filename"`
	Revision                     string    `json:"revision"`
	WebURL                       string    `json:"web_url"`
	ShareURL                     string    `json:"share_url"`
	TranslateURL                 string    `json:"translate_url"`
	RepositoryURL                string    `json:"repository_url"`
	StatisticsURL                string    `json:"statistics_url"`
	FileURL                      string    `json:"file_url"`
	ChangesListURL               string    `json:"changes_list_url"`
	UnitsListURL                 string    `json:"units_list_url"`
	Total                        int       `json:"total"`
	TotalWords                   int       `json:"total_words"`
	TotalChars                   int       `json:"total_chars"`
	Fuzzy                        int       `json:"fuzzy"`
	FuzzyPercent                 float64   `json:"fuzzy_percent"`
	Translated                   int       `json:"translated"`
	TranslatedPercent            float64   `json:"translated_percent"`
	TranslatedWords              int       `json:"translated_words"`
	TranslatedWordsPercent       float64   `json:"translated_words_percent"`
	TranslatedChars              int       `json:"translated_chars"`
	TranslatedCharsPercent       float64   `json:"translated_chars_percent"`
	FailingChecks                int       `json:"failing_checks"`
	FailingChecksPercent         float64   `json:"failing_checks_percent"`
	FailingChecksWords           int       `json:"failing_checks_words"`
	FailingChecksWordsPercent    float64   `json:"failing_checks_words_percent"`
	HaveSuggestion               int       `json:"have_suggestion"`
	HaveComment                  int       `json:"have_comment"`
	LastChange                   time.Time `json:"last_change"`
	LastAuthor                   string    `json:"last_author"`
	RepositoryURLPush            string    `json:"repository_url_push"`
	IsTemplate                   bool      `json:"is_template"`
	IsSource                     bool      `json:"is_source"`
}

// Unit represents a Weblate translation unit
type Unit struct {
	URL                string            `json:"url"`
	ID                 int               `json:"id"`
	Source             []string          `json:"source"`
	Previous           []string          `json:"previous"`
	Target             []string          `json:"target"`
	IDHash             string            `json:"id_hash"`
	ContentHash        string            `json:"content_hash"`
	Location           string            `json:"location"`
	Context            string            `json:"context"`
	Note               string            `json:"note"`
	Flags              string            `json:"flags"`
	State              int               `json:"state"`
	Fuzzy              bool              `json:"fuzzy"`
	Translated         bool              `json:"translated"`
	Approved           bool              `json:"approved"`
	Position           int               `json:"position"`
	HasSuggestion      bool              `json:"has_suggestion"`
	HasComment         bool              `json:"has_comment"`
	HasFailingCheck    bool              `json:"has_failing_check"`
	NumWords           int               `json:"num_words"`
	Priority           int               `json:"priority"`
	ID4                string            `json:"id4"`
	Explanation        string            `json:"explanation"`
	ExtraFlags         string            `json:"extra_flags"`
	WebURL             string            `json:"web_url"`
	SourceUnit         string            `json:"source_unit"`
	PendingSuggestions int               `json:"pending_suggestions"`
	Timestamp          time.Time         `json:"timestamp"`
}

// Change represents a Weblate change record
type Change struct {
	URL            string      `json:"url"`
	Unit           string      `json:"unit,omitempty"`
	Component      string      `json:"component"`
	Translation    string      `json:"translation"`
	User           string      `json:"user,omitempty"`
	Author         string      `json:"author,omitempty"`
	Timestamp      time.Time   `json:"timestamp"`
	Action         int         `json:"action"`
	ActionName     string      `json:"action_name"`
	Target         string      `json:"target"`
	ID             int         `json:"id"`
	Details        interface{} `json:"details,omitempty"`
}

// Statistics represents Weblate statistics
type Statistics struct {
	Name                  string  `json:"name"`
	Slug                  string  `json:"slug,omitempty"`
	TranslatedPercent     float64 `json:"translated_percent"`
	ApprovedPercent       float64 `json:"approved_percent,omitempty"`
	ReadonlyPercent       float64 `json:"readonly_percent,omitempty"`
	NottranslatedPercent  float64 `json:"nottranslated_percent,omitempty"`
	Total                 int     `json:"total"`
	Translated            int     `json:"translated"`
	Approved              int     `json:"approved,omitempty"`
	Nottranslated         int     `json:"nottranslated,omitempty"`
	Readonly              int     `json:"readonly,omitempty"`
	WebURL                string  `json:"web_url,omitempty"`
	RepositoryURL         string  `json:"repository_url,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}