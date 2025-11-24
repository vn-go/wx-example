package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

type ConfigService struct {
	yamFilePath string
	data        *configInfo
}

/*
Get current directory of app
*/
func (c ConfigService) GetCurrentAppDir() string {
	dir, err := os.Getwd()
	if err != nil {
		// Log error or handle it as appropriate for your application
		// For simplicity, we can just panic, but usually, you'd want better error handling
		panic(err)
	}
	return dir
}
func (c *ConfigService) SetConfigFilePath(yamFilePath string) {
	c.yamFilePath = yamFilePath
}

/*
This function reads a YAML file and serializes its content into a configInfo struct.
*/
var configServiceGetOne sync.Once

func (c *ConfigService) Get() *configInfo {
	var err error
	configServiceGetOne.Do(func() {
		c.data, err = c.load()
	})
	if err != nil {
		panic(err)
	}
	return c.data

}

func (c *ConfigService) load() (*configInfo, error) {
	if c.yamFilePath == "" {
		return nil, fmt.Errorf("please call ConfigService.SetConfigFilePath, ref %T", c)
	}

	// 1. Khởi tạo một đối tượng Viper mới
	v := viper.New()

	// Lấy đường dẫn tuyệt đối của file cấu hình
	currentAppDir := c.GetCurrentAppDir()
	fullPath := filepath.Join(currentAppDir, c.yamFilePath)

	// Tách tên file và đường dẫn
	configDir := filepath.Dir(fullPath)                       // Đường dẫn thư mục chứa file
	configBase := filepath.Base(c.yamFilePath)                // Tên file đầy đủ (ví dụ: config.yaml)
	configExt := filepath.Ext(configBase)                     // Phần mở rộng (ví dụ: .yaml)
	configName := configBase[:len(configBase)-len(configExt)] // Tên file không đuôi (ví dụ: config)

	fmt.Println("----load config----------")
	fmt.Printf("File: %s\nDir: %s\n", fullPath, configDir)
	fmt.Println("-------------------------")

	// 2. Cấu hình Viper
	v.SetConfigName(configName)    // Tên file (ví dụ: "config")
	v.SetConfigType(configExt[1:]) // Loại file (ví dụ: "yaml" - bỏ dấu ".")
	v.AddConfigPath(configDir)     // Đường dẫn thư mục

	// Thêm các thiết lập khác nếu cần (ví dụ: Biến môi trường)
	// v.AutomaticEnv()

	// 3. Đọc file cấu hình
	if err := v.ReadInConfig(); err != nil {
		// Kiểm tra xem lỗi có phải là file không tìm thấy hay không
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("file cấu hình không tìm thấy tại %s: %w", fullPath, err)
		}
		return nil, fmt.Errorf("lỗi khi đọc cấu hình: %w", err)
	}

	// 4. Giải mã (Unmarshal) cấu hình vào struct
	var config configInfo
	// Viper sử dụng mapstructure.Decode() bên dưới, nên các tag `mapstructure` sẽ hoạt động
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã cấu hình vào struct: %w", err)
	}

	return &config, nil
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}
