package kiosk

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// GrafanaKioskLocal creates a chrome-based kiosk using a local grafana-server account
func GrafanaKioskLocal(cfg *Config) {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.DisableGPU, // needed?
		chromedp.Flag("noerrdialogs", true),
		chromedp.Flag("kiosk", true),
		chromedp.Flag("bwsi", true),
		chromedp.Flag("incognito", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-overlay-scrollbar", true),
		chromedp.Flag("ignore-certificate-errors", cfg.Target.IgnoreCertificateErrors),
		chromedp.Flag("test-type", cfg.Target.IgnoreCertificateErrors),
		chromedp.Flag("force-device-scale-factor", "0.70"),
		chromedp.Flag("disable-pinch", true),
		chromedp.Flag("check-for-update-interval", "604800"),
		chromedp.Flag("window-position", "0,0"),
		chromedp.UserDataDir(dir),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	listenChromeEvents(taskCtx, targetCrashed)

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	var generatedURL = GenerateURL(cfg.Target.URL, cfg.General.Mode, cfg.General.AutoFit, cfg.Target.IsPlayList)
	log.Println("Navigating to ", generatedURL)
	/*
		Launch chrome and login with local user account

		name=user, type=text
		id=inputPassword, type=password, name=password
	*/
	// Give browser time to load next page (this can be prone to failure, explore different options vs sleeping)
	// time.Sleep(2000 * time.Millisecond)

	if err := chromedp.Run(taskCtx,
		// chromedp.EmulateViewport(960, 640, chromedp.EmulateScale(0.5)),
		// chromedp.Tasks{
		// 	emulation.SetDeviceMetricsOverride(480, 320, 1.0, false),
		// 	emulation.SetPageScaleFactor(0.8),
		// },
		chromedp.Tasks{
			// chromedp.Navigate("https://www.nasa.gov/sites/default/files/thumbnails/image/trajectory_gif_cropped.gif"),
			chromedp.Navigate("https://apod.nasa.gov/apod/image/1909/SpiderFly_Spitzer2Mass_4165.jpg"),
			chromedp.Sleep(time.Minute),
			chromedp.Navigate("https://www.nasa.gov/sites/default/files/styles/stem_hero/public/thumbnails/image/edu_nasa_science_at_home_0.jpg"),
			chromedp.Sleep(time.Minute),
			chromedp.Navigate("https://www.nasa.gov/sites/default/files/styles/full_width_feature/public/thumbnails/image/potw2019a.jpg"),
			chromedp.Sleep(time.Minute),
			chromedp.Navigate("https://www.nasa.gov/sites/default/files/styles/full_width_feature/public/thumbnails/image/tycho.jpg"),
			chromedp.Sleep(time.Minute),
			chromedp.Navigate("https://svs.gsfc.nasa.gov/vis/a030000/a030700/a030792/helix-hst-3240x3240_print.jpg"),
			chromedp.Sleep(time.Minute),
			chromedp.Navigate(generatedURL),
			chromedp.Sleep(2 * time.Second),
			chromedp.WaitVisible(`//input[@name="user"]`, chromedp.BySearch),
			chromedp.SendKeys(`//input[@name="user"]`, cfg.Target.Username, chromedp.BySearch),
			chromedp.SendKeys(`//input[@name="password"]`, cfg.Target.Password+kb.Enter, chromedp.BySearch),
			chromedp.WaitVisible(`notinputPassword`, chromedp.ByID),
		},
	); err != nil {
		panic(err)
	}
}
