package ngitung

import (
	"math"
)

// Location menyimpan koordinat latitude dan longitude dalam derajat
type Location struct {
	Lat float64
	Lon float64
}

// haversine menghitung jarak antara dua titik di permukaan bumi (dalam km)
func haversine(loc1, loc2 Location) float64 {
	const R = 6371 // Radius bumi dalam kilometer

	lat1 := loc1.Lat * math.Pi / 180
	lon1 := loc1.Lon * math.Pi / 180
	lat2 := loc2.Lat * math.Pi / 180
	lon2 := loc2.Lon * math.Pi / 180

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// calculatePHA menghitung Peak Horizontal Acceleration (PHA) menggunakan rumus Fukushima-Tanaka
func calculatePHA(magnitude, distance float64) float64 {
	// Rumus: log10(A) = 0.41*M - log10(R + 0.032 * 10^(0.41*M)) - 0.0034*R + 1.30
	term1 := 0.41 * magnitude
	term2 := math.Log10(distance + 0.032*math.Pow(10, 0.41*magnitude))
	term3 := 0.0034 * distance
	logA := term1 - term2 - term3 + 1.30
	return math.Pow(10, logA) // A dalam cm/s^2
}

// phaToMMI mengonversi PHA ke skala Modified Mercalli Intensity (MMI)
func phaToMMI(pha float64) float64 {
	// Hubungan kasar: MMI ≈ 1 + 2 * log10(PHA)
	return 1 + 2*math.Log10(pha)
}

// isWithinFeltRadius memeriksa apakah lokasi berada dalam radius guncangan yang terasa (MMI >= 3)
func IsWithinFeltRadius(quakeLoc, targetLoc Location, magnitude float64) (bool, float64, float64) {
	// Hitung jarak menggunakan formula Haversine
	distance := haversine(quakeLoc, targetLoc)

	// Jika jarak sangat kecil, gunakan nilai PHA default (620 cm/s^2) sesuai artikel
	var pha float64
	if distance < 1 { // Untuk jarak sangat dekat (< 1 km)
		pha = 620
	} else {
		pha = calculatePHA(magnitude, distance)
	}

	// Konversi PHA ke MMI
	mmi := phaToMMI(pha)

	// Guncangan terasa jika MMI >= 3
	isFelt := mmi >= 3

	return isFelt, distance, mmi
}

//
//func main() {
//	// Data gempa baru
//	quakeLoc := Location{Lat: -5.42, Lon: 123.12} // Episentrum gempa
//	magnitude := 3.2                              // Magnitudo gempa
//
//	// Lokasi target (Siotapina, Buton, perkiraan koordinat)
//	targetLoc := Location{Lat: -5.36, Lon: 123.16} // Koordinat Siotapina (perkiraan)
//
//	// Periksa apakah lokasi berada dalam radius guncangan yang terasa
//	isFelt, distance, mmi := isWithinFeltRadius(quakeLoc, targetLoc, magnitude)
//
//	// Tampilkan hasil
//	fmt.Printf("Jarak ke lokasi: %.2f km\n", distance)
//	fmt.Printf("Intensitas (MMI): %.2f\n", mmi)
//	if isFelt {
//		fmt.Println("Lokasi berada dalam radius guncangan yang terasa (MMI >= III).")
//	} else {
//		fmt.Println("Lokasi berada di luar radius guncangan yang terasa (MMI < III).")
//	}
//}
