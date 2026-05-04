package scrapping

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// EngineConfig agrupa los parámetros de comportamiento del scraper.
// Ningún valor se hardcodea; todos vienen del caller (service).
type EngineConfig struct {
	Limite          int
	DelayMinMS      int // delay mínimo entre acciones (ms)
	DelayMaxMS      int // delay máximo entre acciones (ms)
	TimeoutGlobal   time.Duration
	TimeoutSelector time.Duration
}

// DefaultEngineConfig retorna valores sensatos para producción.
func DefaultEngineConfig(limite int) EngineConfig {
	return EngineConfig{
		Limite:          limite,
		DelayMinMS:      800,
		DelayMaxMS:      2200,
		TimeoutGlobal:   10 * time.Minute,
		TimeoutSelector: 12 * time.Second,
	}
}

// RunScrapingEngine ejecuta el scraping de Google Maps y retorna los leads crudos.
// Es la única función pública de este archivo; el service no conoce chromedp.
func RunScrapingEngine(keyword, ciudad string, cfg EngineConfig) ([]LeadRaw, error) {

	// ── Contexto con timeout global ──────────────────────────────────────────
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("lang", "es-CO,es"),
			chromedp.UserAgent(randomUserAgent()),
		)...,
	)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	globalCtx, cancelGlobal := context.WithTimeout(ctx, cfg.TimeoutGlobal)
	defer cancelGlobal()

	// ── Navegar a Google Maps con la búsqueda ────────────────────────────────
	searchURL := buildSearchURL(keyword, ciudad)

	if err := chromedp.Run(globalCtx,
		chromedp.Navigate(searchURL),
		chromedp.Sleep(humanDelay(cfg)),
	); err != nil {
		return nil, fmt.Errorf("error navegando a Google Maps: %w", err)
	}

	// ── Esperar que el panel de resultados cargue ────────────────────────────
	if err := waitForSelector(globalCtx, `div[role="feed"]`, cfg.TimeoutSelector); err != nil {
		return nil, fmt.Errorf("panel de resultados no apareció: %w", err)
	}

	// ── Recolectar URLs de fichas ────────────────────────────────────────────
	fichaURLs, err := collectFichaURLs(globalCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("error recolectando fichas: %w", err)
	}

	// ── Extraer datos de cada ficha ──────────────────────────────────────────
	var leads []LeadRaw

	for i, fichaURL := range fichaURLs {
		if i >= cfg.Limite {
			break
		}

		lead, err := extractLeadFromFicha(globalCtx, fichaURL, keyword, ciudad, cfg)
		if err != nil {
			// Error en una ficha no detiene el job; el service acumula el contador
			continue
		}

		leads = append(leads, *lead)
		chromedp.Run(globalCtx, chromedp.Sleep(humanDelay(cfg))) //nolint
	}

	return leads, nil
}

// ─── Helpers internos ─────────────────────────────────────────────────────────

func buildSearchURL(keyword, ciudad string) string {
	query := strings.ReplaceAll(keyword+" "+ciudad, " ", "+")
	return "https://www.google.com/maps/search/" + query
}

// collectFichaURLs hace scroll progresivo en el panel de resultados
// y recolecta los href de cada ficha hasta alcanzar cfg.Limite o el final.
func collectFichaURLs(ctx context.Context, cfg EngineConfig) ([]string, error) {

	seen := make(map[string]struct{})
	var urls []string
	noNewRounds := 0
	maxNoNew := 4 // si tras 4 scrolls no hay nada nuevo, asumimos fin

	for len(urls) < cfg.Limite && noNewRounds < maxNoNew {

		var hrefs []string
		err := chromedp.Run(ctx, chromedp.Evaluate(`
			Array.from(
				document.querySelectorAll('div[role="feed"] a[href*="/maps/place/"]')
			).map(a => a.href)
		`, &hrefs))
		if err != nil {
			return urls, err
		}

		prevLen := len(urls)
		for _, href := range hrefs {
			// Normalizar: quedarse solo con la parte hasta el primer parámetro de Google
			clean := cleanMapsURL(href)
			if _, exists := seen[clean]; !exists && clean != "" {
				seen[clean] = struct{}{}
				urls = append(urls, clean)
			}
		}

		if len(urls) == prevLen {
			noNewRounds++
		} else {
			noNewRounds = 0
		}

		// Scroll en el panel de resultados
		chromedp.Run(ctx, //nolint
			chromedp.Evaluate(`
				const feed = document.querySelector('div[role="feed"]');
				if (feed) feed.scrollBy(0, 800);
			`, nil),
			chromedp.Sleep(humanDelay(cfg)),
		)
	}

	return urls, nil
}

