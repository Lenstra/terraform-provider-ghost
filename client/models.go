package client

type Site struct {
	Title               string `json:"title"`
	Description         string `json:"description"`
	Logo                string `json:"logo"`
	Icon                string `json:"icon"`
	CoverImage          string `json:"cover_image"`
	AccentColor         string `json:"accent_color"`
	Locale              string `json:"locale"`
	Url                 string `json:"url"`
	Version             string `json:"version"`
	AllowExternalSignup bool   `json:"allow_external_signup"`
}

type User struct {
	Id                                   string `json:"id"`
	Name                                 string `json:"name"`
	Slug                                 string `json:"slug"`
	Email                                string `json:"email"`
	Status                               string `json:"status"`
	LastSeen                             string `json:"last_seen"`
	CommentNotifications                 bool   `json:"comment_notifications"`
	FreeMemberSignupNotification         bool   `json:"free_member_signup_notification"`
	PaidSubscriptionStartedNotification  bool   `json:"paid_subscription_started_notification"`
	PaidSubscriptionCanceledNotification bool   `json:"paid_subscription_canceled_notification"`
	MentionNotifications                 bool   `json:"mention_notifications"`
	RecommendationNotifications          bool   `json:"recommendation_notifications"`
	MilestoneNotifications               bool   `json:"milestone_notifications"`
	DonationNotifications                bool   `json:"donation_notifications"`
	CreatedAt                            string `json:"created_at"`
	UpdatedAt                            string `json:"updated_at"`
	Url                                  string `json:"url"`
}

type Theme struct {
	Name    string
	Active  bool
	Package any
}

type Webhook struct {
	Id            string `json:"id"`
	Event         string `json:"event"`
	TargetUrl     string `json:"target_url"`
	ApiVersion    string `json:"api_version"`
	IntegrationId string `json:"integration_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
