package global

const (
	dburi       = "mongodb+srv://rohit123:shlocked221b@cluster0.uqlap.mongodb.net/go_blog_db?retryWrites=true&w=majority"
	dbName      = "go_blog_db"
	performance = 100
)

var (
	jwtSecret = []byte("blogSecret")
)
