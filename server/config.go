package server

import (
	"github.com/HeapOfChaos/goondvr/entity"
	"github.com/HeapOfChaos/goondvr/supabase"
)

var Config *entity.Config
var SupabaseClient *supabase.Client

// InitSupabase initializes the Supabase client if configured
func InitSupabase() {
	if Config != nil && Config.EnableSupabase && Config.SupabaseURL != "" && Config.SupabaseAPIKey != "" {
		SupabaseClient = supabase.NewClient(Config.SupabaseURL, Config.SupabaseAPIKey)
	}
}
