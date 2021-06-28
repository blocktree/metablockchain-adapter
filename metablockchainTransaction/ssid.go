package metablockchainTransaction

type SsidVc struct {
	Did string `json:"did"`
	PublicKey string `json:"public_key"`
}

type KycVc struct {
	did string
	uid string
	public_key string
}

type RegistVcInfo struct {
	ssid_vc SsidVc
	kyc_vc KycVc
}
