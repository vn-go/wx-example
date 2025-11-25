package middleware

import "net/http"

func CorsMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var cors = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// Cho phép tất cả origin (cẩn thận với sản phẩm thật!)
		if len(r.Header["Origin"]) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization,view-path")
		//w.Header().Set("Access-Control-Allow-Origin", "https://frontend.example.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Nếu là preflight request (OPTIONS), chỉ phản hồi 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Gọi tiếp handler chính
		next.ServeHTTP(w, r)
		//w.Header().Set("X-Process-Time", fmt.Sprintf("%.2fms", float64(time.Since(start).Nanoseconds())/1e6))

	}
	return cors
}
