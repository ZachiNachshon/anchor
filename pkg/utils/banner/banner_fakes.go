package banner

var CreateFakeBanner = func() *fakeBannerImpl {
	return &fakeBannerImpl{}
}

type fakeBannerImpl struct {
	Banner
	PrintAnchorMock func()
}

func (b *fakeBannerImpl) PrintAnchor() {
	b.PrintAnchorMock()
}
