package scanner

type InjectionReport struct {
	VulnerableGETParams  []string `json:"vulnerableGetParams"`
	VulnerablePOSTParams []string `json:"vulnerablePostParams"`
	VulnerableHeaders    []string `json:"vulnerableHeaders"`
	VulnerableCookies    []string `json:"vulnerableCookies"`
}
