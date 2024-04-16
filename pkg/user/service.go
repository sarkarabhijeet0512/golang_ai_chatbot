package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"
	"uber_fx_init_folder_structure/pkg/cache"
	"uber_fx_init_folder_structure/utils"
	"uber_fx_init_folder_structure/utils/bot"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_pg "github.com/go-pg/pg/v10"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Service struct {
	conf     *viper.Viper
	log      *logrus.Logger
	Repo     Repository
	s3Config *AWSS3Config
	cache    *cache.Service
}

type AWSS3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

// NewService returns a user service object.
func NewService(conf *viper.Viper, log *logrus.Logger, Repo Repository, cache *cache.Service) *Service {
	s3Config := AWSS3Config{
		AccessKeyID:     conf.GetString(utils.AccessKeyEnv),
		SecretAccessKey: conf.GetString(utils.SecretAccessKey),
		Region:          conf.GetString(utils.Region),
		Bucket:          conf.GetString(utils.BucketName),
	}
	return &Service{
		s3Config: &s3Config,
		conf:     conf,
		log:      log,
		Repo:     Repo,
		cache:    cache,
	}
}

func (s *Service) UpsertUserRegistration(ctx context.Context, user *User) error {
	return s.Repo.upsertUserRegistration(ctx, user)
}
func (s *Service) FetchUserByUsername(ctx context.Context, username string) (*User, error) {
	return s.Repo.fetchUserByUsername(ctx, username)
}
func (s *Service) UserUploadPhoto(ctx context.Context, user UserImages, file multipart.File, fileName, contentType string, sess *session.Session) error {
	uploader := s3manager.NewUploader(sess)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.s3Config.Bucket),
		// ACL:         aws.String("public-read"),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		s.log.Error("Failed to upload file to S3: " + err.Error())
		s.log.WithContext(ctx).Info(up.UploadID)
		err = errors.New("failed to upload file")
		return err
	}

	filepath := "https://" + s.s3Config.Bucket + "." + "s3-" + s.s3Config.Region + ".amazonaws.com/" + fileName
	user.Url = filepath
	return s.Repo.userUploadPhoto(ctx, user)
}

func (s *Service) ProcessMessage(ctx context.Context, message string) (string, error) {
	client := openai.NewClient(s.conf.GetString("OPEN_AI_API_KEY"))
	t := s.CustomFunctionOpenAiParams()

	dialogue := bot.Dialogue(message)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   50,
		Messages:    dialogue,
		Temperature: 0.1,
		Tools:       t,
	})
	if err != nil || len(resp.Choices) != 1 {
		return "", fmt.Errorf("completion error: %v len(choices): %v", err, len(resp.Choices))
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		dialogue = append(dialogue, msg)
		call := msg.ToolCalls[0]
		s.log.Infof("OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
			call.Function.Name, call.Function.Arguments)

		var user User
		if err := json.Unmarshal([]byte(call.Function.Arguments), &user); err != nil {
			return "", err
		}

		var toolResp string
		switch call.Function.Name {
		case "CreateUsername":
			toolResp = s.CreateUsername(user)
		case "FetchPhotos":
			toolResp = fmt.Sprint(s.FetchPhotos(user))
		default:
			return "", fmt.Errorf("unsupported tool call: %s", call.Function.Name)
		}

		dialogue = append(dialogue, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResp,
			Name:       call.Function.Name,
			ToolCallID: call.ID,
		})

		resp, err = client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: dialogue,
			Tools:    t,
		})
		if err != nil || len(resp.Choices) != 1 {
			return "", fmt.Errorf("2nd completion error: %v len(choices): %v", err, len(resp.Choices))
		}
		return resp.Choices[0].Message.Content, nil
	}

	return resp.Choices[0].Message.Content, nil
}

// Function to call the API to retrieve photos based on username
func (s *Service) RetrievePhotos(ctx context.Context, username string) ([]string, error) {
	user, err := s.FetchUserByUsername(ctx, username)
	if err != nil {
		if err == _pg.ErrNoRows {
			return []string{"username not found ask to create a new user to upload photos"}, nil
		}
		return nil, err
	}
	userImages, err := s.Repo.retrievePhotos(ctx, user.ID)
	if err != nil {
		if err == _pg.ErrNoRows {
			return []string{"photos not found ask to upload photos"}, nil
		}
		return nil, err
	}
	arr := []string{}
	for _, image := range userImages {
		arr = append(arr, image.Url)
	}
	return arr, nil
}

func (s *Service) CustomFunctionOpenAiParams() []openai.Tool {
	usernameParam := jsonschema.Definition{
		Type:        jsonschema.String,
		Description: "the username, e.g., Abhi0512",
	}
	createUsernameFunction := openai.FunctionDefinition{
		Name:        "CreateUsername",
		Description: "creates or for uploding photos this will be used ",
		Parameters: jsonschema.Definition{
			Type:       jsonschema.Object,
			Properties: map[string]jsonschema.Definition{"username": usernameParam},
			Required:   []string{"username"},
		},
	}
	fetchPhotosFunction := openai.FunctionDefinition{
		Name:        "FetchPhotos",
		Description: "fetches photos for a given username",
		Parameters: jsonschema.Definition{
			Type:       jsonschema.Object,
			Properties: map[string]jsonschema.Definition{"username": usernameParam},
			Required:   []string{"username"},
		},
	}

	return []openai.Tool{
		{Type: openai.ToolTypeFunction, Function: &createUsernameFunction},
		{Type: openai.ToolTypeFunction, Function: &fetchPhotosFunction},
	}
}

// GetCurrentWeather returns the current weather in a given location
func (s *Service) CreateUsername(user User) string {
	ctx := context.Background()

	userdata, err := s.FetchUserByUsername(ctx, user.Username)
	if err != nil && err != _pg.ErrNoRows {
		return "please try again later"
	}
	if userdata.Username == "" {
		err := s.UpsertUserRegistration(ctx, &user)
		s.log.Info("Creating new user")
		if err != nil {
			return "unable to create new user ask to try again"
		}
		err = s.cache.Set("username", user.Username, 5*time.Minute)
		if err != nil {
			return "something went wrong please try again"
		}
		return user.Username
	}
	err = s.cache.Set("username", userdata.Username, 5*time.Minute)
	if err != nil {
		return "something went wrong please try again"
	}
	return "greet the user welcome back ,if you want to upload photos you can do so by clicking on the upload button in the chat box" + userdata.Username
}

// GetCurrentWeather returns the current weather in a given location
func (s *Service) FetchPhotos(user User) []string {
	ctx := context.Background()
	imagedata, err := s.RetrievePhotos(ctx, user.Username)
	if err != nil {
		return []string{}
	}
	return imagedata
}
