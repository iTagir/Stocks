package common

import (
	"net/http"
)

type HTTPResponseFunc func(w http.ResponseWriter, r *http.Request)
