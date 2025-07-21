// Package config provides configuration management for the invoice generator.
// It supports YAML, JSON, and environment variable configuration sources.
package config

import (
	"time"
)

// AppConfig represents the main application configuration
type AppConfig struct {
    Server   ServerConfig   `yaml:"server" json:"server" mapstructure:"server"`
    Invoice  InvoiceConfig  `yaml:"invoice" json:"invoice" mapstructure:"invoice"`
    PDF      PDFConfig      `yaml:"pdf" json:"pdf" mapstructure:"pdf"`
    Template TemplateConfig `yaml:"template" json:"template" mapstructure:"template"`
    Logging  LoggingConfig  `yaml:"logging" json:"logging" mapstructure:"logging"`
}

// ServerConfig represents server configuration for API mode
type ServerConfig struct {
    Port         int           `yaml:"port" json:"port" mapstructure:"port" validate:"min=1,max=65535"`
    Host         string        `yaml:"host" json:"host" mapstructure:"host"`
    ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout" mapstructure:"read_timeout"`
    WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout" mapstructure:"write_timeout"`
    EnableCORS   bool          `yaml:"enable_cors" json:"enable_cors" mapstructure:"enable_cors"`
}

// InvoiceConfig represents invoice-specific configuration
type InvoiceConfig struct {
    DefaultCurrency    string  `yaml:"default_currency" json:"default_currency" mapstructure:"default_currency" validate:"len=3"`
    DefaultDueDays     int     `yaml:"default_due_days" json:"default_due_days" mapstructure:"default_due_days" validate:"gte=0"`
    NumberingStrategy  string  `yaml:"numbering_strategy" json:"numbering_strategy" mapstructure:"numbering_strategy" validate:"oneof=sequential date-based custom"`
    NumberPrefix       string  `yaml:"number_prefix" json:"number_prefix" mapstructure:"number_prefix"`
    DefaultTaxRate     float64 `yaml:"default_tax_rate" json:"default_tax_rate" mapstructure:"default_tax_rate" validate:"gte=0,lte=100"`
}

// PDFConfig represents PDF generation configuration
type PDFConfig struct {
    DPI         int    `yaml:"dpi" json:"dpi" mapstructure:"dpi" validate:"min=72,max=600"`
    PageSize    string `yaml:"page_size" json:"page_size" mapstructure:"page_size" validate:"oneof=A4 A3 A5 Letter Legal"`
    Orientation string `yaml:"orientation" json:"orientation" mapstructure:"orientation" validate:"oneof=Portrait Landscape"`
    MarginTop   string `yaml:"margin_top" json:"margin_top" mapstructure:"margin_top"`
    MarginBottom string `yaml:"margin_bottom" json:"margin_bottom" mapstructure:"margin_bottom"`
    MarginLeft   string `yaml:"margin_left" json:"margin_left" mapstructure:"margin_left"`
    MarginRight  string `yaml:"margin_right" json:"margin_right" mapstructure:"margin_right"`
}

// TemplateConfig represents template configuration
type TemplateConfig struct {
    Theme           string            `yaml:"theme" json:"theme" mapstructure:"theme" validate:"oneof=modern classic minimal"`
    CustomCSS       string            `yaml:"custom_css" json:"custom_css" mapstructure:"custom_css"`
    ShowLogo        bool              `yaml:"show_logo" json:"show_logo" mapstructure:"show_logo"`
    CustomVariables map[string]string `yaml:"custom_variables" json:"custom_variables" mapstructure:"custom_variables"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
    Level      string `yaml:"level" json:"level" mapstructure:"level" validate:"oneof=debug info warn error"`
    Format     string `yaml:"format" json:"format" mapstructure:"format" validate:"oneof=json text"`
    OutputPath string `yaml:"output_path" json:"output_path" mapstructure:"output_path"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *AppConfig {
    return &AppConfig{
        Server: ServerConfig{
            Port:         8080,
            Host:         "localhost",
            ReadTimeout:  30 * time.Second,
            WriteTimeout: 30 * time.Second,
            EnableCORS:   true,
        },
        Invoice: InvoiceConfig{
            DefaultCurrency:   "EUR",
            DefaultDueDays:    30,
            NumberingStrategy: "sequential",
            NumberPrefix:      "",
            DefaultTaxRate:    19.0, // 19% VAT for Germany
        },
        PDF: PDFConfig{
            DPI:         300,
            PageSize:    "A4",
            Orientation: "Portrait",
            MarginTop:   "1cm",
            MarginBottom: "1cm",
            MarginLeft:  "1cm",
            MarginRight: "1cm",
        },
        Template: TemplateConfig{
            Theme:           "modern",
            ShowLogo:        true,
            CustomVariables: make(map[string]string),
        },
        Logging: LoggingConfig{
            Level:      "info",
            Format:     "json",
            OutputPath: "stdout",
        },
    }
}