// extractLeadFromFicha navega a la URL de la ficha y extrae todos los campos.
func extractLeadFromFicha(
	ctx context.Context,
	fichaURL, keyword, ciudad string,
	cfg EngineConfig,
) (*LeadRaw, error) {

	if err := chromedp.Run(ctx,
		chromedp.Navigate(fichaURL),
		chromedp.Sleep(humanDelay(cfg)),
	); err != nil {
		return nil, fmt.Errorf("error navegando a ficha %s: %w", fichaURL, err)
	}

	// Esperar que el nombre del negocio esté visible
	if err := waitForSelector(ctx, `h1.DUwDvf`, cfg.TimeoutSelector); err != nil {
		return nil, fmt.Errorf("ficha no cargó: %w", err)
	}

	lead := &LeadRaw{
		TipoNegocio: keyword,
		Ciudad:      ciudad,
		MapsURL:     fichaURL,
	}

	// ── Nombre ───────────────────────────────────────────────────────────────
	chromedp.Run(ctx, chromedp.Text(`h1.DUwDvf`, &lead.Nombre, chromedp.ByQuery)) //nolint

	// ── Dirección ────────────────────────────────────────────────────────────
	var direccion string
	chromedp.Run(ctx, chromedp.AttributeValue( //nolint
		`button[data-item-id="address"]`, "aria-label", &direccion, nil,
	))
	lead.Direccion = cleanAriaLabel(direccion, "Dirección: ")

	// ── Teléfono ─────────────────────────────────────────────────────────────
	var telefono string
	chromedp.Run(ctx, chromedp.AttributeValue( //nolint
		`button[data-item-id^="phone:tel:"]`, "aria-label", &telefono, nil,
	))
	lead.Telefono = cleanAriaLabel(telefono, "Número de teléfono: ")

	// ── Sitio web ────────────────────────────────────────────────────────────
	var sitioWeb string
	chromedp.Run(ctx, chromedp.AttributeValue( //nolint
		`a[data-item-id="authority"]`, "href", &sitioWeb, nil,
	))
	lead.SitioWeb = sitioWeb

	// ── Rating ───────────────────────────────────────────────────────────────
	var ratingStr string
	chromedp.Run(ctx, chromedp.Text( //nolint
		`div.F7nice span[aria-hidden="true"]`, &ratingStr, chromedp.ByQuery,
	))
	lead.Rating = parseFloatSafe(strings.ReplaceAll(ratingStr, ",", "."))

	// ── Total reviews ─────────────────────────────────────────────────────────
	var reviewsStr string
	chromedp.Run(ctx, chromedp.Text( //nolint
		`div.F7nice span[aria-label]`, &reviewsStr, chromedp.ByQuery,
	))
	lead.TotalReviews = parseReviewCount(reviewsStr)

	if lead.Nombre == "" {
		return nil, fmt.Errorf("no se pudo extraer nombre de la ficha")
	}

	return lead, nil
}

// waitForSelector espera a que un selector aparezca en el DOM.
func waitForSelector(ctx context.Context, sel string, timeout time.Duration) error {
	tCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return chromedp.Run(tCtx, chromedp.WaitVisible(sel, chromedp.ByQuery))
}

// humanDelay retorna un tiempo aleatorio entre DelayMinMS y DelayMaxMS.
func humanDelay(cfg EngineConfig) time.Duration {
	ms := cfg.DelayMinMS + rand.Intn(cfg.DelayMaxMS-cfg.DelayMinMS+1)
	return time.Duration(ms) * time.Millisecond
}

// cleanMapsURL normaliza la URL de una ficha eliminando parámetros innecesarios.
func cleanMapsURL(href string) string {
	if idx := strings.Index(href, "?"); idx != -1 {
		return href[:idx]
	}
	return href
}

// cleanAriaLabel quita el prefijo textual que Google Maps agrega al aria-label.
func cleanAriaLabel(value, prefix string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), prefix)
}

func parseFloatSafe(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// parseReviewCount extrae el número de reseñas de strings como "(1.234)" o "1234 reseñas".
func parseReviewCount(s string) int {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.Trim(s, "()")
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0
	}
	return n
}

// randomUserAgent rota entre varios user agents reales para reducir detección.
func randomUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	}
	return agents[rand.Intn(len(agents))]
}
