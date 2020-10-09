package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/irohit427/Blog/backend/global"
	blog "github.com/irohit427/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type authServer struct{}

var userCollection mongo.Collection

func (authServer) Login(_ context.Context, in *blog.LoginRequest) (*blog.AuthResponse, error) {
	username, password := in.GetUsername(), in.GetPassword()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var user global.User
	userCollection.FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": username}}}).Decode(&user)
	if user == global.NilUser {
		return &blog.AuthResponse{}, errors.New("Wrong Login Credentials provided")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &blog.AuthResponse{}, errors.New("Wrong Login Credentials provided")
	}
	return &blog.AuthResponse{Token: user.GetToken()}, nil
}

func (server authServer) Signup(_ context.Context, in *blog.SignupRequest) (*blog.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(username) < 4 || len(username) > 20 || len(email) < 7 || len(email) > 35 || len(password) < 8 || len(password) > 128 || !match {
		return &blog.AuthResponse{}, errors.New("Validation failed")
	}

	res, err := server.UsernameUsed(context.Background(), &blog.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Println("Error returned from UsernameUsed: ", err.Error())
		return &blog.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		return &blog.AuthResponse{}, errors.New("Username is used")
	}

	res, err = server.EmailUsed(context.Background(), &blog.EmailUsedRequest{Email: email})
	if err != nil {
		log.Println("Error returned from EmailUsed: ", err.Error())
		return &blog.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		return &blog.AuthResponse{}, errors.New("Email is used")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := global.User{ID: primitive.NewObjectID(), Username: username, Email: email, Password: string(pw)}

	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error inserting newUser: ", err.Error())
		return &blog.AuthResponse{}, errors.New("Something went wrong")
	}

	return &blog.AuthResponse{Token: newUser.GetToken()}, nil
}

func (authServer) UsernameUsed(_ context.Context, in *blog.UsernameUsedRequest) (*blog.UsedResponse, error) {
	username := in.GetUsername()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	return &blog.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) EmailUsed(_ context.Context, in *blog.EmailUsedRequest) (*blog.UsedResponse, error) {
	email := in.GetEmail()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return &blog.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) AuthUser(_ context.Context, in *blog.AuthUserRequest) (*blog.AuthUserResponse, error) {
	token := in.GetToken()
	user := global.UserFromToken(token)
	return &blog.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil
}

func main() {
	userCollection = *global.DB.Collection("user")
	server := grpc.NewServer()
	blog.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error creating listener: ", err.Error())
	}
	go func() {
		log.Fatal("Serving gRPC: ", server.Serve(listener).Error())
	}()

	grpcWebServer := grpcweb.WrapServer(server)
	httpServer := &http.Server{
		Addr: ":9001",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())
}
