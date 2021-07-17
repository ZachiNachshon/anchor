package banner

var CreateFakeBanner = func() *fakeBannerImpl {
	return &fakeBannerImpl{}
}

type fakeBannerImpl struct {
	Banner
	PrintAnchorBannerMock func()
}

func (b *fakeBannerImpl) PrintAnchorBanner() {
	b.PrintAnchorBannerMock()
}
