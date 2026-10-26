package ratelimiter

import "time"

// Kumpulan Policy terpusat sesuai desain kamu
var Policies = map[string]Policy{
	"global":              {MaxRequests: 1200, Window: 1 * time.Minute},
	"register_ip":         {MaxRequests: 60, Window: 1 * time.Hour},
	"login_ip":            {MaxRequests: 600, Window: 1 * time.Minute},
	"login_failure":       {MaxRequests: 10, Window: 15 * time.Minute},
	"verify_email_ip":     {MaxRequests: 60, Window: 1 * time.Minute},
	"email_send_ip":       {MaxRequests: 60, Window: 10 * time.Minute},
	"email_send_account":  {MaxRequests: 3, Window: 10 * time.Minute},
	"password_reset_ip":   {MaxRequests: 60, Window: 10 * time.Minute},
	"otp_request":         {MaxRequests: 3, Window: 1 * time.Hour},
	"user_action_std":     {MaxRequests: 30, Window: 1 * time.Minute},
	"resend_verification": {MaxRequests: 1, Window: 60 * time.Second},
}
