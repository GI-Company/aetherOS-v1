package main

type DesktopApp struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Window string `json:"window"`
}

func getDesktopApps() []DesktopApp {
	return []DesktopApp{
		{
			ID:     "com.aether.marketplace",
			Name:   "Marketplace",
			Icon:   "marketplace-icon.png",
			Window: "MarketplaceWindow",
		},
		{
			ID:     "com.aether.ai_agent",
			Name:   "AI Agent",
			Icon:   "compute-icon.png",
			Window: "AIAgentWindow",
		},
		{
			ID:     "com.aether.compute",
			Name:   "Compute",
			Icon:   "compute-icon.png",
			Window: "ComputeWindow",
		},
	}
}
