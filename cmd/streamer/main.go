package main
import ("flag"; "log"; "time"; "os")
func main() {
	csv := flag.String("csv", "", "path to DCGM csv")
	flag.Parse()
	log.Println("streamer: starting")
	for {
		if _, err := os.Stat(*csv); os.IsNotExist(err) {
			log.Printf("streamer: %s not found, retrying in 5s", *csv)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("streamer: found %s, processing...", *csv)
		time.Sleep(10 * time.Second)
	}
}
