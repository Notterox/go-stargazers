package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func humanize(number int) string {
	switch {
	case number >= 1e9:
		return fmt.Sprintf("%.[2]*[1]fgb", float64(number)/1e9, boolToInt(number%1e9 != 0))
	case number >= 1e6:
		return fmt.Sprintf("%.[2]*[1]fgm", float64(number)/1e6, boolToInt(number%1e6 != 0))
	case number >= 1e3:
		return fmt.Sprintf("%.[2]*[1]fk", float64(number)/1e3, boolToInt(number%1e3 != 0))
	default:
		return fmt.Sprint(number)
	}
}

type StargazersPart struct {
	Stargazers_count float64
}

func queryStargazers(repo string) int {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.github.com/repos/"+repo,
		nil)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	req.Header = map[string][]string{
		"Accept": {"application/vnd.github.v3+json"},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 4048)
	body := resp.Body
	builder := make([]byte, 0, 256)
	for {
		read, err := body.Read(buf)

		builder = append(builder, buf[:read]...)

		if err == io.EOF || read == 0 {
			break
		}
	}

	var stargazer StargazersPart
	err = json.Unmarshal(builder, &stargazer)
	if err != nil {
		fmt.Println(err)
	}
	return int(stargazer.Stargazers_count)
}

func handle(w http.ResponseWriter, r *http.Request) {
	const prefix = "/github.com/"
	u := r.URL.Path
	repo := strings.TrimPrefix(u, prefix)
	stargazers := queryStargazers(repo)
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(fmt.Sprintf(template, humanize(stargazers))))
}

func main() {
	http.HandleFunc("/github.com/", handle)

	log.Println("Serving on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

const template = `
<svg width="122" height="28" viewBox="0 0 122 28" fill="none" xmlns="http://www.w3.org/2000/svg">
<text x="96" y="14" font-family="sans-serif" text-anchor="middle" alignment-baseline="middle">10.7k</text>
<rect x="0.5" y="0.5" width="121" height="27" rx="3.5" fill="white"/>
<rect x="0.5" y="0.5" width="63" height="27" fill="#FAFBFC"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M18 6.25C18.14 6.24991 18.2773 6.28901 18.3962 6.36289C18.5151 6.43676 18.611 6.54245 18.673 6.668L20.555 10.483L24.765 11.095C24.9035 11.1151 25.0335 11.1736 25.1405 11.2637C25.2475 11.3539 25.3271 11.4722 25.3704 11.6052C25.4137 11.7383 25.4189 11.8808 25.3854 12.0167C25.3519 12.1525 25.2811 12.2763 25.181 12.374L22.135 15.344L22.854 19.536C22.8777 19.6739 22.8624 19.8157 22.8097 19.9454C22.757 20.0751 22.6691 20.1874 22.5559 20.2697C22.4427 20.352 22.3087 20.401 22.1691 20.4111C22.0295 20.4212 21.8899 20.3921 21.766 20.327L18 18.347L14.234 20.327C14.1102 20.392 13.9707 20.4211 13.8312 20.411C13.6917 20.4009 13.5578 20.352 13.4447 20.2699C13.3315 20.1877 13.2436 20.0755 13.1908 19.946C13.138 19.8165 13.1225 19.6749 13.146 19.537L13.866 15.343L10.818 12.374C10.7176 12.2763 10.6465 12.1525 10.6128 12.0165C10.5792 11.8805 10.5843 11.7378 10.6276 11.6045C10.6709 11.4713 10.7507 11.3528 10.8579 11.2626C10.965 11.1724 11.0953 11.114 11.234 11.094L15.444 10.483L17.327 6.668C17.389 6.54245 17.4849 6.43676 17.6038 6.36289C17.7227 6.28901 17.86 6.24991 18 6.25V6.25ZM18 8.695L16.615 11.5C16.5612 11.6089 16.4818 11.7031 16.3835 11.7745C16.2853 11.8459 16.1712 11.8924 16.051 11.91L12.954 12.36L15.194 14.544C15.2811 14.6289 15.3464 14.7337 15.384 14.8493C15.4216 14.965 15.4305 15.0881 15.41 15.208L14.882 18.292L17.651 16.836C17.7586 16.7794 17.8784 16.7499 18 16.7499C18.1216 16.7499 18.2414 16.7794 18.349 16.836L21.119 18.292L20.589 15.208C20.5685 15.0881 20.5774 14.965 20.615 14.8493C20.6526 14.7337 20.7178 14.6289 20.805 14.544L23.045 12.361L19.949 11.911C19.8288 11.8934 19.7147 11.8469 19.6165 11.7755C19.5182 11.7041 19.4388 11.6099 19.385 11.501L18 8.694V8.695Z" fill="#575D64"/>
<text x="40" y="14" font-size="0.8em" font-family="-apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji" text-anchor="middle" dy=".3em" fill="#575D64">Star</text>
<rect x="0.5" y="0.5" width="63" height="27" stroke="#D9DADC"/>
<rect x="0.5" y="0.5" width="121" height="27" rx="3.5" stroke="#DDDFE1"/>
<text x="96" y="14" font-size="0.8em" font-family="-apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif, Apple Color Emoji, Segoe UI Emoji" text-anchor="middle" dy=".3em" font-weight="600" fill="black">%s</text>
</svg>`
